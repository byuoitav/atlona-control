package actions

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
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

Response format:
{"volume":10}
{"muted":false}  (the variable comes in as either a 1/0 or true/false)


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
func SetMute(address string, output string, status bool, device string) (Mute, error) {
	fmt.Printf("Incoming vars: address: %s, output: %s, device: %s\r\n", address, output, device)
	var state Mute

	if status {
		payload := []byte{0x3A, 0x30, 0x31, 0x53, 0x39, 0x30, 0x30, 0x31, 0x0d}
		_, err := sendCommand(address, payload)
		if err != nil {
			return state, err
		}
		status = false
	} else {
		payload := []byte{0x3A, 0x30, 0x31, 0x53, 0x39, 0x30, 0x30, 0x30, 0x0d}
		_, err := sendCommand(address, payload)
		if err != nil {
			return state, err
		}
		status = true
	}
	return state, nil
}

// gain 60 - zoneOut1 - VOUTMute sta
// 62 - zoneOut1, zoneOut2 - VOUTMutex sta - x(output#)=(1, 2, 3, 4)
// 52 - Out1 - VOUTMute sta
func GetMute(address string, output string, device string) (Mute, error) {
	fmt.Printf("Incoming vars: address: %s, output: %s, device: %s\r\n", address, output, device)
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

func convertVolume(volume string, device string) string {
	//fmt.Printf("\nconvertVolume Incoming Values volume: %s, device: %s\n", volume, device)

	vtmp, err := strconv.Atoi(volume)
	v := float64(vtmp)
	if v > 100 {
		v = 100
	}
	if v < 0 {
		v = 0
	}
	outMax := 100.0
	outMin := 0.0
	devHi := 100.0
	devLo := 0.0
	if err != nil {
		fmt.Println(err)
		return "0"
	}

	switch device {
	case "AT-UHD-SW-52ED": //range -80 - 15
		devHi = 15
		devLo = -80
	case "AT-OME-PS62": //range -90 - 10
		devHi = 10
		devLo = -90
	case "AT-GAIN-60": //range   0 - 100
		devHi = 100
		devLo = 0

	}

	vol := ((devHi-devLo)*(v-outMin))/(outMax-outMin) + devLo
	fmt.Printf("Incoming Volume: %s, devHi: %f, devLo: %f, vol: %f\n", volume, devHi, devLo, vol)
	volToSend := int(vol)
	return fmt.Sprint(volToSend)
}

func convertReceiveVolume(volume string, device string) string {
	//fmt.Printf("\nconvertVolume Incoming Values volume: %s, device: %s\n", volume, device)

	vtmp, err := strconv.Atoi(volume)
	v := float64(vtmp)
	// if v > 100 {
	// 	v = 100
	// }
	// if v < 0 {
	// 	v = 0
	// }
	outMax := 100.0
	outMin := 0.0
	devHi := 100.0
	devLo := 0.0
	if err != nil {
		fmt.Println(err)
		return "0"
	}

	switch device {
	case "AT-UHD-SW-52ED": //range -80 - 15
		devHi = 15.0
		devLo = -80.0
	case "AT-OME-PS62": //range -90 - 10
		devHi = 10.0
		devLo = -90
	case "AT-GAIN-60": //range   0 - 100
		devHi = 100
		devLo = 0

	}
	//volToSend := ((devHi-devLo)*(v-outMin))/(outMax-outMin) + devLo
	vol := ((outMax-outMin)*(v-devLo))/(devHi-devLo) + outMin
	fmt.Printf("Incoming Volume: %s, devHi: %f, devLo: %f, vol: %f\n", volume, devHi, devLo, vol)
	volToSend := int(vol)
	return fmt.Sprint(volToSend)
}

// gain 60 - zoneOut1 - VOL x (0-100)
// 62 - zoneOut1, zoneOut2 - VOUTx y - x(output#)=(1, 2, 3, 4; y= -90-10, sta=status)
// 52 - Out1 - VOUT1 y (y= -80-15, sta=status)
func SetVolume(address string, output string, volume string, device string) (Volume, error) {
	fmt.Printf("Incoming vars: address: %s, output: %s, volume: %s, device: %s\r\n", address, output, volume, device)
	var level Volume

	vol := convertVolume(volume, device)
	fmt.Println(vol)
	//zoneOut1, zoneOut2

	cmd := ""
	parseCmd := ""
	switch device {
	case "AT-UHD-SW-52ED":
		cmd = "VOUT1 " + vol + "\r"
		parseCmd = "VOUT1"
	case "AT-OME-PS62":
		switch output {
		case "zoneOut1":
			cmd = "VOUT1 " + vol + "\r" + "VOUT3 " + vol + "\r"
			parseCmd = "VOUT1"
		case "zoneOut2":
			cmd = "VOUT2 " + vol + "\r" + "VOUT4 " + vol + "\r"
			parseCmd = "VOUT2"
		}
	case "AT-GAIN-60":
		cmd = "VOL " + vol + "\r"
		parseCmd = "VOUT1"
	}

	resp, err := sendCommand(address, []byte(cmd))
	if err != nil {
		return level, err
	}
	fmt.Printf("The response is: %s", resp)
	respLevel, err := parseVolumeResponse(resp, output, parseCmd)
	if err != nil {
		return level, err
	}
	stringVol := strconv.Itoa(respLevel)
	tmpVol := convertReceiveVolume(stringVol, device)
	level.Volume, err = strconv.Atoi(tmpVol)
	if err != nil {
		return level, err
	}

	fmt.Println("respLevel: ", respLevel)
	fmt.Println("level: ", level)

	return level, nil
}

// gain 60 - zoneOut1 - VOL sta
// 62 - zoneOut1, zoneOut2 - VOUTx sta (x(output#)=1, 2, 3, 4)
// 52 - Out1 - VOUT1 sta
func GetVolume(address string, output string, device string) (Volume, error) {
	var level Volume
	fmt.Printf("Incoming vars: address: %s, output: %s, device: %s\r\n", address, output, device)

	cmd := ""
	parseCmd := ""
	switch device {
	case "AT-UHD-SW-52ED":
		cmd = "VOUT1 sta\r"
		parseCmd = "VOUT1"
	case "AT-OME-PS62":
		switch output {
		case "zoneOut1":
			cmd = "VOUT1 sta\r"
			parseCmd = "VOUT1"
		case "zoneOut2":
			cmd = "VOUT2 sta\r"
			parseCmd = "VOUT2"
		}
	case "AT-GAIN-60":
		cmd = "VOL sta\r"
		parseCmd = "VOL"
	}

	resp, err := sendCommand(address, []byte(cmd))
	if err != nil {
		return level, err
	}
	fmt.Printf("The response is: %s", resp)

	//parse return
	respLevel, err := parseVolumeResponse(resp, output, parseCmd)
	if err != nil {
		return level, err
	}
	stringVol := strconv.Itoa(respLevel)
	tmpVol := convertReceiveVolume(stringVol, device)
	level.Volume, err = strconv.Atoi(tmpVol)
	if err != nil {
		return level, err
	}

	fmt.Println("respLevel: ", respLevel)
	fmt.Println("level: ", level)

	//convert response to 0-100 instead of the Atlona levels

	return level, nil
}

func parseVolumeResponse(resp []byte, output string, parseCmd string) (input int, err error) {
	//fmt.Printf("Response: %s, output: %s, parseCmd: %s\n", string(resp), output, parseCmd)
	responses := strings.Split(string(resp), "\r\n")
	responseContainsOut := false
	for _, value := range responses {
		fmt.Println("Slice: ", value)
		if len(value) > 5 {
			responseContainsOut = strings.Contains(string(value), parseCmd)
		} else {
			continue
		}
		if responseContainsOut {
			v := strings.Split(string(value), " ")
			input, err = strconv.Atoi(v[1])
			fmt.Println("input: ", input)
			if err != nil {
				return input, err
			}
			return input, nil
		} else {
			err = fmt.Errorf("invalid response: %s", resp)
			continue
		}
	}
	return input, err
}
