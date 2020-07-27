package main

import (
	"fmt"
	"time"
	"errors"
	"encoding/hex"
	"github.com/distributed/sers"
)

var (
	ecu3InitCommand = []byte {0x1A, 0x9A}
	// ecu3StartDiagnostic = []byte {0x10, 0xA0}
	// ecu3RequestSeed = []byte {0x27, 0x01}
	// ecu3SendKey = []byte {0x27, 0x02}
	// ecu3PingCommand = []byte {0x3E}
	// ecu3ClearFaultsCommand = []byte {0x31, 0xCB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//
	// ecu3RequestFaultsCommand = []byte {0x21, 0x19}
	//
	// ecu3RequestData00 = []byte {0x21, 0x00}
	// ecu3RequestData01 = []byte {0x21, 0x01}
	// ecu3RequestData02 = []byte {0x21, 0x02}
	// ecu3RequestData03 = []byte {0x21, 0x03}
	// ecu3RequestData05 = []byte {0x21, 0x05}
	// ecu3RequestData06 = []byte {0x21, 0x06}
	// ecu3RequestData07 = []byte {0x21, 0x07}
	// ecu3RequestData08 = []byte {0x21, 0x08}
	// ecu3RequestData09 = []byte {0x21, 0x09}
	// ecu3RequestData0A = []byte {0x21, 0x0A}
	// ecu3RequestData0B = []byte {0x21, 0x0B}
	// ecu3RequestData0C = []byte {0x21, 0x0C}
	// ecu3RequestData0D = []byte {0x21, 0x0D}
	// ecu3RequestData0F = []byte {0x21, 0x0F}
	// ecu3RequestData10 = []byte {0x21, 0x10}
	// ecu3RequestData11 = []byte {0x21, 0x11}
	// ecu3RequestData12 = []byte {0x21, 0x12}
	// ecu3RequestData13 = []byte {0x21, 0x13}
	// ecu3RequestData21 = []byte {0x21, 0x21}
	// ecu3RequestData25 = []byte {0x21, 0x25}
	// ecu3RequestData3A = []byte {0x21, 0x3A}
	//
	// ecu3WokeResponse = []byte {0xc1, 0xd5, 0x8f}
	// ecu3StartDiagResponse = []byte {0x50}
	// ecu3SeedResponse = []byte {0x67, 0x01}
	// ecu3Seed = 0
	// ecu3Key = 0
	// ecu3KeyAcceptResponse = []byte {0x67, 0x02}
	// ecu3PongResponse = []byte {0x7E}
	// ecu3FaultsClearedResponse = []byte {0x71, 0xCB}
	//
	// ecu3FaultsResponse = []byte {0x61, 0x19}
	//
	// ecu3ResponseData00 = []byte {0x61, 0x00}
	// ecu3ResponseData01 = []byte {0x61, 0x01}
	// ecu3ResponseData02 = []byte {0x61, 0x02}
	// ecu3ResponseData03 = []byte {0x61, 0x03}
	// ecu3ResponseData05 = []byte {0x61, 0x05}
	// ecu3ResponseData06 = []byte {0x61, 0x06}
	// ecu3ResponseData07 = []byte {0x61, 0x07}
	// ecu3ResponseData08 = []byte {0x61, 0x08}
	// ecu3ResponseData09 = []byte {0x61, 0x09}
	// ecu3ResponseData0A = []byte {0x61, 0x0A}
	// ecu3ResponseData0B = []byte {0x61, 0x0B}
	// ecu3ResponseData0C = []byte {0x61, 0x0C}
	// ecu3ResponseData0D = []byte {0x61, 0x0D}
	// ecu3ResponseData0F = []byte {0x61, 0x0F}
	// ecu3ResponseData10 = []byte {0x61, 0x10}
	// ecu3ResponseData11 = []byte {0x61, 0x11}
	// ecu3ResponseData12 = []byte {0x61, 0x12}
	// ecu3ResponseData13 = []byte {0x61, 0x13}
	// ecu3ResponseData21 = []byte {0x61, 0x21}
	// ecu3ResponseData25 = []byte {0x61, 0x25}
	// ecu3ResponseData3A = []byte {0x61, 0x3A}

	// ecu3UserCommands = map[string] []byte{
	// 	"clearfaults": ecu3ClearFaultsCommand,
	// }
	//
	// ecu3Faults = map[int]string{
	// 	0x150A: "Driver airbag shorted to battery positive",
	// 	0x150B: "Driver airbag shorted to battery negative",
	// 	0x150C: "Driver airbag high resistance",
	// 	0x150D: "Driver airbag low resistance",
	// 	0x150E: "Driver airbag squib circuit",
	// 	0x1512: "Passenger airbag squib short to battery positive",
	// 	0x1513: "Passenger airbag 1 short to battery negative",
	// 	0x1514: "Passenger airbag 1 high resistance",
	// 	0x1515: "Passenger airbag 1 low resistance",
	// 	0x1516: "Passenger airbag 1 squib circuit",
	// 	0x151A: "Pretensioner short to battery positive",
	// 	0x151B: "Pretensioner short to battery negative",
	// 	0x151C: "Passenger airbag 2 high resistance",
	// 	0x151D: "Passenger airbag 2 low resistance",
	// 	0x151E: "Passenger airbag 2 squib circuit",
	// 	0x1524: "Right pretensioner high resistance",
	// 	0x1525: "Right pretensioner low resistance",
	// 	0x1526: "Right pretensioner squib circuit",
	// 	0x152C: "Left pretensioner high resistance",
	// 	0x152D: "Left pretensioner low resistance",
	// 	0x152E: "Left pretensioner squib circuit",
	// 	0x160C: "SRS warning lamp short circuit",
	// 	0x160D: "SRS warning lamp open circuit",
	// 	0x160E: "SRS warning lamp driver",
	// 	0x0000: "0x0000 Unknown fault, power cycle and try again",
	// }
)
//
// func ecu3SendCommand(sp sers.SerialPort, command []byte) {
//
// 	finalCommand := []byte {byte(len(command))}
//
// 	for i := 0; i < len(command); i++ {
// 		finalCommand = append(finalCommand, command[i])
// 	}
//
//   checksum := 0
// 	for i := 0; i < len(finalCommand); i++ {
//     checksum += int(finalCommand[i])
//   }
//   checksum = checksum & 0xFF
// 	finalCommand = append(finalCommand, byte(checksum))
// 	// fmt.Printf("sending %d bytes \n%s", len(finalCommand), hex.Dump(finalCommand))
// 	sp.Write(finalCommand)
// }

