package main

import (
	"fmt"
	"time"
	"errors"
	// "encoding/hex"
	"github.com/distributed/sers"
	"strconv"
)

var (
	rc5PingCommand = []byte {0x82, 0x00, 0x7D}
	rc5RequestFaultsCommand = []byte {0x82, 0x33, 0x4A}
	rc5ClearFaultsCommand = []byte {0x82, 0xC3, 0xBA}

	rc5WokeResponse = []byte {0x55, 0x06, 0x3B}
	rc5PongResponse = []byte {0xC2, 0x00, 0x3D}
	rc5FaultsResponse = []byte {0x33}
	rc5FaultsClearedResponse = []byte {0xC2, 0xC3, 0x7A}

	rc5UserCommands = map[string] []byte{
		"clearfaults": rc5ClearFaultsCommand,
	}

	rc5Faults = map[int]string{
		0x150A: "Driver airbag shorted to battery positive",
		0x150B: "Driver airbag shorted to battery negative",
		0x150C: "Driver airbag high resistance",
		0x150D: "Driver airbag low resistance",
		0x150E: "Driver airbag squib circuit",
		0x1512: "Passenger airbag squib short to battery positive",
		0x1513: "Passenger airbag 1 short to battery negative",
		0x1514: "Passenger airbag 1 high resistance",
		0x1515: "Passenger airbag 1 low resistance",
		0x1516: "Passenger airbag 1 squib circuit",
		0x151A: "Pretensioner short to battery positive",
		0x151B: "Pretensioner short to battery negative",
		0x151C: "Passenger airbag 2 high resistance",
		0x151D: "Passenger airbag 2 low resistance",
		0x151E: "Passenger airbag 2 squib circuit",
		0x1524: "Right pretensioner high resistance",
		0x1525: "Right pretensioner low resistance",
		0x1526: "Right pretensioner squib circuit",
		0x152C: "Left pretensioner high resistance",
		0x152D: "Left pretensioner low resistance",
		0x152E: "Left pretensioner squib circuit",
		0x160C: "SRS warning lamp short circuit",
		0x160D: "SRS warning lamp open circuit",
		0x160E: "SRS warning lamp driver",
		0x0000: "0x0000 Unknown fault, power cycle and try again",
	}
)


func rc5SendNextCommand(sp sers.SerialPort, previousResponse []byte) {
	if globalUserCommand != "" {
		command, ok := rc5UserCommands[globalUserCommand];
		if ok {
			globalUserCommand = ""
			sp.Write(command)
			return
		} else {
			fmt.Println("Asked to perform a user command but don't understand it")
		}
	}

	globalUserCommand = ""
	if slicesEqual(previousResponse, rc5PongResponse) {
		sp.Write(rc5RequestFaultsCommand)

	} else if slicesEqual(previousResponse, rc5WokeResponse) || slicesEqual(previousResponse, rc5FaultsResponse) {
		sp.Write(rc5PingCommand)

	} else if slicesEqual(previousResponse, rc5FaultsClearedResponse) {
		sp.Write(rc5RequestFaultsCommand)

	} else { // fall back to ping
		sp.Write(rc5PingCommand)
	}


}

