package config


// ***** Network-configuration *****
const GlobalPort int = 20007

// ***** System specifications *****
const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ClearRequestVariant int

const (
	CV_all    ClearRequestVariant = iota // Assumes customers enter the elevator even though its moving in the wrong direction
	CV_InDirn                            // Assumes customers only enter the elevator when its moving in the correct direction
)

const SystemsClearRequestVariant  ClearRequestVariant = CV_InDirn
const DoorOpenDurationSec int = 3