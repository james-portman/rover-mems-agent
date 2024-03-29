package main

// TODO: fix the code to skip over failing 0x21 data reads and go to next instead of restarting

import (
	"fmt"
	"time"
	"errors"
	// "encoding/hex"
	"github.com/distributed/sers"
)

var (

	twojLastSentCommand = []byte {}

	twojInitCommand = []byte {0x81, 0x13, 0xF7, 0x81, 0x0C}

	twojStartDiagnostic = []byte {0x10, 0xA0}
	/*
	related/alternates
	default seems to be 4
	just 0x82 sets it to 4 or 5
	0x10 0x80 sets it to 2 or 4
	0x10 0x90 sets it to 0 or 4
	0x10 0xa0 sets it to 1 or 4
	just 0x13 sets it to 3 or 4

	*/

	twojRequestSeed = []byte {0x27, 0x01}
	twojSendKey = []byte {0x27, 0x02}
	twojPingCommand = []byte {0x3E}

	twojClearFaultsCommand = []byte {0x31, 0xCB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	twojFaultsClearedResponse = []byte {0x71, 0xCB}


	twojLearnImmoCommand = []byte {0x31, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // doesn't work when it was 31,d1, trying 31d0 now, needs more auth or something
	twojResponseLearnImmoCommand = []byte {0x71, 0xD0}

	twojRead722Command = []byte {0x23, 0x00, 0x07, 0x22, 0x01}

	// // read a load from the start of the ROM 0x100000
	// // eventually needs to get to 0x1FFFFF
	// twojReadRomCommand = []byte {0x23, 0x10, 0x00, 0x00, 32} // 34 bytes might actually be allowed at a time but 32 will be neater blocks
	// // 4096 reads needed at that size
	// twojReadRomCommandNextAddress = 0x100000
	// twojReadRomCommandContinued = []byte {0x23, 0x10, 0x00, 0x00, 32} // used to loop with new addresses
	// twojReadRomInProgress = false
	// twojReadRomFilename = "rom-dump.bin"
	// twojReadRomStartedTime = time.Now().Unix() // placeholder

	twojRequestService13 = []byte {0x13}

	/*
	serivce 0x30
	most seem to just set some ram values

	*/

	twojRequestService31_d5 = []byte {0x31, 0xd5} // data from start of cal


	// service 0x33 requests - gather secretive data, either done by index or the first value
	twojRequestService33_d5 = []byte {0x33, 0xd5} // data from 0x1a5 ram
	twojRequestService33_c0 = []byte {0x33, 0xc0} // data from 0x3e2 - immobiliser status
	// possibles:
	// 5 = learn or disabled (0xF0F0)
	// 3 = code correct
	// 4 = code wrong
	// 6 = learn or disabled or pin shorted?  (0x0000 or 0xFFFF)
	// unknown - 0, 1, 2, 7

	twojRequestService33_c8 = []byte {0x33, 0xc8} // data from 0x514
	twojRequestService33_d2 = []byte {0x33, 0xd2} // data from 0x4b4
	twojRequestService33_d4 = []byte {0x33, 0xd4} // data loop - 8 bytes of zeros on the bench?
	twojRequestService33_da = []byte {0x33, 0xda} // data from 0x4b4 (seems same as d2 code but separate function)
	twojRequestService33_c1 = []byte {0x33, 0xc1} // data from 509, 386, 387
	twojRequestService33_d7 = []byte {0x33, 0xd7} // checks multiple bit fields to see if clear, returns OK if so, 7F otherwise


/*
Service 0x31

d6 - swap to normal calibration
d7 - swap to NOSELECT calibration and set immo code to 0xffff (disable it?)
*/


	twojRequestData00 = []byte {0x21, 0x00}
	twojRequestData01 = []byte {0x21, 0x01}
	twojRequestData02 = []byte {0x21, 0x02}
	twojRequestData03 = []byte {0x21, 0x03}
	twojRequestData05 = []byte {0x21, 0x05}
	twojRequestData06 = []byte {0x21, 0x06}
	twojRequestData07 = []byte {0x21, 0x07}
	twojRequestData08 = []byte {0x21, 0x08}
	twojRequestData09 = []byte {0x21, 0x09}
	twojRequestData0A = []byte {0x21, 0x0A}
	twojRequestData0B = []byte {0x21, 0x0B}
	twojRequestData0C = []byte {0x21, 0x0C}
	twojRequestData0D = []byte {0x21, 0x0D}
	twojRequestData0F = []byte {0x21, 0x0F}
	twojRequestData10 = []byte {0x21, 0x10}
	twojRequestData11 = []byte {0x21, 0x11}
	twojRequestData12 = []byte {0x21, 0x12}
	twojRequestData13 = []byte {0x21, 0x13}
	twojRequestFaultsCommand = []byte {0x21, 0x19}
	twojRequestData21 = []byte {0x21, 0x21}
	twojRequestData25 = []byte {0x21, 0x25}
	twojRequestData3A = []byte {0x21, 0x3A}

	twojWokeResponse = []byte {0xc1, 0xd5, 0x8f}
	twojStartDiagResponse = []byte {0x50}
	twojSeedResponse = []byte {0x67, 0x01}
	twojSeed = 0
	twojKey = 0
	twojKeyAcceptResponse = []byte {0x67, 0x02}
	twojPongResponse = []byte {0x7E}


	twojFaultsResponse = []byte {0x61, 0x19}

	twojResponseData00 = []byte {0x61, 0x00}
	twojResponseData01 = []byte {0x61, 0x01}
	twojResponseData02 = []byte {0x61, 0x02}
	twojResponseData03 = []byte {0x61, 0x03}
	twojResponseData05 = []byte {0x61, 0x05}
	twojResponseData06 = []byte {0x61, 0x06}
	twojResponseData07 = []byte {0x61, 0x07}
	twojResponseData08 = []byte {0x61, 0x08}
	twojResponseData09 = []byte {0x61, 0x09}
	twojResponseData0A = []byte {0x61, 0x0A}
	twojResponseData0B = []byte {0x61, 0x0B}
	twojResponseData0C = []byte {0x61, 0x0C}
	twojResponseData0D = []byte {0x61, 0x0D}
	twojResponseData0F = []byte {0x61, 0x0F}
	twojResponseData10 = []byte {0x61, 0x10}
	twojResponseData11 = []byte {0x61, 0x11}
	twojResponseData12 = []byte {0x61, 0x12}
	twojResponseData13 = []byte {0x61, 0x13}
	twojResponseData21 = []byte {0x61, 0x21}
	twojResponseData25 = []byte {0x61, 0x25}
	twojResponseData3A = []byte {0x61, 0x3A}

	twojRefusePing = []byte {0x7F, 0x3e, 0x10}

	twojUserCommands = map[string] []byte{
		"clearfaults": twojClearFaultsCommand,
		"learnimmo": twojLearnImmoCommand,
		"read722": twojRead722Command,
		// "readrom": twojReadRomCommand,
		"service13": twojRequestService13,

		"service31_d5": twojRequestService31_d5,

		"service33_d5": twojRequestService33_d5,
	  "service33_c0": twojRequestService33_c0,
	  "service33_c8": twojRequestService33_c8,
	  "service33_d2": twojRequestService33_d2,
	  "service33_d4": twojRequestService33_d4,
	  "service33_da": twojRequestService33_da,
	  "service33_c1": twojRequestService33_c1,
	  "service33_d7": twojRequestService33_d7,
	}

)

func twojSendCommand(sp sers.SerialPort, command []byte) {
	// fmt.Println("twojSendCommand")
	finalCommand := []byte {byte(len(command))}

	for i := 0; i < len(command); i++ {
		finalCommand = append(finalCommand, command[i])
	}

  checksum := 0
	for i := 0; i < len(finalCommand); i++ {
    checksum += int(finalCommand[i])
  }
  checksum = checksum & 0xFF
	finalCommand = append(finalCommand, byte(checksum))
	// fmt.Printf("sending %d bytes \n%s", len(finalCommand), hex.Dump(finalCommand))
	twojLastSentCommand = finalCommand
	sp.Write(finalCommand)
}


// this is the logic of what to do next based on what is received
func twojSendNextCommand(sp sers.SerialPort, previousResponse []byte) {

	if globalUserCommand != "" {
		command, ok := twojUserCommands[globalUserCommand];
		if ok {
			fmt.Println("Running 2J user command")
			globalUserCommand = ""
			twojSendCommand(sp, command)
			return
		} else {
			fmt.Println("Asked to perform a user command but don't understand it")
		}
	}

	globalUserCommand = ""
	if slicesEqual(previousResponse, twojWokeResponse) {
		twojSendCommand(sp, twojStartDiagnostic)

	} else if slicesEqual(previousResponse, twojStartDiagResponse) {
		twojSendCommand(sp, twojRequestSeed)

	} else if slicesEqual(previousResponse[0:2], twojSeedResponse) {
		command := append(twojSendKey, byte(twojKey >> 8))
		command = append(command, byte(twojKey & 0xFF))
		twojSendCommand(sp, command)

	} else if slicesEqual(previousResponse, twojKeyAcceptResponse) {
		twojSendCommand(sp, twojPingCommand)

	} else if slicesEqual(previousResponse, twojPongResponse) {
		twojSendCommand(sp, twojRequestFaultsCommand)

	} else if slicesEqual(previousResponse, twojFaultsClearedResponse) {
		twojSendCommand(sp, twojRequestFaultsCommand)

	} else if slicesEqual(previousResponse, twojResponseLearnImmoCommand) {
		twojSendCommand(sp, twojRequestData00)

	} else if slicesEqual(previousResponse[0:2], twojFaultsResponse[0:2]) { twojSendCommand(sp, twojRequestData00)

	} else if slicesEqual(previousResponse[0:2], twojResponseData00) { twojSendCommand(sp, twojRequestData01)
	} else if slicesEqual(previousResponse[0:2], twojResponseData01) { twojSendCommand(sp, twojRequestData02)
	} else if slicesEqual(previousResponse[0:2], twojResponseData02) { twojSendCommand(sp, twojRequestData03)
	} else if slicesEqual(previousResponse[0:2], twojResponseData03) { twojSendCommand(sp, twojRequestData05)
	} else if slicesEqual(previousResponse[0:2], twojResponseData05) { twojSendCommand(sp, twojRequestData06)
	} else if slicesEqual(previousResponse[0:2], twojResponseData06) { twojSendCommand(sp, twojRequestData07)
	} else if slicesEqual(previousResponse[0:2], twojResponseData07) { twojSendCommand(sp, twojRequestData08)
	} else if slicesEqual(previousResponse[0:2], twojResponseData08) { twojSendCommand(sp, twojRequestData09)
	} else if slicesEqual(previousResponse[0:2], twojResponseData09) { twojSendCommand(sp, twojRequestData0A)
	} else if slicesEqual(previousResponse[0:2], twojResponseData0A) { twojSendCommand(sp, twojRequestData0B)
	} else if slicesEqual(previousResponse[0:2], twojResponseData0B) { twojSendCommand(sp, twojRequestData0C)
	} else if slicesEqual(previousResponse[0:2], twojResponseData0C) { twojSendCommand(sp, twojRequestData0D)
	} else if slicesEqual(previousResponse[0:2], twojResponseData0D) { twojSendCommand(sp, twojRequestData0F)
	} else if slicesEqual(previousResponse[0:2], twojResponseData0F) { twojSendCommand(sp, twojRequestData10)
	} else if slicesEqual(previousResponse[0:2], twojResponseData10) { twojSendCommand(sp, twojRequestData11)
	} else if slicesEqual(previousResponse[0:2], twojResponseData11) { twojSendCommand(sp, twojRequestData12)
	} else if slicesEqual(previousResponse[0:2], twojResponseData12) { twojSendCommand(sp, twojRequestData13)
	} else if slicesEqual(previousResponse[0:2], twojResponseData13) { twojSendCommand(sp, twojRequestData21)
	} else if slicesEqual(previousResponse[0:2], twojResponseData21) { twojSendCommand(sp, twojRequestData25)
	} else if slicesEqual(previousResponse[0:2], twojResponseData25) { twojSendCommand(sp, twojRequestData3A)
	} else if slicesEqual(previousResponse[0:2], twojResponseData3A) { twojSendCommand(sp, twojPingCommand)

	} else if slicesEqual(previousResponse, twojRefusePing) {
		// cope with not authed? it can refuse 3e ping
		twojSendCommand(sp, twojRequestSeed)	

	} else { // fall back to ping
		fmt.Println("Falling back to ping command")
		twojSendCommand(sp, twojPingCommand)
	}

}


func readFirstBytesFromPortTwoj(fn string) ([]byte, error) {

	fmt.Println("Connecting to 2J ECU")
	globalConnected = false

	sp, err := sers.Open(fn)
	if err != nil {
		return nil, err
	}
	defer sp.Close()

	err = sp.SetMode(10400, 8, sers.N, 1, sers.NO_HANDSHAKE)
	if err != nil {
		return nil, err
	}
	// setting:
	// minread = 0: minimal buffering on read, return characters as early as possible
	// timeout = 1.0: time out if after 1.0 seconds nothing is received
	err = sp.SetReadParams(0, 0)
	if err != nil {
		return nil, err
	}
	mode, err := sp.GetMode()
	fmt.Println("Serial cable set to:")
	fmt.Println(mode)

	// TODO: this is not great, it will keep spawning new ones on every reconnect
	// but they are completely locked by the blocking read in Linux
	// they might behave weirdly in Windows, test it
	go serialReadRoutine(sp)

	sp.SetBreak(false)
	time.Sleep(200 * time.Millisecond)

	sp.SetBreak(true)
	time.Sleep(25 * time.Millisecond)
	sp.SetBreak(false)
	time.Sleep(25 * time.Millisecond)

	time.Sleep(50 * time.Millisecond)

	sp.Write(twojInitCommand)
	fmt.Println("Done sending init command")


	buffer := make([]byte, 0)

	lastReceivedData := timestampMs()
	timeoutMs := int64(1000)

	for timestampMs() < lastReceivedData + timeoutMs {
		// this timeout needs changing to time without an answer rather than number of loops
		// readLoops++

		newData := nonBlockingSerialRead()
		if len(newData) > 0 {
			lastReceivedData = timestampMs()
		}
		buffer = append(buffer, newData...)

		// clear leading zeros (from our wake up)
		for len(buffer) > 0 && buffer[0] == 0x00 {
			fmt.Println("Cleared leading zeros")
			buffer = buffer[1:]
		}

		if len(buffer) == 0 {
			// fmt.Println("buffer empty")
			time.Sleep(1 * time.Millisecond) // dont hammer CPU
			continue
		}

		// check for our init echo
		if len(buffer) >= 5 && slicesEqual(buffer[0:len(twojInitCommand)], twojInitCommand) {
			fmt.Println("Got our init echo")
			buffer = buffer[len(twojInitCommand):]
			// fmt.Print(hex.Dump(buffer))
			continue
		}

		// check for full commands - our echos and responses too

		packetSize := int(buffer[0])
		// TODO: check for implausible packet size

		if len(buffer) < packetSize + 2 {
			// fmt.Println("waiting for rest of data packet")
			time.Sleep(1 * time.Millisecond) // dont hammer CPU
			continue
		}

		// TODO: check checksum ?

		actualData := buffer[1:packetSize+1] // doesn't include len or checksum
		fullPacket := buffer[0:packetSize+2] // entire packet
		// fmt.Printf("actual data: got %d bytes \n%s", len(actualData), hex.Dump(actualData))
		// fmt.Printf("fullPacket: got %d bytes \n%s", len(fullPacket), hex.Dump(fullPacket))

		// check for our echo
		if len(twojLastSentCommand) > 0 && len(fullPacket) >= len(twojLastSentCommand) && slicesEqual(fullPacket[0:len(twojLastSentCommand)], twojLastSentCommand) {
			buffer = buffer[len(twojLastSentCommand):]
			// fmt.Println("Got our last command echo, buffer now:")
			// fmt.Print(hex.Dump(buffer))
			// fmt.Println(len(buffer))
			continue
		}

		// fmt.Println("must have a response that needs parsing:")
		// fmt.Print(hex.Dump(fullPacket))
		// fmt.Print(hex.Dump(actualData))
		twojParseResponse(actualData)
		buffer = nil
		time.Sleep(25 * time.Millisecond)
		twojSendNextCommand(sp, actualData)
	}

	if timestampMs() >= lastReceivedData + timeoutMs {
		// fmt.Printf("had buffer data: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		return nil, errors.New("MEMS 2J timed out")
	}
	fmt.Println("fell out of readloop somehow")

	return nil, err
}
