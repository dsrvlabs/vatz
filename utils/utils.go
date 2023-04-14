package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"strconv"
)

func MakeUniqueValue(pName, pAddr string, pPort int) string {
	return pName + pAddr + strconv.Itoa(pPort)
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

// This is internal purpose
func ConvertHashToInput(hashValue string) string {
	// Decode the hash value from a hexadecimal string to a byte array
	hashBytes, _ := hex.DecodeString(hashValue)

	// Convert the byte array to a string
	originalString := string(hashBytes)

	return originalString
}

func GetClients(plugins []config.Plugin) []pluginpb.PluginClient {
	var grpcClients []pluginpb.PluginClient

	if len(plugins) > 0 {
		for _, plugin := range plugins {
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", plugin.Address, plugin.Port), grpc.WithInsecure())
			if err != nil {
				log.Fatal().Str("module", "main").Msgf("gRPC Dial Error(%s): %s", plugin.Name, err)
			}
			grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
		}
	} else {
		// TODO: Is this really neccessary???
		defaultConnectedTarget := "localhost:9091"
		conn, err := grpc.Dial(defaultConnectedTarget, grpc.WithInsecure())
		if err != nil {
			log.Fatal().Str("module", "main").Msgf("gRPC Dial Error: %s", err)
		}

		//TODO: Please, Create a better client functions with static
		grpcClients = append(grpcClients, pluginpb.NewPluginClient(conn))
	}

	return grpcClients
}
