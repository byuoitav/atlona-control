package actions

import (
	"bytes"
	"fmt"
)

/*
2023/03/02 16:40:13.020914 Setting mute on "zoneOut1" to false
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:13.149921 Set mute
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:13.153742 Setting mute on "zoneOut2" to false
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:13.592596 Set mute
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:18.783588 Setting mute on "zoneOut1" to true
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:18.856657 Set mute
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:18.858791 Setting mute on "zoneOut2" to true
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:19.355018 Set mute
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:21.741883 Setting mute on "zoneOut1" to false
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:21.845413 Set mute
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:21.849429 Setting mute on "zoneOut2" to false
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:22.342521 Set mute
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:24.528289 Setting volume on "zoneOut1" to 35
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:24.599127 Set volume
[ITB-1106-SW1.byu.edu] 2023/03/02 16:40:24.603355 Setting volume on "zoneOut2" to 35


gain 60 - zoneOut1 - VOL x (0-100, sta=status) VOUTMute x (on, off, sta)
62 - zoneOut1, zoneOut2 - VOUTx y (x(output#)=1, 2, 3, 4; y= -90-10, sta=status) VOUTMute x (on, off, sta)
52 - Out1 - VOUTx y (x(output#)=1; y= -80-15, sta=status) VOUTMute x (on, off, sta)
*/

type Volume struct {
	Volume int `json:"volume"`
}

type Mute struct {
	Muted bool `json:"mute"`
}

// gain 60 - zoneOut1 - VOUTMute x - x=(on, off)
// 62 - zoneOut1, zoneOut2 - VOUTMutex y x(output#)=(1, 2, 3, 4)  y=(on, off)
// 52 - Out1 - VOUTMute x - x=(on, off)
func SetMute(address string, output string, status bool, device ...string) error {

	if status {
		payload := []byte{0x3A, 0x30, 0x31, 0x53, 0x39, 0x30, 0x30, 0x31, 0x0d}
		_, err := sendCommand(address, payload)
		if err != nil {
			return err
		}
		status = false
	} else {
		payload := []byte{0x3A, 0x30, 0x31, 0x53, 0x39, 0x30, 0x30, 0x30, 0x0d}
		_, err := sendCommand(address, payload)
		if err != nil {
			return err
		}
		status = true
	}
	return nil
}

// gain 60 - zoneOut1 - VOUTMute sta
// 62 - zoneOut1, zoneOut2 - VOUTMutex sta - x(output#)=(1, 2, 3, 4)
// 52 - Out1 - VOUTMute sta
func GetMute(address string, output string, device ...string) (Mute, error) {
	var state Mute

	mute := []byte("")

	unmute := []byte("")

	payload := []byte{0x3A, 0x30, 0x31, 0x47, 0x39, 0x30, 0x30, 0x30, 0x0D}
	resp, err := sendCommand(address, payload)
	if err != nil {
		return Mute{}, err
	} else if bytes.Equal(resp, mute) {
		state.Muted = true
	} else if bytes.Equal(resp, unmute) {
		state.Muted = false
	} else {
		return Mute{}, err
	}

	return state, nil
}

// gain 60 - zoneOut1 - VOL x (0-100)
// 62 - zoneOut1, zoneOut2 - VOUTx y - x(output#)=(1, 2, 3, 4; y= -90-10, sta=status)
// 52 - Out1 - VOUT1 y (y= -80-15, sta=status)
func SetVolume(address string, output string, volume string, device ...string) error {
	//zoneOut1, zoneOut2
	vol := volume
	if len(vol) == 1 {
		vol = "00" + vol
	} else if len(vol) == 2 {
		vol = "0" + vol
	}

	payload := []byte("send")
	_, err := sendCommand(address, payload)

	if err != nil {
		return err
	}

	return nil
}

// gain 60 - zoneOut1 - VOL sta
// 62 - zoneOut1, zoneOut2 - VOUTx sta (x(output#)=1, 2, 3, 4)
// 52 - Out1 - VOUT1 sta
func GetVolume(address string, output string, device ...string) (Volume, error) {
	var level Volume

	//build command
	gain60 := false
	if len(device) > 0 {
		gain60 = true
	}

	command := "VOUT"
	if gain60 {
		command = "VOL"
	}

	out := ""
	switch output {
	case "zoneOut1":
		out = "1"
	case "zoneOut2":
		out = "2"
	case "Out1":
		out = "1"
	}

	//send command
	payload := []byte(command + out + " sta\r")
	resp, err := sendCommand(address, payload)
	if err != nil {
		return Volume{}, err
	}
	fmt.Println(resp)

	//parse return

	return level, nil
}
