package main

import (
	"fmt"
	"time"
	"errors"
	"encoding/hex"
	"github.com/distributed/sers"
)

var (
	twojInitCommand = []byte {0x81, 0x13, 0xF7, 0x81, 0x0C}
	twojStartDiagnostic = []byte {0x10, 0xA0}
	twojRequestSeed = []byte {0x27, 0x01}
	twojSendKey = []byte {0x27, 0x02}
	twojPingCommand = []byte {0x3E}
	twojClearFaultsCommand = []byte {0x31, 0xCB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	twojRequestFaultsCommand = []byte {0x21, 0x19}

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
	twojFaultsClearedResponse = []byte {0x71, 0xCB}

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

	twojUserCommands = map[string] []byte{
		"clearfaults": twojClearFaultsCommand,
	}

	twojFaults = map[int]string{
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

func twojSendCommand(sp sers.SerialPort, command []byte) {

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
	sp.Write(finalCommand)
}


func twojSendNextCommand(sp sers.SerialPort, previousResponse []byte) {
	if globalUserCommand != "" {
		command, ok := twojUserCommands[globalUserCommand];
		if ok {
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

	} else if slicesEqual(previousResponse, twojSeedResponse) {
		command := append(twojSendKey, byte(twojKey >> 8))
		command = append(command, byte(twojKey & 0xFF))
		twojSendCommand(sp, command)

	} else if slicesEqual(previousResponse, twojKeyAcceptResponse) {
		twojSendCommand(sp, twojPingCommand)

	} else if slicesEqual(previousResponse, twojPongResponse) {
		twojSendCommand(sp, twojRequestFaultsCommand)

	} else if slicesEqual(previousResponse, twojFaultsClearedResponse) {
		twojSendCommand(sp, twojRequestFaultsCommand)

	} else if slicesEqual(previousResponse, twojFaultsResponse) { twojSendCommand(sp, twojRequestData00)

	} else if slicesEqual(previousResponse, twojResponseData00) { twojSendCommand(sp, twojRequestData01)
	} else if slicesEqual(previousResponse, twojResponseData01) { twojSendCommand(sp, twojRequestData02)
	} else if slicesEqual(previousResponse, twojResponseData02) { twojSendCommand(sp, twojRequestData03)
	} else if slicesEqual(previousResponse, twojResponseData03) { twojSendCommand(sp, twojRequestData05)
	} else if slicesEqual(previousResponse, twojResponseData05) { twojSendCommand(sp, twojRequestData06)
	} else if slicesEqual(previousResponse, twojResponseData06) { twojSendCommand(sp, twojRequestData07)
	} else if slicesEqual(previousResponse, twojResponseData07) { twojSendCommand(sp, twojRequestData08)
	} else if slicesEqual(previousResponse, twojResponseData08) { twojSendCommand(sp, twojRequestData09)
	} else if slicesEqual(previousResponse, twojResponseData09) { twojSendCommand(sp, twojRequestData0A)
	} else if slicesEqual(previousResponse, twojResponseData0A) { twojSendCommand(sp, twojRequestData0B)
	} else if slicesEqual(previousResponse, twojResponseData0B) { twojSendCommand(sp, twojRequestData0C)
	} else if slicesEqual(previousResponse, twojResponseData0C) { twojSendCommand(sp, twojRequestData0D)
	} else if slicesEqual(previousResponse, twojResponseData0D) { twojSendCommand(sp, twojRequestData0F)
	} else if slicesEqual(previousResponse, twojResponseData0F) { twojSendCommand(sp, twojRequestData10)
	} else if slicesEqual(previousResponse, twojResponseData10) { twojSendCommand(sp, twojRequestData11)
	} else if slicesEqual(previousResponse, twojResponseData11) { twojSendCommand(sp, twojRequestData12)
	} else if slicesEqual(previousResponse, twojResponseData12) { twojSendCommand(sp, twojRequestData13)
	} else if slicesEqual(previousResponse, twojResponseData13) { twojSendCommand(sp, twojRequestData21)
	} else if slicesEqual(previousResponse, twojResponseData21) { twojSendCommand(sp, twojRequestData25)
	} else if slicesEqual(previousResponse, twojResponseData25) { twojSendCommand(sp, twojRequestData3A)
	} else if slicesEqual(previousResponse, twojResponseData3A) { twojSendCommand(sp, twojPingCommand)

	// } else if slicesEqual(previousResponse, twojWokeResponse) || slicesEqual(previousResponse, twojFaultsResponse) {
	// 	sp.Write(twojPingCommand)
	//
	// } else if slicesEqual(previousResponse, twojFaultsClearedResponse) {
	// 	sp.Write(twojRequestFaultsCommand)
	// 	globalAlert = "ECU reports faults cleared"

	} else { // fall back to ping
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
	time.Sleep(25 * time.Millisecond)
	sp.SetBreak(false)
	time.Sleep(25 * time.Millisecond)

	sp.Write(twojInitCommand)


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

		// clear leading zeros (from our wake up)
		for len(buffer) > 0 && buffer[0] == 0x00 {
			buffer = buffer[1:]
		}

		if len(buffer) == 0 { continue }

		// check for init echo
		if len(buffer) >= 5 && slicesEqual(buffer[0:5], twojInitCommand) {
			// fmt.Println("Got our init echo")
			buffer = buffer[5:]
			continue
		}

		// check for full commands - our echos and responses too

		if len(buffer) < int(buffer[0]) + 2 {
			// fmt.Println("waiting for rest of data packet")
			continue
		}

		// TODO: check checksum ?

		actualData := buffer[1:int(buffer[0])+1]
		// fmt.Printf("actual data: got %d bytes \n%s", len(actualData), hex.Dump(actualData))

		// our echos
		if slicesEqual(actualData, twojPingCommand) {
			// fmt.Println("Got our ping echo")
			buffer = buffer[(len(twojPingCommand)+2):]
			continue
		}
		if slicesEqual(actualData, twojStartDiagnostic) {
			// fmt.Println("Got our start diag echo")
			buffer = buffer[(len(twojStartDiagnostic)+2):]
			continue
		}
		if slicesEqual(actualData, twojRequestSeed) {
			// fmt.Println("Got our seed req echo")
			buffer = buffer[(len(twojRequestSeed)+2):]
			continue
		}
		if slicesEqual(actualData[0:2], twojSendKey) {
			// fmt.Println("Got our key send echo")
			buffer = buffer[(len(twojSendKey)+2+2):] // extra 2 for key
			continue
		}
		if slicesEqual(actualData, twojClearFaultsCommand) {
			// fmt.Println("Got our clear faults echo")
			buffer = buffer[(len(twojClearFaultsCommand)+2):]
			continue
		}
		if slicesEqual(actualData, twojRequestFaultsCommand) {
			// fmt.Println("Got our request faults echo")
			buffer = buffer[(len(twojRequestFaultsCommand)+2):]
			continue
		}
		if slicesEqual(actualData, twojRequestData00) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData01) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData02) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData03) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData05) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData06) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData07) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData08) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData09) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData0A) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData0B) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData0C) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData0D) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData0F) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData10) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData11) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData12) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData13) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData21) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData25) { buffer = buffer[4:]; continue }
		if slicesEqual(actualData, twojRequestData3A) { buffer = buffer[4:]; continue }


		// todo - all the point values are being lost

		// actual responses

		if slicesEqual(actualData, twojWokeResponse) {
			fmt.Println("< ECU woke up")
			buffer = nil
			globalConnected = true
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojWokeResponse)
			continue
		}
		if slicesEqual(actualData, twojStartDiagResponse) {
			fmt.Println("< Diag mode accepted")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojStartDiagResponse)
			continue
		}
		if slicesEqual(actualData[0:2], twojSeedResponse) {
			fmt.Println("< seed")
			twojSeed = int(actualData[2]) << 8
			twojSeed += int(actualData[3])
			// do key generation
			twojKey = generateKey(twojSeed)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojSeedResponse)
			continue
		}
		if slicesEqual(actualData, twojKeyAcceptResponse) {
			fmt.Println("< Key accepted")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojKeyAcceptResponse)
			continue
		}
		if slicesEqual(actualData, twojPongResponse) {
			fmt.Println("< PONG")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojPongResponse)
			continue
		}
		if slicesEqual(actualData, twojFaultsClearedResponse) {
			fmt.Println("< FAULT CLEARED")
			globalAlert = "ECU reports faults cleared"
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojFaultsClearedResponse)
			continue
		}

		if slicesEqual(actualData[0:2], twojFaultsResponse) {
			fmt.Println("< Faults")
			twojParseFaults(actualData)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojFaultsResponse)
			continue
		}


		if slicesEqual(actualData[0:2], twojResponseData00) {
			fmt.Println("got data packet 00")
			// don't care?
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData00)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData01) {
			fmt.Println("got data packet 01")
			coolant := int(actualData[2]) << 8
			coolant += int(actualData[3])
			coolantFloat := float32(coolant) - 2732
			coolantFloat /= 10
			globalDataOutput["coolant_temp"] = coolantFloat
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData01)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData02) {
			fmt.Println("got data packet 02")
			oiltemp := int(actualData[2]) << 8
			oiltemp += int(actualData[3])
			oiltempFloat := float32(oiltemp) - 2732
			oiltempFloat /= 10
			globalDataOutput["oil_temp"] = oiltempFloat
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData02)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData03) {
			fmt.Println("got data packet 03")
			iat := int(actualData[2]) << 8
			iat += int(actualData[3])
			iatFloat := float32(iat) - 2732
			iatFloat /= 10
			globalDataOutput["intake_air_temp"] = iatFloat
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData03)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData05) {
			fmt.Println("got data packet 05")
			fueltemp := int(actualData[2]) << 8
			fueltemp += int(actualData[3])
			globalDataOutput["fuel_temp"] = float32(fueltemp)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData05)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData06) {
			fmt.Println("got data packet 06")
			// don't care?
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData06)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData07) {
			fmt.Println("got data packet 07")
			mapkpa := int(actualData[2]) << 8
			mapkpa += int(actualData[3])
			globalDataOutput["map_sensor_kpa"] = float32(mapkpa)/100
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData07)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData08) {
			fmt.Println("got data packet 08")
			tps := int(actualData[2]) << 8
			tps += int(actualData[3])
			tpsFloat := float32(tps) / 100
			globalDataOutput["tps_degrees"] = tpsFloat
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData08)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData09) {
			fmt.Println("got data packet 09")
			rpm := int(actualData[2]) << 8
			rpm += int(actualData[3])
			globalDataOutput["rpm"] = float32(rpm)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData09)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData0A) {
			fmt.Println("got data packet 0A")
			feedback := int(actualData[2]) << 8
			feedback += int(actualData[3])
			feedbackFloat := float32(feedback) / 100
			globalDataOutput["fuelling_feedback_percent"] = feedbackFloat

			o2mv := int(actualData[4]) << 8
			o2mv += int(actualData[5])
			globalDataOutput["o2_mv"] = float32(o2mv)

			airFuel := ((float32(o2mv) / 1000) * 2) + 10
			globalDataOutput["estimate_air_fuel"] = airFuel

			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData0A)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData0B) {
			fmt.Println("got data packet 0B")
			globalDataOutput["coil_1_charge_time"] = float32(actualData[2]) / 1000
			globalDataOutput["coil_2_charge_time"] = float32(actualData[3]) / 1000
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData0B)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData0C) {
			fmt.Println("got data packet 0C")
			globalDataOutput["injector_1_pw"] = float32(actualData[2])
			globalDataOutput["injector_2_pw"] = float32(actualData[3])
			globalDataOutput["injector_3_pw"] = float32(actualData[4])
			globalDataOutput["injector_4_pw"] = float32(actualData[5])
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData0C)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData0D) {
			fmt.Println("got data packet 0D")
			globalDataOutput["vehicle_speed"] = float32(actualData[2])
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData0D)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData0F) {
			fmt.Println("got data packet 0F")
			globalDataOutput["throttle_switch"] = float32(int(actualData[2]) & 1) // 0b00000001
			globalDataOutput["ignition"] = float32((int(actualData[2]) >> 1) & 1) // 0b00000010
			globalDataOutput["ac_button"] = float32((int(actualData[2]) >> 3) & 1) // 0b00001000
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData0F)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData10) {
			fmt.Println("got data packet 10")
			battery := int(actualData[4]) << 8
			battery += int(actualData[5])
			batteryFloat := float32(battery) / 1000
			globalDataOutput["battery_voltage"] = batteryFloat
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData10)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData11) {
			fmt.Println("got data packet 11")
			// 0 means OK, 1 bad
			// will swap for our purposes
			// output is 1 for yes
			camSync := actualData[2] & 1 // 0b00000001
			globalDataOutput["cam_sync"] = float32(1 - camSync)
			crankSync := (actualData[2] >> 1) & 1 //0b00000010
			globalDataOutput["crank_sync"] = float32(1 - crankSync)

			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData11)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData12) {
			fmt.Println("got data packet 12")
			idleValvePos := int(actualData[2]) << 8
			idleValvePos += int(actualData[3])
			idleValveFloat := float32(idleValvePos) / 2
			globalDataOutput["idle_valve_pos"] = idleValveFloat
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData12)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData13) {
			fmt.Println("got data packet 13")
			globalDataOutput["closed_loop"] = float32(actualData[2] & 0b00000001)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData13)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData21) {
			fmt.Println("got data packet 21")
			rpmError := int(actualData[2]) << 8
			rpmError += int(actualData[3])
			if rpmError > 32768 {
				rpmError -= 65535
			}
			globalDataOutput["rpm_error"] = float32(rpmError)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData21)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData25) {
			fmt.Println("got data packet 25")
			camPercent := int(actualData[2]) << 8
			camPercent += int(actualData[3])
			globalDataOutput["cam_percent"] = float32(camPercent)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData25)
			continue
		}
		if slicesEqual(actualData[0:2], twojResponseData3A) {
			fmt.Println("got data packet 3A")
			idleTimingOffset := int(actualData[2]) << 8
			idleTimingOffset += int(actualData[3])
			idleTimingOffsetFloat := float32(idleTimingOffset) / 10
			globalDataOutput["idle_timing_offset"] = idleTimingOffsetFloat

			idleAdjusterRpm := int(actualData[4]) << 8
			idleAdjusterRpm += int(actualData[5])
			globalDataOutput["idle_adjuster_rpm"] = float32(idleAdjusterRpm)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			twojSendNextCommand(sp, twojResponseData3A)
			continue
		}

		// if we get here then something is wrong with the data

		// todo: cope with 7F (fail/no)

		fmt.Printf("unknown command in buffer (burning it): got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		buffer = nil
		twojSendNextCommand(sp, twojPongResponse)

	}
	if readLoops >= readLoopsLimit {
		fmt.Printf("had buffer data: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		return nil, errors.New("readloop timed out")
	}
	fmt.Println("fell out of readloop")

	return nil, err
}

