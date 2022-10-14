package main

import "fmt"

// Version strings
var (
	Version = "dev"
	Commit  = "N/A"
)

// GetVersion returns version information of VATZ
func GetVersion() string {
	return fmt.Sprintf("%s-%s", Version, Commit)
}
