package actions

import (
	"bufio"
	"fmt"
	"net"

	"time"
)

// func init() {
// 	fmt.Println("Actions Package Init")
// }

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
		//conn.Close()
		return nil, err
	}

	timeoutDuration := 100 * time.Millisecond

	// Set Read Connection Duration
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//read intial device response
	reader := bufio.NewReader(conn)
	resp := make([]byte, 64)

	readTime := time.Millisecond * 100
	for start := time.Now(); ; {
		_, err := reader.Read(resp)
		//fmt.Println("***************************************************Initial read loop")
		time.Sleep(time.Duration(5 * time.Millisecond))
		if err != nil {
			//err = fmt.Errorf("error reading from system: %s", err.Error())
			//fmt.Printf(err.Error())
			break
		}
		fmt.Printf("The initial response is: %s\r\n", resp)
		if time.Since(start) > readTime {
			break
		}
	}
	//fmt.Println("***************************************************Initial read loop end")

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

	timeoutDuration = 200 * time.Millisecond
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	//fmt.Println("Read Command")
	readTime = time.Millisecond * 100
	for start := time.Now(); ; {
		if time.Since(start) > readTime {
			break
		}
		tempResp, err := reader.ReadBytes('\n')
		if len(tempResp) > 0 {
			resp = append(resp, tempResp...)
		}
		//fmt.Println("***************************************************After write read loop")
		time.Sleep(time.Duration(5 * time.Millisecond))
		if err != nil {
			err = fmt.Errorf("error reading from system: %s", err.Error())
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("The second response is: %s\r\n", resp)

	}
	//fmt.Println("***************************************************After write read loop end")

	if err != nil {
		err = fmt.Errorf("error reading from system: %s", err.Error())
		//fmt.Printf(err.Error())
		conn.Close()
		return nil, err
	}

	//catch for failed command response from Atlona
	if string(resp) == "Command FAILED" {
		err = fmt.Errorf("failed command, please check the correct device is selected. device response: %s", string(resp))
		//fmt.Printf(err.Error())
		conn.Close()
		return nil, err
	}

	fmt.Printf("Response from device: %s\n", resp)
	fmt.Println(resp)

	conn.Close()

	return resp, nil
}
