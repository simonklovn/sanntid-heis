package config

import (
	"driver/network/localip"
	"flag"
	"fmt"
	"os"
)

const (
	// System specifications
	N_FLOORS               int = 4
	N_BUTTONS              int = 3
	DoorOpenDurationSec    int = 3
	IdleTimeOutDurationSec int = 10

	// Heartbeat-configuration
	HeartbeatSleepMillisec int = 100

	// Network-configuration
	DefaultPortPeer   int = 22017
	DefaultPortBcast  int = 22018

	// Backup-configuration
	BackupPort = "22019"
	BackupSendAddr    = "255.255.255.255:" + BackupPort
)

func InitializeConfig() (string, string) {
	var nodeID, port string
	flag.StringVar(&nodeID, "id", getDefaultID(), "ID of this peer")
	flag.StringVar(&port, "port", "15657", "Port of this peer")
	flag.Parse()
	return nodeID, port
}

func getDefaultID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error obtaining local IP:", err)
		return "DISCONNECTED"
	}
	return fmt.Sprintf("peer_%s:%d", localIP, os.Getpid())
}
