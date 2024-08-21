package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type GClientWithPlugin struct {
	GRPCClient pluginpb.PluginClient
	PluginInfo config.Plugin
}

func MakeUniqueValue(pName, pAddr string, pPort int) string {
	return pName + pAddr + strconv.Itoa(pPort)
}

func ParseBool(str string) bool {
	switch str {
	case "true", "1", "on":
		return true
	case "false", "0", "off":
		return false
	default:
		return false
	}
}

/*
	Those funcs are internal purpose only.

func: ConvertHashToInput
func: UniqueHashValue
*/
func ConvertHashToInput(hashValue string) string {
	// Decode the hash value from a hexadecimal string to a byte array
	hashBytes, _ := hex.DecodeString(hashValue)

	// Convert the byte array to a string
	originalString := string(hashBytes)

	return originalString
}

func UniqueHashValue(inputString string) string {
	// Create a SHA-256 hash object
	h := sha256.New()
	// Write the input string to the hash object
	h.Write([]byte(inputString))
	// Get the 256-bit hash value as a byte array
	hashBytes := h.Sum(nil)
	// Encode the hash value as a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)
	// Truncate the string to 16 characters
	hashString = hashString[:16]
	return hashString
}

func GetClients(plugins []config.Plugin) []GClientWithPlugin {
	var (
		grpcClientWithPlugins []GClientWithPlugin
		wg                    sync.WaitGroup
		connectionCancel      = 10
	)

	for _, plugin := range plugins {
		wg.Add(1)
		pluginAddress := fmt.Sprintf("%s:%d", plugin.Address, plugin.Port)

		go func(addr string, configPlugin config.Plugin) {
			defer wg.Done()
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal().Str("module", "main").Msgf("gRPC Dial Error(%s): %s", configPlugin.Name, err)
			}
			// Create a context for the connection check.

			grpcClientWithPlugins = append(grpcClientWithPlugins, GClientWithPlugin{GRPCClient: pluginpb.NewPluginClient(conn),
				PluginInfo: configPlugin})

			executeTicker := time.Duration(connectionCancel) * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), executeTicker)
			defer cancel()

			// Block until the connection is ready or until the context times out.
			if err := waitForConnection(ctx, conn); err != nil {
				log.Error().Str("module", "util").Msgf("Connection to %s (plugin:%s) failed: %v\n", addr, configPlugin.Name, err)
				return
			}

			if conn.GetState() == connectivity.Ready {
				log.Info().Str("module", "util").Msgf("Client successfully connected to %s (plugin:%s).", addr, configPlugin.Name)
			}
		}(pluginAddress, plugin)
	}
	wg.Wait()

	return grpcClientWithPlugins
}

// waitForConnection blocks until the gRPC connection is ready or the context times out.
func waitForConnection(ctx context.Context, conn *grpc.ClientConn) error {
	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return nil // Connection is ready
		}

		select {
		case <-ctx.Done():
			log.Error().Str("module", "util").Msg("Connection is timed out. Please Check your plugins' status. ")
			return errors.New("")
		default:
			// Wait a short period before checking the connection state again.
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func InitializeChannel() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}

func IsNotifiedEnabledAndSend(dispatcherNotificationFlag string, pluginNotificationFlag string) (bool, bool) {
	isNotifiedEnabled := false
	isSameFlagExists := false
	if dispatcherNotificationFlag != "" {
		isNotifiedEnabled = true
		if dispatcherNotificationFlag == pluginNotificationFlag {
			isSameFlagExists = true
		}
	}
	return isNotifiedEnabled, isSameFlagExists
}