//
// func ecu3SendNextCommand(sp sers.SerialPort, previousResponse []byte) {
// 	if globalUserCommand != "" {
// 		command, ok := ecu3UserCommands[globalUserCommand];
// 		if ok {
// 			globalUserCommand = ""
// 			ecu3SendCommand(sp, command)
// 			return
// 		} else {
// 			fmt.Println("Asked to perform a user command but don't understand it")
// 		}
// 	}
//
// 	globalUserCommand = ""
// 	if slicesEqual(previousResponse, ecu3WokeResponse) {
// 		ecu3SendCommand(sp, ecu3StartDiagnostic)
//
// 	} else if slicesEqual(previousResponse, ecu3StartDiagResponse) {
// 		ecu3SendCommand(sp, ecu3RequestSeed)
//
// 	} else if slicesEqual(previousResponse, ecu3SeedResponse) {
// 		command := append(ecu3SendKey, byte(ecu3Key >> 8))
// 		command = append(command, byte(ecu3Key & 0xFF))
// 		ecu3SendCommand(sp, command)
//
// 	} else if slicesEqual(previousResponse, ecu3KeyAcceptResponse) {
// 		ecu3SendCommand(sp, ecu3PingCommand)
//
// 	} else if slicesEqual(previousResponse, ecu3PongResponse) {
// 		ecu3SendCommand(sp, ecu3RequestFaultsCommand)
//
// 	} else if slicesEqual(previousResponse, ecu3FaultsClearedResponse) {
// 		ecu3SendCommand(sp, ecu3RequestFaultsCommand)
//
// 	} else if slicesEqual(previousResponse, ecu3FaultsResponse) { ecu3SendCommand(sp, ecu3RequestData00)
//
// 	} else if slicesEqual(previousResponse, ecu3ResponseData00) { ecu3SendCommand(sp, ecu3RequestData01)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData01) { ecu3SendCommand(sp, ecu3RequestData02)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData02) { ecu3SendCommand(sp, ecu3RequestData03)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData03) { ecu3SendCommand(sp, ecu3RequestData05)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData05) { ecu3SendCommand(sp, ecu3RequestData06)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData06) { ecu3SendCommand(sp, ecu3RequestData07)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData07) { ecu3SendCommand(sp, ecu3RequestData08)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData08) { ecu3SendCommand(sp, ecu3RequestData09)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData09) { ecu3SendCommand(sp, ecu3RequestData0A)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData0A) { ecu3SendCommand(sp, ecu3RequestData0B)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData0B) { ecu3SendCommand(sp, ecu3RequestData0C)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData0C) { ecu3SendCommand(sp, ecu3RequestData0D)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData0D) { ecu3SendCommand(sp, ecu3RequestData0F)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData0F) { ecu3SendCommand(sp, ecu3RequestData10)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData10) { ecu3SendCommand(sp, ecu3RequestData11)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData11) { ecu3SendCommand(sp, ecu3RequestData12)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData12) { ecu3SendCommand(sp, ecu3RequestData13)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData13) { ecu3SendCommand(sp, ecu3RequestData21)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData21) { ecu3SendCommand(sp, ecu3RequestData25)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData25) { ecu3SendCommand(sp, ecu3RequestData3A)
// 	} else if slicesEqual(previousResponse, ecu3ResponseData3A) { ecu3SendCommand(sp, ecu3PingCommand)
//
// 	// } else if slicesEqual(previousResponse, ecu3WokeResponse) || slicesEqual(previousResponse, ecu3FaultsResponse) {
// 	// 	sp.Write(ecu3PingCommand)
// 	//
// 	// } else if slicesEqual(previousResponse, ecu3FaultsClearedResponse) {
// 	// 	sp.Write(ecu3RequestFaultsCommand)
// 	// 	globalAlert = "ECU reports faults cleared"
//
// 	} else { // fall back to ping
// 		ecu3SendCommand(sp, ecu3PingCommand)
// 	}
//
// }


