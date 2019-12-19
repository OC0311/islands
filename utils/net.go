package utils

import (
	"fmt"
	"log"
	"os"
)

func GetNodeID() string {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Panic("unset node id")
	}
	return nodeID
}

func GetDBPath(filename string) string {
	return fmt.Sprintf("%s_%s", GetNodeID(), filename)
}

func GetWalletPath(filename string) string {
	fmt.Println(os.Getenv("NODE_ID"))
	return fmt.Sprintf("%s_%s", GetNodeID(), filename)
}
