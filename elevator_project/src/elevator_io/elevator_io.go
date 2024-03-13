package elevator_io

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const _pollRate = 20 * time.Millisecond

var _initialized bool = false
var _numFloors int = 4
var _mtx sync.Mutex
var _conn net.Conn

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down MotorDirection = -1
	MD_Stop MotorDirection = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2
)

type ButtonEvent struct {
	BtnFloor int
	BtnType  ButtonType
}

func Init(addr string, numFloors int) {
	if _initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	_numFloors = numFloors
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial("tcp", addr)

	if err != nil {
		panic(err.Error())
	}
	_initialized = true
}

func SetMotorDirection(dir MotorDirection) {
	write([4]byte{1, byte(dir), 0, 0})
}

func SetButtonLamp(button ButtonType, floor int, value bool) {
	write([4]byte{2, byte(button), byte(floor), toByte(value)})
}

func SetFloorIndicator(floor int) {
	write([4]byte{3, byte(floor), 0, 0})
}

func SetDoorOpenLamp(value bool) {
	write([4]byte{4, toByte(value), 0, 0})
}

func SetStopLamp(value bool) {
	write([4]byte{5, toByte(value), 0, 0})
}

func PollButtons(ch_receiver chan<- ButtonEvent) {
	prevButtons := make([][3]bool, _numFloors)
	for {
		time.Sleep(_pollRate)
		for floor := 0; floor < _numFloors; floor++ {
			for btnType := ButtonType(0); btnType < 3; btnType++ {
				button := GetButton(btnType, floor)
				if button != prevButtons[floor][btnType] && !button {
					ch_receiver <- ButtonEvent{floor, ButtonType(btnType)}
				}
				prevButtons[floor][btnType] = button
			}
		}
	}
}

func PollFloorSensor(ch_receiver chan<- int) {
	prevFloor := -1
	for {
		time.Sleep(_pollRate)
		floor := GetFloor()
		if floor != prevFloor && floor != -1 {
			ch_receiver <- floor
		}
		prevFloor = floor
	}
}

func PollStopButton(ch_receiver chan<- bool) {
	prevStop := false
	for {
		time.Sleep(_pollRate)
		stop := GetStop()
		if stop != prevStop {
			ch_receiver <- stop

		}
		prevStop = stop
	}
}

func PollObstructionSwitch(ch_receiver chan<- bool) {
	prevObstruction := false
	for {
		time.Sleep(_pollRate)
		obstruction := GetObstruction()
		if obstruction != prevObstruction {
			ch_receiver <- obstruction
		}
		prevObstruction = obstruction
	}
}

func GetButton(button ButtonType, floor int) bool {
	a := read([4]byte{6, byte(button), byte(floor), 0})
	return toBool(a[1])
}

func GetFloor() int {
	a := read([4]byte{7, 0, 0, 0})
	if a[1] != 0 {
		return int(a[2])
	} else {
		return -1
	}
}

func GetStop() bool {
	a := read([4]byte{8, 0, 0, 0})
	return toBool(a[1])
}

func GetObstruction() bool {
	a := read([4]byte{9, 0, 0, 0})
	return toBool(a[1])
}

func read(in [4]byte) [4]byte {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	var out [4]byte
	_, err = _conn.Read(out[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	return out
}

func write(in [4]byte) {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 {
		b = true
	}
	return b
}