func twojParseFaults(buffer []byte) {
	// fmt.Printf("got %d bytes \n%s", len(buffer), hex.Dump(buffer))

	faults := []string {}

	if len(buffer) >= 5 {
		if (buffer[4] & 0b01000000) > 0 {	faults = append(faults, "Outside air temp (low voltage)") }
	  if (buffer[4] & 0b00100000) > 0 { faults = append(faults, "Power supply (low voltage)") }
	  if (buffer[4] & 0b00010000) > 0 { faults = append(faults, "Engine oil temp (low voltage)") }
	  if (buffer[4] & 0b00000100) > 0 { faults = append(faults, "Coolant temp (low voltage)") }
	  if (buffer[4] & 0b00000001) > 0 { faults = append(faults, "System (low voltage)") }
	}


	if len(buffer) >= 6 {
	  if (buffer[5] & 0b10000000) > 0 { faults = append(faults, "Battery (low voltage)") }
	  if (buffer[5] & 0b00010000) > 0 { faults = append(faults, "Lambda 1 bank 1 (low voltage)") }
	  if (buffer[5] & 0b00000100) > 0 { faults = append(faults, "Throttle pot (low voltage)") }
	  if (buffer[5] & 0b00000010) > 0 { faults = append(faults, "Air intake (low voltage)") }
	  if (buffer[5] & 0b00000001) > 0 { faults = append(faults, "MAP sensor (low voltage)") }
	}

	if len(buffer) >= 9 {
	  if (buffer[8] & 0b01000000) > 0 { faults = append(faults, "Outside air temp (high voltage)") }
	  if (buffer[8] & 0b00100000) > 0 { faults = append(faults, "Power supply (high voltage)") }
	  if (buffer[8] & 0b00010000) > 0 { faults = append(faults, "Oil temperature (high voltage)") }
	  if (buffer[8] & 0b00000100) > 0 { faults = append(faults, "Coolant temperature (high voltage)") }
	  if (buffer[8] & 0b00000001) > 0 { faults = append(faults, "System (high voltage)") }
	}

	if len(buffer) >= 10 {
	  if (buffer[9] & 0b10000000) > 0 { faults = append(faults, "Battery (high voltage)") }
	  if (buffer[9] & 0b10000) > 0 { faults = append(faults, "Lambda 1 bank 1 (high voltage)") }
	  if (buffer[9] & 0b100) > 0 { faults = append(faults, "Throttle pot (high voltage)") }
	  if (buffer[9] & 0b10) > 0 { faults = append(faults, "Intake air temp (high voltage)") }
	  if (buffer[9] & 0b1) > 0 { faults = append(faults, "MAP sensor (high voltage)") }
	}

	if len(buffer) >= 13 {
	  if (buffer[12] & 0b1000000) > 0 { faults = append(faults, "Outside temp sensor (present)") }
	  if (buffer[12] & 0b100000) > 0 { faults = append(faults, "Power supply (present)") }
	  if (buffer[12] & 0b10000) > 0 { faults = append(faults, "Oil temp (present)") }
	  if (buffer[12] & 0b100) > 0 { faults = append(faults, "Coolant temp (present)") }
	  if (buffer[12] & 0b100) > 0 { faults = append(faults, "System voltage (present)") }
	}

	if len(buffer) >= 14 {
	  if (buffer[13] & 0b10000000) > 0 { faults = append(faults, "Battery voltage (present)") }
	  if (buffer[13] & 0b10000) > 0 { faults = append(faults, "Lambda 1 bank 1 (present)") }
	  if (buffer[13] & 0b100) > 0 { faults = append(faults, "Throttle pot (present)") }
	  if (buffer[13] & 0b10) > 0 { faults = append(faults, "Intake air temp (present)") }
	  if (buffer[13] & 0b1) > 0 { faults = append(faults, "MAP sensor (present)") }
	}

	if len(buffer) >= 24 {
	  if (buffer[23] & 0b1000) > 0 { faults = append(faults, "MAP sensor (present 2)") }
	  if (buffer[23] & 0b100) > 0 { faults = append(faults, "Oil temp (present 2)") }
	  if (buffer[23] & 0b10) > 0 { faults = append(faults, "Intake air temp (present 2)") }
	  if (buffer[23] & 0b1) > 0 { faults = append(faults, "Coolant temp (present 2)") }
	}

	if len(buffer) >= 26 {
	  if (buffer[25] & 0b1000) > 0 { faults = append(faults, "MAP sensor (historic)") }
	  if (buffer[25] & 0b100) > 0 { faults = append(faults, "Oil temp (historic)") }
	  if (buffer[25] & 0b10) > 0 { faults = append(faults, "Intake air temp (historic)") }
	  if (buffer[25] & 0b1) > 0 { faults = append(faults, "Coolant temp (historic)") }
	}

	if len(buffer) >= 27 {
	  if (buffer[26] & 0b00000001) > 0 { faults = append(faults, "Road speed sensor (present)") }
	  if (buffer[26] & 0b00000010) > 0 { faults = append(faults, "Comm. with AT (present)") }
	  if (buffer[26] & 0b00010000) > 0 { faults = append(faults, "Feedback (present)") }
	}

	if len(buffer) >= 28 {
	  if (buffer[27] & 0b00000001) > 0 { faults = append(faults, "Road speed sensor (historic)") }
	  if (buffer[27] & 0b00000010) > 0 { faults = append(faults, "Comm. with AT (historic)") }
	  if (buffer[27] & 0b00010000) > 0 { faults = append(faults, "Feedback (historic)") }
	}

	globalFaults = faults
}
