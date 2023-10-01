package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"time"
)

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

func GetClients(plugins []config.Plugin) []pluginpb.PluginClient {
	var (
		grpcClients      []pluginpb.PluginClient
		wg               sync.WaitGroup
		connectionCancel = 10
	)

	for _, plugin := range plugins {
		wg.Add(1)
		pluginAddress := fmt.Sprintf("%s:%d", plugin.Address, plugin.Port)
		go func(addr string) {
			defer wg.Done()
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal().Str("module", "main").Msgf("gRPC Dial Error(%s): %s", plugin.Name, err)
			}
			// Create a context for the connection check.
			grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
			executeTicker := time.Duration(connectionCancel) * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), executeTicker)
			defer cancel()

			// Block until the connection is ready or until the context times out.
			if err := waitForConnection(ctx, conn); err != nil {
				fmt.Printf("Connection to %s failed: %v\n", addr, err)
				return
			}

			if conn.GetState() == connectivity.Ready {
				log.Info().Str("module", "util").Msgf("Client connected to plugin: %s successfully with address %s", plugin.Name, addr)
			}
		}(pluginAddress)
	}
	wg.Wait()

	return grpcClients
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
			return fmt.Errorf("Connection timeout")
		default:
			// Wait a short period before checking the connection state again.
			time.Sleep(100 * time.Millisecond)
		}
	}
}
