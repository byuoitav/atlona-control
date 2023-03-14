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
		out  string
		inpt Input
	)

	switch output {
	case "hdmiOutA": //for 6x2
		out = "x1"
	case "hdmiOutB": //for 6x2
		out = "x2"
	case "mirror": //for 6x2 mirrored
		out = "x1,x2"
	case "1": //for 5x1 or other single output switcher
		out = "x1"
	}
	if out == "" {
		err := fmt.Errorf("invalid output: %s", output)
		return inpt, err
	}
	fmt.Println("******************************* Out:", out)

	cmdString := "x" + input + "AV" + out + "\r" //syntax is xYAVxZ Y=input number, Z=output number
	fmt.Println("String to send: ", cmdString)

	resp, err := sendCommand(address, []byte(cmdString))
	if err != nil {
		return inpt, err
	}
	inpt.Input, err = parseResponse(resp, output, out)
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
	out := ""
	fmt.Println("******************************* Output:", output)
	switch output {
	case "hdmiOutA": //for 6x2
		out = "x1"
	case "hdmiOutB": //for 6x2
		out = "x2"
	case "mirror": //for 6x2 mirrored
		out = "x1"
	case "1": //for 5x1 or other single output switcher
		out = "x1"
		//fmt.Println("5x1 case")
	}
	if out == "" {
		err := fmt.Errorf("invalid output: %s", output)
		return input, err
	}
	fmt.Println("******************************* Out:", out)

	payload := []byte("Status\r")
	//fmt.Println("Payload: ", payload)
	resp, err := sendCommand(address, payload)
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
	fmt.Println("Response: ", string(resp))
	responses := strings.Split(string(resp), "\r\n")
	fmt.Printf("Responses: %x \r\n", responses)
	fmt.Printf("Responses: %s \r\n", responses)
	responseContainsOut := false
	for index, value := range responses {
		fmt.Println(index)
		fmt.Println("Slice: ", value)
		if len(value) > 5 {
			fmt.Println("Output Port: ", string(value[4:]))
			fmt.Println("Out: ", out)
			responseContainsOut = strings.Contains(string(value[4:]), out)
		} else {
			continue
		}
		fmt.Println(len(value), responseContainsOut)
		if responseContainsOut {
			fmt.Println("true dat")
			fmt.Println(string(value[1]))
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