func readFirstBytesFromPortEcu3(fn string) ([]byte, error) {

	fmt.Println("Connecting to MEMS 3 ECU")
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
	err = sp.SetReadParams(0, 0.001)
	if err != nil {
		return nil, err
	}

	mode, err := sp.GetMode()
	fmt.Println("Serial cable set to:")
	fmt.Println(mode)

	sp.Write(ecu3InitCommand)


	buffer := make([]byte, 0)

	readLoops := 0
	readLoopsLimit := 200
	for readLoops < readLoopsLimit {
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

		// // clear leading zeros (from our wake up)
		// for len(buffer) > 0 && buffer[0] == 0x00 {
		// 	buffer = buffer[1:]
		// }

		if len(buffer) == 0 { continue }

		// check for init echo
		// if len(buffer) >= 5 && slicesEqual(buffer[0:5], ecu3InitCommand) {
		// 	// fmt.Println("Got our init echo")
		// 	buffer = buffer[5:]
		// 	continue
		// }

		// check for full commands - our echos and responses too

		// if len(buffer) < int(buffer[0]) + 2 {
		// 	// fmt.Println("waiting for rest of data packet")
		// 	continue
		// }

		// TODO: check checksum ?

		// actualData := buffer[1:int(buffer[0])+1]
		fmt.Printf("buffer: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		//
		// // our echos
		// if slicesEqual(actualData, ecu3PingCommand) {
		// 	// fmt.Println("Got our ping echo")
		// 	buffer = buffer[(len(ecu3PingCommand)+2):]
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3StartDiagnostic) {
		// 	// fmt.Println("Got our start diag echo")
		// 	buffer = buffer[(len(ecu3StartDiagnostic)+2):]
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3RequestSeed) {
		// 	// fmt.Println("Got our seed req echo")
		// 	buffer = buffer[(len(ecu3RequestSeed)+2):]
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3SendKey) {
		// 	// fmt.Println("Got our key send echo")
		// 	buffer = buffer[(len(ecu3SendKey)+2+2):] // extra 2 for key
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3ClearFaultsCommand) {
		// 	// fmt.Println("Got our clear faults echo")
		// 	buffer = buffer[(len(ecu3ClearFaultsCommand)+2):]
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3RequestFaultsCommand) {
		// 	// fmt.Println("Got our request faults echo")
		// 	buffer = buffer[(len(ecu3RequestFaultsCommand)+2):]
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3RequestData00) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData01) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData02) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData03) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData05) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData06) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData07) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData08) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData09) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData0A) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData0B) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData0C) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData0D) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData0F) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData10) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData11) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData12) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData13) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData21) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData25) { buffer = buffer[4:]; continue }
		// if slicesEqual(actualData, ecu3RequestData3A) { buffer = buffer[4:]; continue }
		//
		// // actual responses
		//
		// if slicesEqual(actualData, ecu3WokeResponse) {
		// 	fmt.Println("< ECU woke up")
		// 	buffer = nil
		// 	globalConnected = true
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3WokeResponse)
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3StartDiagResponse) {
		// 	fmt.Println("< Diag mode accepted")
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3StartDiagResponse)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3SeedResponse) {
		// 	fmt.Println("< seed")
		// 	ecu3Seed = int(actualData[2]) << 8
		// 	ecu3Seed += int(actualData[3])
		// 	// do key generation
		// 	ecu3Key = generateKey(ecu3Seed)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3SeedResponse)
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3KeyAcceptResponse) {
		// 	fmt.Println("< Key accepted")
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3KeyAcceptResponse)
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3PongResponse) {
		// 	fmt.Println("< PONG")
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3PongResponse)
		// 	continue
		// }
		// if slicesEqual(actualData, ecu3FaultsClearedResponse) {
		// 	fmt.Println("< FAULT CLEARED")
		// 	globalAlert = "ECU reports faults cleared"
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3FaultsClearedResponse)
		// 	continue
		// }
		//
		// if slicesEqual(actualData[0:2], ecu3FaultsResponse) {
		// 	fmt.Println("< Faults")
		// 	ecu3ParseFaults(actualData)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3FaultsResponse)
		// 	continue
		// }
		//
		//
		// if slicesEqual(actualData[0:2], ecu3ResponseData00) {
		// 	fmt.Println("got data packet 00")
		// 	// don't care?
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData00)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData01) {
		// 	fmt.Println("got data packet 01")
		// 	coolant := int(actualData[2]) << 8
		// 	coolant += int(actualData[3])
		// 	coolantFloat := float32(coolant) - 2732
		// 	coolantFloat /= 10
		// 	globalDataOutput["coolant_temp"] = coolantFloat
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData01)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData02) {
		// 	fmt.Println("got data packet 02")
		// 	oiltemp := int(actualData[2]) << 8
		// 	oiltemp += int(actualData[3])
		// 	oiltempFloat := float32(oiltemp) - 2732
		// 	oiltempFloat /= 10
		// 	globalDataOutput["oil_temp"] = oiltempFloat
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData02)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData03) {
		// 	fmt.Println("got data packet 03")
		// 	iat := int(actualData[2]) << 8
		// 	iat += int(actualData[3])
		// 	iatFloat := float32(iat) - 2732
		// 	iatFloat /= 10
		// 	globalDataOutput["intake_air_temp"] = iatFloat
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData03)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData05) {
		// 	fmt.Println("got data packet 05")
		// 	fueltemp := int(actualData[2]) << 8
		// 	fueltemp += int(actualData[3])
		// 	globalDataOutput["fuel_temp"] = float32(fueltemp)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData05)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData06) {
		// 	fmt.Println("got data packet 06")
		// 	// don't care?
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData06)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData07) {
		// 	fmt.Println("got data packet 07")
		// 	mapkpa := int(actualData[2]) << 8
		// 	mapkpa += int(actualData[3])
		// 	globalDataOutput["map_sensor_kpa"] = float32(mapkpa)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData07)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData08) {
		// 	fmt.Println("got data packet 08")
		// 	tps := int(actualData[2]) << 8
		// 	tps += int(actualData[3])
		// 	tpsFloat := float32(tps) / 100
		// 	globalDataOutput["tps_degrees"] = tpsFloat
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData08)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData09) {
		// 	fmt.Println("got data packet 09")
		// 	rpm := int(actualData[2]) << 8
		// 	rpm += int(actualData[3])
		// 	globalDataOutput["rpm"] = float32(rpm)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData09)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData0A) {
		// 	fmt.Println("got data packet 0A")
		// 	feedback := int(actualData[2]) << 8
		// 	feedback += int(actualData[3])
		// 	feedbackFloat := float32(feedback) / 100
		// 	globalDataOutput["fuelling_feedback_percent"] = feedbackFloat
		//
		// 	o2mv := int(actualData[4]) << 8
		// 	o2mv += int(actualData[5])
		// 	globalDataOutput["o2_mv"] = float32(o2mv)
		//
		// 	airFuel := ((float32(o2mv) / 1000) * 2) + 10
		// 	globalDataOutput["estimate_air_fuel"] = airFuel
		//
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData0A)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData0B) {
		// 	fmt.Println("got data packet 0B")
		// 	globalDataOutput["coil_1_charge_time"] = float32(actualData[2]) / 1000
		// 	globalDataOutput["coil_2_charge_time"] = float32(actualData[3]) / 1000
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData0B)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData0C) {
		// 	fmt.Println("got data packet 0C")
		// 	globalDataOutput["injector_1_pw"] = float32(actualData[2])
		// 	globalDataOutput["injector_2_pw"] = float32(actualData[3])
		// 	globalDataOutput["injector_3_pw"] = float32(actualData[4])
		// 	globalDataOutput["injector_4_pw"] = float32(actualData[5])
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData0C)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData0D) {
		// 	fmt.Println("got data packet 0D")
		// 	globalDataOutput["vehicle_speed"] = float32(actualData[2])
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData0D)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData0F) {
		// 	fmt.Println("got data packet 0F")
		// 	globalDataOutput["throttle_switch"] = float32(int(actualData[2]) & 1) // 0b00000001
		// 	globalDataOutput["ignition"] = float32((int(actualData[2]) >> 1) & 1) // 0b00000010
		// 	globalDataOutput["ac_button"] = float32((int(actualData[2]) >> 3) & 1) // 0b00001000
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData0F)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData10) {
		// 	fmt.Println("got data packet 10")
		// 	battery := int(actualData[4]) << 8
		// 	battery += int(actualData[5])
		// 	batteryFloat := float32(battery) / 1000
		// 	globalDataOutput["battery_voltage"] = batteryFloat
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData10)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData11) {
		// 	fmt.Println("got data packet 11")
		// 	// 0 means OK, 1 bad
		// 	// will swap for our purposes
		// 	// output is 1 for yes
		// 	camSync := actualData[2] & 1 // 0b00000001
		// 	globalDataOutput["cam_sync"] = float32(1 - camSync)
		// 	crankSync := (actualData[2] >> 1) & 1 //0b00000010
		// 	globalDataOutput["crank_sync"] = float32(1 - crankSync)
		//
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData11)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData12) {
		// 	fmt.Println("got data packet 12")
		// 	idleValvePos := int(actualData[2]) << 8
		// 	idleValvePos += int(actualData[3])
		// 	idleValveFloat := float32(idleValvePos) / 2
		// 	globalDataOutput["idle_valve_pos"] = idleValveFloat
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData12)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData13) {
		// 	fmt.Println("got data packet 13")
		// 	globalDataOutput["closed_loop"] = float32(actualData[2] & 0b00000001)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData13)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData21) {
		// 	fmt.Println("got data packet 21")
		// 	rpmError := int(actualData[2]) << 8
		// 	rpmError += int(actualData[3])
		// 	if rpmError > 32768 {
		// 		rpmError -= 65535
		// 	}
		// 	globalDataOutput["rpm_error"] = float32(rpmError)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData21)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData25) {
		// 	fmt.Println("got data packet 25")
		// 	camPercent := int(actualData[2]) << 8
		// 	camPercent += int(actualData[3])
		// 	globalDataOutput["cam_percent"] = float32(camPercent)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData25)
		// 	continue
		// }
		// if slicesEqual(actualData[0:2], ecu3ResponseData3A) {
		// 	fmt.Println("got data packet 3A")
		// 	idleTimingOffset := int(actualData[2]) << 8
		// 	idleTimingOffset += int(actualData[3])
		// 	idleTimingOffsetFloat := float32(idleTimingOffset) / 10
		// 	globalDataOutput["idle_timing_offset"] = idleTimingOffsetFloat
		//
		// 	idleAdjusterRpm := int(actualData[4]) << 8
		// 	idleAdjusterRpm += int(actualData[5])
		// 	globalDataOutput["idle_adjuster_rpm"] = float32(idleAdjusterRpm)
		// 	buffer = nil
		// 	time.Sleep(50 * time.Millisecond)
		// 	ecu3SendNextCommand(sp, ecu3ResponseData3A)
		// 	continue
		// }

		// if we get here then something is wrong with the data

		// todo: cope with 7F (fail/no)

		// fmt.Printf("unknown command in buffer (burning it): got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		// buffer = nil
		// ecu3SendNextCommand(sp, ecu3PongResponse)

	}
	if readLoops >= readLoopsLimit {
		fmt.Printf("had buffer data: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		return nil, errors.New("readloop timed out")
	}
	fmt.Println("fell out of readloop")

	return nil, err
}
//
// func ecu3ParseFaults(buffer []byte) {
// 	// fmt.Printf("got %d bytes \n%s", len(buffer), hex.Dump(buffer))
//
// 	faults := []string {}
//
// 	if len(buffer) >= 5 {
// 		if (buffer[4] & 0b01000000) > 0 {	faults = append(faults, "Outside air temp (low voltage)") }
// 	  if (buffer[4] & 0b00100000) > 0 { faults = append(faults, "Power supply (low voltage)") }
// 	  if (buffer[4] & 0b00010000) > 0 { faults = append(faults, "Engine oil temp (low voltage)") }
// 	  if (buffer[4] & 0b00000100) > 0 { faults = append(faults, "Coolant temp (low voltage)") }
// 	  if (buffer[4] & 0b00000001) > 0 { faults = append(faults, "System (low voltage)") }
// 	}
//
//
// 	if len(buffer) >= 6 {
// 	  if (buffer[5] & 0b10000000) > 0 { faults = append(faults, "Battery (low voltage)") }
// 	  if (buffer[5] & 0b00010000) > 0 { faults = append(faults, "Lambda 1 bank 1 (low voltage)") }
// 	  if (buffer[5] & 0b00000100) > 0 { faults = append(faults, "Throttle pot (low voltage)") }
// 	  if (buffer[5] & 0b00000010) > 0 { faults = append(faults, "Air intake (low voltage)") }
// 	  if (buffer[5] & 0b00000001) > 0 { faults = append(faults, "MAP sensor (low voltage)") }
// 	}
//
// 	if len(buffer) >= 9 {
// 	  if (buffer[8] & 0b01000000) > 0 { faults = append(faults, "Outside air temp (high voltage)") }
// 	  if (buffer[8] & 0b00100000) > 0 { faults = append(faults, "Power supply (high voltage)") }
// 	  if (buffer[8] & 0b00010000) > 0 { faults = append(faults, "Oil temperature (high voltage)") }
// 	  if (buffer[8] & 0b00000100) > 0 { faults = append(faults, "Coolant temperature (high voltage)") }
// 	  if (buffer[8] & 0b00000001) > 0 { faults = append(faults, "System (high voltage)") }
// 	}
//
// 	if len(buffer) >= 10 {
// 	  if (buffer[9] & 0b10000000) > 0 { faults = append(faults, "Battery (high voltage)") }
// 	  if (buffer[9] & 0b10000) > 0 { faults = append(faults, "Lambda 1 bank 1 (high voltage)") }
// 	  if (buffer[9] & 0b100) > 0 { faults = append(faults, "Throttle pot (high voltage)") }
// 	  if (buffer[9] & 0b10) > 0 { faults = append(faults, "Intake air temp (high voltage)") }
// 	  if (buffer[9] & 0b1) > 0 { faults = append(faults, "MAP sensor (high voltage)") }
// 	}
//
// 	if len(buffer) >= 13 {
// 	  if (buffer[12] & 0b1000000) > 0 { faults = append(faults, "Outside temp sensor (present)") }
// 	  if (buffer[12] & 0b100000) > 0 { faults = append(faults, "Power supply (present)") }
// 	  if (buffer[12] & 0b10000) > 0 { faults = append(faults, "Oil temp (present)") }
// 	  if (buffer[12] & 0b100) > 0 { faults = append(faults, "Coolant temp (present)") }
// 	  if (buffer[12] & 0b100) > 0 { faults = append(faults, "System voltage (present)") }
// 	}
//
// 	if len(buffer) >= 14 {
// 	  if (buffer[13] & 0b10000000) > 0 { faults = append(faults, "Battery voltage (present)") }
// 	  if (buffer[13] & 0b10000) > 0 { faults = append(faults, "Lambda 1 bank 1 (present)") }
// 	  if (buffer[13] & 0b100) > 0 { faults = append(faults, "Throttle pot (present)") }
// 	  if (buffer[13] & 0b10) > 0 { faults = append(faults, "Intake air temp (present)") }
// 	  if (buffer[13] & 0b1) > 0 { faults = append(faults, "MAP sensor (present)") }
// 	}
//
// 	if len(buffer) >= 24 {
// 	  if (buffer[23] & 0b1000) > 0 { faults = append(faults, "MAP sensor (present 2)") }
// 	  if (buffer[23] & 0b100) > 0 { faults = append(faults, "Oil temp (present 2)") }
// 	  if (buffer[23] & 0b10) > 0 { faults = append(faults, "Intake air temp (present 2)") }
// 	  if (buffer[23] & 0b1) > 0 { faults = append(faults, "Coolant temp (present 2)") }
// 	}
//
// 	if len(buffer) >= 26 {
// 	  if (buffer[25] & 0b1000) > 0 { faults = append(faults, "MAP sensor (historic)") }
// 	  if (buffer[25] & 0b100) > 0 { faults = append(faults, "Oil temp (historic)") }
// 	  if (buffer[25] & 0b10) > 0 { faults = append(faults, "Intake air temp (historic)") }
// 	  if (buffer[25] & 0b1) > 0 { faults = append(faults, "Coolant temp (historic)") }
// 	}
//
// 	if len(buffer) >= 27 {
// 	  if (buffer[26] & 0b00000001) > 0 { faults = append(faults, "Road speed sensor (present)") }
// 	  if (buffer[26] & 0b00000010) > 0 { faults = append(faults, "Comm. with AT (present)") }
// 	  if (buffer[26] & 0b00010000) > 0 { faults = append(faults, "Feedback (present)") }
// 	}
//
// 	if len(buffer) >= 28 {
// 	  if (buffer[27] & 0b00000001) > 0 { faults = append(faults, "Road speed sensor (historic)") }
// 	  if (buffer[27] & 0b00000010) > 0 { faults = append(faults, "Comm. with AT (historic)") }
// 	  if (buffer[27] & 0b00010000) > 0 { faults = append(faults, "Feedback (historic)") }
// 	}
//
// 	globalFaults = faults
// }
