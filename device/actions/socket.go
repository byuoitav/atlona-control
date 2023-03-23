package actions

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"time"
)

// Creating Connection
func createConnection(address string, port string) (*net.TCPConn, error) {
	fmt.Printf("Opening raw socket connection with address %s port %s\n", address, port)
	radder, err := net.ResolveTCPAddr("tcp", address+":"+port)
	if err != nil {
		err = fmt.Errorf("error resolving address : %s", err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, radder)
	if err != nil {
		err = fmt.Errorf("error dialing address : %s", err.Error())
		return nil, err
	}

	return conn, nil
}

// SendCommand opens a connection with <addr> and sends the <command> to the via, returning the response, or an error if one occured.
func sendCommand(address string, cmd []byte) ([]byte, error) {
	port := "23"

	// get the connection
	fmt.Printf("\n\nOpening telnet connection with address: %s command decimal: %d, command: %s\n", address, cmd, string(cmd))
	conn, err := createConnection(address, port)
	if err != nil {
		return nil, err
	}

	timeoutDuration := 100 * time.Millisecond

	// Set Read Connection Duration
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//read intial device response upon opening a connection
	reader := bufio.NewReader(conn)
	resp := make([]byte, 64)
	readTime := time.Millisecond * 100
	for start := time.Now(); ; {
		_, err := reader.Read(resp)
		time.Sleep(time.Duration(5 * time.Millisecond))
		if err != nil {
			break
		}

		fmt.Printf("The initial response is: %s\r\n", resp)
		if strings.Contains(string(resp), "Full Connections") {
			err := fmt.Errorf("error in switcher response:  %s", string(resp))
			conn.Close()
			return nil, err
		}

		if time.Since(start) > readTime {
			break
		}
	}
	if len(resp) < 1 {
		err := fmt.Errorf("no initial response, closing connection")
		conn.Close()
		return nil, err
	}

	// write command
	fmt.Println("Write Command: ", string(cmd))
	if len(cmd) > 0 {
		_, err = conn.Write(cmd)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	resp = nil //clear out the connection initialized read

	//read response to command - may get multiple responses (aka multiple \n's)
	timeoutDuration = 500 * time.Millisecond
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	readTime = time.Millisecond * 100
	for start := time.Now(); ; {
		tempResp, err := reader.ReadBytes('\n')
		if len(tempResp) > 0 {
			resp = append(resp, tempResp...)
		}
		if err != nil {
			err = fmt.Errorf("error reading from system: %s", err.Error())
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("The second response is: %s\r\n", resp)
		time.Sleep(time.Duration(10 * time.Millisecond))

		if time.Since(start) > readTime {
			break
		}
	}

	if err != nil {
		conn.Close()
		return nil, err
	}

	//catch for failed command response from Atlona
	if string(resp) == "Command FAILED" {
		err = fmt.Errorf("failed command, please check the correct device is selected. device response: %s", string(resp))
		conn.Close()
		return nil, err
	}

	fmt.Printf("Response from device: %s\n", resp)
	fmt.Println(resp)
	conn.Close()
	return resp, nil
}
