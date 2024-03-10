package config

import (
	"driver/network/localip"
	"flag"
	"fmt"
	"os"
)

const (
	// System specifications
	N_FLOORS            int = 4
	N_BUTTONS           int = 3
	DoorOpenDurationSec int = 3
	// Network-configuration
	DefaultPortPeer   int = 22017
	DefaultPortBcast  int = 22018
	DefaultPortBackup int = 22019
)

func InitializeConfig() (string, string) {
	var id, port string
	flag.StringVar(&id, "id", getDefaultID(), "ID of this peer")
	flag.StringVar(&port, "port", "15657", "Port of this peer")
	flag.Parse()
	return id, port
}

func getDefaultID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error obtaining local IP:", err)
		return "DISCONNECTED"
	}
	return fmt.Sprintf("peer_%s:%d", localIP, os.Getpid())
}
