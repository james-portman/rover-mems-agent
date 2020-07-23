package main

// useful: https://github.com/bugst/go-serial/blob/master/serial_windows.go

import (
	// "encoding/binary"
	// "encoding/hex"
	"fmt"
	// "log"
	"time"
	"errors"

	// "github.com/distributed/sers"
	// "net/http"
	// "github.com/gin-gonic/gin"
)

var (
	globalConnected = false
	globalFaults = []string{"not-checked-yet"}
	globalEcuType = ""
	globalUserCommand = ""
	globalAlert = ""

	globalDataOutput = map[string] float32{}
)

func main() {

	fmt.Println("Rover MEMS Diagnostic Agent version 0.0.0")
	fmt.Println("https://rovermems.com/web-app/")

	go runWebserver()
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

	globalEcuType = "1.9"

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
		default:
			return errors.New("Unknown ECU type set")
	}
	if err != nil {
		return err
	}

	return errors.New("Connect loop finished")

}
