package main

// useful: https://github.com/bugst/go-serial/blob/master/serial_windows.go

import (
	"fmt"
	"time"
	"errors"
	"os/exec"
	// "sync"
)

var (
	globalConnected = false
	globalFaults = []string{"not-checked-yet"}
	globalEcuType = ""
	globalUserCommand = ""
	globalAlert = ""

	globalDataOutput = map[string] float32{}
	// globalDataOutputLock = sync.RWMutex{}

	globalAgentVersion = "1.0.0-RC"
)

func main() {

	fmt.Println("Rover MEMS Diagnostic Agent version "+globalAgentVersion)
	fmt.Println("Going to run on http://localhost:8080/")
	fmt.Println("It should automatically open a browser for you")

	go runWebserver()
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8080").Start()

	for true {
		err := connectLoop();
		if err != nil {
			fmt.Println(err)
		}

		// clear the output
		// for k := range globalDataOutput {
	  //   delete(globalDataOutput, k)
		// }

		time.Sleep(3 * time.Second)
	}

}

func connectLoop() error {

	if globalEcuType == "" {
		return errors.New("No ECU type selected yet, go back to the web interface to choose one")
	}

	portList, err := nativeGetPortsList()
	if err != nil {
		return err
	}
	fmt.Println("Found the following ports that I can use:")
	fmt.Println(portList)

	portname := ""

	if len(portList) == 1 {
		fmt.Println("Only found one port so I'm going to use it")
		portname = portList[0]
	} else if len(portList) > 1 {
		return errors.New("TODO: ask the user which port to use")
	} else {
		return errors.New("No serial ports found, check device manager, do you need to install a driver?")
	}

	fmt.Println("Using port:")
	fmt.Println(portname)

	switch globalEcuType {
		case "rc5":
			_, err = readFirstBytesFromPortRc5(portname)
			break
		case "2J":
			_, err = readFirstBytesFromPortTwoj(portname)
			break
		case "1.9":
			_, err = readFirstBytesFromPortEcu19(portname)
			break
		case "3":
			_, err = readFirstBytesFromPortEcu3(portname)
			break
		default:
			return errors.New("Unknown ECU type set")
	}
	if err != nil {
		return err
	}

	return errors.New("Connect loop finished")

}
