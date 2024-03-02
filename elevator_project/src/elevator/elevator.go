package elevator

import (
	"driver/config"
	"driver/elevator_io"
	"fmt"
)

// SK: Elevator states
type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

// SK:
type ClearRequestVariant int

const (
	/* GIVEN: Assume everyone waiting for the elevator gets on the elevator, even if
	   they will be traveling in the "wrong" direction for a while */
	CV_all ClearRequestVariant = iota
	CV_InDirn
)

// Struct contain
type Elevator struct {
	Floor     int
	Dirn      elevator_io.MotorDirection
	Requests  [config.N_FLOORS][config.N_BUTTONS]bool
	Behaviour ElevatorBehaviour

	Config struct {
		ClearRequestVariant ClearRequestVariant
		DoorOpenDurationSec float64
	}
}

func ElevBehaviourToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

func ElevDirnToString(d elevator_io.MotorDirection) string {
	switch d {
	case elevator_io.MD_Down:
		return "Down"
	case elevator_io.MD_Stop:
		return "Stop"
	case elevator_io.MD_Up:
		return "Up"
	default:
		return "Unknown"
	}
}

func ElevButtonToString(b elevator_io.ButtonType) string {
	switch b {
	case elevator_io.BT_HallUp:
		return "HallUp"
	case elevator_io.BT_HallDown:
		return "HallDown"
	case elevator_io.BT_Cab:
		return "Cab"
	default:
		return "Unknown"
	}
}

// Prints the state of the elevator
func (es *Elevator) Elevator_print() {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.Floor,
		ElevDirnToString(es.Dirn),
		ElevBehaviourToString(es.Behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := config.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			btnType := elevator_io.ButtonType(btn)
			if ((f == config.N_FLOORS-1) && (btnType == elevator_io.BT_HallUp)) ||
				(f == 0 && btnType == elevator_io.BT_HallDown) {
				fmt.Print("|     ")
			} else {
				if es.Requests[f][btn] {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func UninitializedElevator() Elevator {
	return Elevator{
		Floor:     -1,
		Dirn:      elevator_io.MD_Stop,
		Behaviour: EB_Idle,
	}
}

func GetCabRequests(elevator Elevator) []bool {
	// Create a new slice to store the last column elements
	cabRequests := make([]bool, len(elevator.Requests))

	// Loop through each row and access the last element
	for i, row := range elevator.Requests {
		cabRequests[i] = row[len(row)-1]
	}

	return cabRequests
}
