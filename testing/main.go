package main

import (
	"fmt"
	"os"
	"time"

	"github.com/byuoitav/atlona-control/device/actions"
)

func main() {
	count := 0
	commandTime := 1500 //time in ms to send a command
	for {
		input, err := actions.GetInput("169.254.96.49", "B")
		time.Sleep(time.Duration(time.Duration(commandTime) * time.Millisecond))

		count += 1
		dt := time.Now()
		filename := "log.txt"

		text := fmt.Sprintf("\nCount: %s,     Time: %s,     Response: %s,     Error: %s", fmt.Sprint(count), dt.String(), input, err)

		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
		f.Close()
		fmt.Print(text)
	}
}
