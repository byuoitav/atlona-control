package actions

import (
	"fmt"
	"strings"
)

type Input struct {
	Input string `json:"input,omitempty"`
}

var CurrentInput = "hdmi1"

// SetInput sets the input on the display
// 5x1 response format: {"input":"4:1"}
// 6x2 response format: {"input":"1:hdmiOutA"}
func SetInput(address string, input string, output string) (Input, error) {
	var (
		out       string
		inpt      Input
		parseResp string
	)
	port := "23"
	out = "x1" //default to x1 if nothing matches
	parseResp = "x1"

	switch {
	case strings.Contains(output, "A") || strings.Contains(output, "a") || strings.Contains(output, "1") || strings.Contains(output, "0"):
		//A/a for 6x1, 1 for 5x1 0 for 4x1
		out = "x1"
		parseResp = "x1"
	case strings.Contains(output, "B") || strings.Contains(output, "b"): //for 6x2 out B
		out = "x2"
		parseResp = "x2"
	case strings.Contains(output, "mirror"): //for 6x2 mirrored
		out = "x1,x2"
		parseResp = "x1"
	default:
		err := fmt.Errorf("invalid device. expected to contain A, a, 1, 0, B, b, or mirror: %s", output)
		return inpt, err
	}

	payload := "x" + input + "AV" + out + "\r" //syntax is xYAVxZ Y=input number, Z=output number

	resp, err := sendCommand(address, port, []byte(payload))
	if err != nil {
		return inpt, err
	}
	inpt.Input, err = parseResponse(resp, output, parseResp)
	if err != nil {
		return inpt, err
	}
	return inpt, nil
}

// GetInput returns the input being shown on the display
// 5x1 response format: {"input":"4:1"}
// 6x2 response format: {"input":"1:hdmiOutA"}
func GetInput(address string, output string) (Input, error) {
	var input Input
	port := "23"
	out := ""
	switch {
	case strings.Contains(output, "A") || strings.Contains(output, "a") || strings.Contains(output, "1") || strings.Contains(output, "0"):
		out = "x1"
	case strings.Contains(output, "B") || strings.Contains(output, "b"): //for 6x2
		out = "x2"
	case strings.Contains(output, "mirror"): //for 6x2 mirrored
		out = "x1"
	default:
		err := fmt.Errorf("invalid device. expected to contain A, a, 1, 0, B, b, or mirror: %s", output)
		return input, err
	}

	payload := []byte("Status\r")
	resp, err := sendCommand(address, port, payload)
	if err != nil {
		return input, err
	}

	input.Input, err = parseResponse(resp, output, out)
	if err != nil {
		return input, err
	}
	return input, nil

}

func parseResponse(resp []byte, output string, out string) (input string, err error) {
	responses := strings.Split(string(resp), "\r\n")
	responseContainsOut := false
	for _, value := range responses {
		if len(value) > 5 {
			responseContainsOut = strings.Contains(string(value[4:]), out)
		} else {
			continue
		}
		if responseContainsOut {
			respValue := string(value[1]) + ":" + output
			input = respValue
			return input, nil
		} else {
			err = fmt.Errorf("invalid response: %s", resp)
			continue
		}
	}
	return input, err
}
