package main

// useful: https://github.com/bugst/go-serial/blob/master/serial_windows.go

import (
	"fmt"
	"time"
	"errors"
	"os/exec"
	"sync"
)

var (
	globalConnected = false
	globalFaults = []string{"not-checked-yet"}
	globalSerialPorts = []string{}
	globalSelectedSerialPort = ""
	globalEcuType = ""
	globalUserCommand = ""
	globalAlert = "" // pops up on web UI then closes itself
	globalError = "" // pops up on web UI and stays until closed

	globalDataOutput = map[string] float32{}
	globalDataOutputLock = sync.RWMutex{}

	globalAgentVersion = "1.2.0"

	outgoingData chan string // for pushing data out of the websocket
)

func main() {

	outgoingData = make(chan string, 1000) // buffer on it in case the web browser is slow?
	fmt.Println("################################################################################")
	fmt.Println("# Rover MEMS Diagnostic Agent version "+globalAgentVersion)
	fmt.Println("################################################################################")
	fmt.Println("")
	fmt.Println("Going to run the application at web address http://localhost:8080/")
	fmt.Println("")
	fmt.Println("It should automatically open in your default browser in a moment,")
	fmt.Println("if not then open a browser and go to that address")
	fmt.Println("")
	fmt.Println("If you end up with multiple copies of the web application open then close all but one copy")
	fmt.Println("")
	go runWebserver()
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8080").Start()

	for true {
		err := connectLoop();
		if err != nil {
			// fmt.Println(err)
			globalDataOutputLock.Lock()
			globalError = err.Error()
			globalDataOutputLock.Unlock()
		}

		time.Sleep(1 * time.Second)
	}

}

func connectLoop() error {

	// if globalEcuType == "" {
	// 	return nil
	// 	// return errors.New("No ECU type selected")
	// }

	portList, err := nativeGetPortsList()
	if err != nil {
		return err
	}
	// if len(portList) > 0 {
	// 	fmt.Println("Found the following ports that I can use:")
	// 	fmt.Println(portList)
	// }

	globalDataOutputLock.Lock()
	globalSerialPorts = portList
	globalDataOutputLock.Unlock()

	portname := ""

	if len(portList) == 1 {
		// fmt.Println("Only found one port so I'm going to use it")
		portname = portList[0]

		globalDataOutputLock.Lock()
		globalSelectedSerialPort = portname
		globalDataOutputLock.Unlock()

	} else if len(portList) > 1 {
		globalDataOutputLock.Lock()
		if globalSelectedSerialPort == "" {
			globalDataOutputLock.Unlock()
			// return errors.New("Multiple COM ports found, select one")
			return nil
		} else {
			portname = globalSelectedSerialPort
		}
		globalDataOutputLock.Unlock()
	} else {
		return errors.New("No serial ports found, check device manager, do you need to install a driver?")
	}

	// TODO: send normal logging data straight to UI using "outgoingData"
	// fmt.Println("Using port:")
	// fmt.Println(portname)

	switch globalEcuType {
		case "1.x":
			_, err = readFirstBytesFromPortEcu1x(portname)
			break;
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
		case "":
			return nil
		default:
			return errors.New("Unknown ECU type set")
	}
	if err != nil {
		return err
	}

	return errors.New("Connect loop finished")

}