func readFirstBytesFromPortRc5(fn string) ([]byte, error) {

	fmt.Println("Connecting to RC5 ECU")
	globalConnected = false

	sp, err := sers.Open(fn)
	if err != nil {
		return nil, err
	}
	defer sp.Close()

	err = sp.SetMode(2400, 8, sers.N, 1, sers.NO_HANDSHAKE)
	if err != nil {
		return nil, err
	}

	// setting:
	// minread = 0: minimal buffering on read, return characters as early as possible
	// timeout = 1.0: time out if after 1.0 seconds nothing is received
	err = sp.SetReadParams(0, 0.001)
	if err != nil {
		return nil, err
	}

	mode, err := sp.GetMode()
	fmt.Println("Serial cable set to:")
	fmt.Println(mode)

	sp.SetBreak(false)
	time.Sleep(2000 * time.Millisecond)

	sp.SetBreak(true)
	time.Sleep(200 * time.Millisecond)

	sp.SetBreak(false)
	time.Sleep(400 * time.Millisecond)
	sp.SetBreak(true)
	time.Sleep(400 * time.Millisecond)
	sp.SetBreak(false)
	time.Sleep(400 * time.Millisecond)
	sp.SetBreak(true)
	time.Sleep(400 * time.Millisecond)

	sp.SetBreak(false)
	time.Sleep(200 * time.Millisecond)

	// TODO: get rid of this
	// time.Sleep(1000 * time.Millisecond)

	// buffer to read into
	initBuffer := make([]byte, 0)

	initLoops := 0
	initLoopsLimit := 100
	for initLoops < initLoopsLimit {
		initLoops++
		if initLoops > 1 {
			time.Sleep(10 * time.Millisecond)
		}

		var rb []byte
		rb = make([]byte, 128)

		// read
		n, err := sp.Read(rb[:])
		if err != nil { continue }
		if n == 0 { continue }
		// fmt.Println("got some bytes:")
		// fmt.Println(n)
		// fmt.Printf("got %d bytes \n%s", len(initBuffer), hex.Dump(initBuffer))

		// chop down to actual data size
		rb = rb[0:n]
		initBuffer = append(initBuffer, rb...)


		// strip leading zeros (wake breaks causing them)
		for (len(initBuffer) > 0 && initBuffer[0] == 0x00) {
			// fmt.Println("Stripping leading zero")
			initBuffer = initBuffer[1:]
		}

		if len(initBuffer) < 3 {
			continue
		}

		if slicesEqual(initBuffer[0:3], rc5WokeResponse) {
			fmt.Println("RC5 woke up")
			globalConnected = true
			break
		} else {
			return nil, errors.New("Unsure what RC5 sent back, aborting")
		}

	}
	if initLoops >= initLoopsLimit {
		return nil, errors.New("Timed out waiting for RC5 to wake up")
	}

	// wait 700ms after first connect
	// normall sleep 200ms, so 500 extra this time
	time.Sleep(500 * time.Millisecond)






	// go into proper read/write loop here, with ping as first command

	time.Sleep(200 * time.Millisecond)
	rc5SendNextCommand(sp, rc5WokeResponse)

	buffer := make([]byte, 0)

	readLoops := 0
	for readLoops < 100 {
		readLoops++
		if readLoops > 1 {
			time.Sleep(10 * time.Millisecond)
		}

		rb := make([]byte, 128)
		n, _ := sp.Read(rb[:])
		rb = rb[0:n] // chop down to actual data size
		buffer = append(buffer, rb...)
		if n > 0 {
			readLoops = 0 // reset timeout
		}

		if len(buffer) == 0 { continue }

		if len(buffer) >= 3 {

			// check for our echos and throw them away
			if slicesEqual(buffer[0:3], rc5PingCommand) {
				buffer = buffer[3:]
				continue;
			}	else if slicesEqual(buffer[0:3], rc5RequestFaultsCommand) {
				buffer = buffer[3:]
				continue;
			} else if slicesEqual(buffer[0:3], rc5ClearFaultsCommand) {
				buffer = buffer[3:]
				continue;
			}


			if slicesEqual(buffer[0:3], rc5PongResponse) {
				fmt.Println("< PONG from ECU")
				buffer = buffer[3:]
				time.Sleep(200 * time.Millisecond)
				rc5SendNextCommand(sp, rc5PongResponse)
				continue
			}

			if slicesEqual(buffer[0:3], rc5FaultsClearedResponse) {
				fmt.Println("< FAULT CODES CLEARED")
				globalAlert = "ECU reports faults cleared"
				buffer = buffer[3:]
				time.Sleep(200 * time.Millisecond)
				rc5SendNextCommand(sp, rc5FaultsClearedResponse)
				continue
			}

			// faults returned
			if len(buffer) > 2 && buffer[1] == rc5FaultsResponse[0] {
				expectedLength := buffer[0]
				expectedLength = expectedLength - 0xC0 + 1
				if len(buffer) < int(expectedLength) {
					continue
				}
				fmt.Println("< FAULTS Got fault codes!")

				rc5ParseFaults(buffer)

				buffer = buffer[expectedLength:]
				time.Sleep(200 * time.Millisecond)
				rc5SendNextCommand(sp, rc5FaultsResponse)
				continue
			}

		}

	}
	if readLoops == 100 {
		return nil, errors.New("readloop timed out")
	}
	fmt.Println("fell out of readloop")

	return nil, err
}

func rc5ParseFaults(buffer []byte) {
	// fmt.Printf("got %d bytes \n%s", len(buffer), hex.Dump(buffer))
	// remove first 2 bytes which are length/type
	buffer = buffer[2:]
	numFaults := len(buffer)/2
	fmt.Println("num faults:")
	fmt.Println(numFaults)

	faults := []string {}

	i := 0
	for i < numFaults {
		fault := int(buffer[i*2]) << 8
		fault += int(buffer[(i*2)+1])

		faultText, ok := rc5Faults[fault];
		if !ok {
			faultText = "unknown fault: "+strconv.Itoa(fault);
		}
		// fmt.Println(faultText)

		faults = append(faults, faultText)
		i++
	}
	globalFaults = faults
}
