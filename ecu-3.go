package main

import (
	"fmt"
	"time"
	"errors"
	"encoding/hex"
	"github.com/distributed/sers"
)

var (
	ecu3RequestHeader = []byte {0xB8, 0x13, 0xF7}

	ecu3InitCommand = []byte {0x1A, 0x9A}
	ecu3InitAccepted = []byte {0x5A, 0x9A}

	ecu3StartDiagnostic = []byte {0x10, 0xA0}
	ecu3StartDiagResponse = []byte {0x50}

	ecu3RequestSeed = []byte {0x27, 0x01}
	ecu3SeedResponse = []byte {0x67, 0x01}
	ecu3Seed = 0

	ecu3SendKey = []byte {0x27, 0x02}
	ecu3KeyAcceptResponse = []byte {0x67, 0x02}
	ecu3Key = 0

	ecu3PingCommand = []byte {0x3E}
	ecu3PongResponse = []byte {0x7E}

	ecu3ClearFaultsCommand = []byte {0x14, 0x00, 0x00}
	ecu3FaultsClearedResponse = []byte {0x54, 0x00, 0x00}

	ecu3RequestFaultsCommand = []byte {0x18, 0x0, 0x0, 0x0}
	ecu3ResponseFaults = []byte {0x58}

	ecu3RequestData00 = []byte {0x21, 0x00}
	ecu3RequestData06 = []byte {0x21, 0x06}
	ecu3RequestData0A = []byte {0x21, 0x0A}
	ecu3RequestData0B = []byte {0x21, 0x0B}
	ecu3RequestData21 = []byte {0x21, 0x21}

	ecu3ResponseData00 = []byte {0x61, 0x00}
	ecu3ResponseData06 = []byte {0x61, 0x06}
	ecu3ResponseData0A = []byte {0x61, 0x0A}
	ecu3ResponseData0B = []byte {0x61, 0x0B}
	ecu3ResponseData21 = []byte {0x61, 0x21}

	ecu3UserCommands = map[string] []byte{
		"clearfaults": ecu3ClearFaultsCommand,
	}


  ecu3FaultTypes = map[int]string {
    0x20: "historical",
    0x74: "present, test not complete",
    0x30: "historical, test not complete",
    0x58: "present, test not complete",
    0x61: "present",
    0x62: "present",
    0x64: "present",
    0x71: "present, test not complete",
  }
  ecu3Faults = map[int]string {
    0x1232: "fuel pump relay, open circuit",
    0x0650: "MIL control circuit malfunction",
    0x0481: "A/C condensor fan",
    0x1508: "IACV driver open circuit",
    0x1186: "front lambda heater",
    0x1185: "front lambda heater",
    0x1192: "rear lambda heater",
    0x0445: "purge valve drive",
    0x0480: "cooling fan",
    0x1610: "main relay - open circuit",
    0x0113: "IAT shorted",
    0x0118: "coolant temp sensor shorted",
    0x0122: "throttle pot shorted",
    0x0562: "system voltage malfunction",
    0x0197: "oil temp sensor shorted",
    0x0462: "fuel tank level sensor shorted to ground",
    0x0340: "cam position sensor",
    0x0106: "manifold pressure - incorrect reading",
    0x1316: "misfire causing excess emissions",
    0x0170: "fuel system",
    0x0655: "warning lamp - engine bay temperature - open circuit",
  }

)


func ecu3SendNextCommand(sp sers.SerialPort, previousResponse []byte) {
	if globalUserCommand != "" {
		command, ok := ecu3UserCommands[globalUserCommand];
		if ok {
			globalUserCommand = ""
			ecu3SendCommand(sp, command)
			return
		} else {
			fmt.Println("Asked to perform a user command but don't understand it")
		}
	}

	if slicesEqual(previousResponse, ecu3InitAccepted) {
		ecu3SendCommand(sp, ecu3StartDiagnostic)

	} else if slicesEqual(previousResponse, ecu3StartDiagResponse) {
		ecu3SendCommand(sp, ecu3RequestSeed)

  } else if slicesEqual(previousResponse, ecu3SeedResponse) {
		command := append(ecu3SendKey, byte(ecu3Key >> 8))
		command = append(command, byte(ecu3Key & 0xFF))
		ecu3SendCommand(sp, command)

	} else if slicesEqual(previousResponse, ecu3KeyAcceptResponse) {
		ecu3SendCommand(sp, ecu3PingCommand)

   } else if slicesEqual(previousResponse, ecu3PongResponse) {
		ecu3SendCommand(sp, ecu3RequestFaultsCommand)

	} else if slicesEqual(previousResponse, ecu3ResponseFaults) { ecu3SendCommand(sp, ecu3RequestData00)
	} else if slicesEqual(previousResponse, ecu3ResponseData00) { ecu3SendCommand(sp, ecu3RequestData06)
	} else if slicesEqual(previousResponse, ecu3ResponseData06) { ecu3SendCommand(sp, ecu3RequestData0A)
	} else if slicesEqual(previousResponse, ecu3ResponseData0A) { ecu3SendCommand(sp, ecu3RequestData0B)
	} else if slicesEqual(previousResponse, ecu3ResponseData0B) { ecu3SendCommand(sp, ecu3RequestData21)
	} else if slicesEqual(previousResponse, ecu3ResponseData21) { ecu3SendCommand(sp, ecu3PingCommand)

	} else if slicesEqual(previousResponse, ecu3FaultsClearedResponse) {
		sp.Write(ecu3RequestFaultsCommand)
		globalAlert = "ECU reports faults cleared"

	} else { // fall back to ping
		ecu3SendCommand(sp, ecu3PingCommand)
	}
//
}

func ecu3SendCommand(sp sers.SerialPort, data []byte) {
  output := ecu3RequestHeader // []byte {0xB8, 0x13, 0xF7} // always
	output = append(output, byte(len(data)))
	output = append(output, data...)
  output = append(output, xor_all_bytes(output))
	sp.Write(output)
	// fmt.Println("> command sent")
}


func readFirstBytesFromPortEcu3(fn string) ([]byte, error) {

	fmt.Println("Connecting to MEMS 3 ECU")
	globalConnected = false

	sp, err := sers.Open(fn)
	if err != nil {
		return nil, err
	}
	defer sp.Close()

	err = sp.SetMode(9600, 8, sers.E, 1, sers.NO_HANDSHAKE)
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

	ecu3SendCommand(sp, ecu3InitCommand)

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

		if len(buffer) == 0 { continue }

		// check for init echo
		if len(buffer) >= 2 && slicesEqual(buffer[0:2], ecu3InitCommand) {
			fmt.Println("Got our init echo")
			buffer = buffer[2:]
			continue
		}

		// check for full commands - our echos and responses too
		// [always 3 byte header] [data length] [data...] [checksum]

		if len(buffer) < 4 { continue } // no length byte yet
		dataLength := int(buffer[3])
		totalLength := 3 + 1 + dataLength + 1
		if len(buffer) < totalLength { continue } // have length byte but not whole packet yet

		// TODO: check checksum ?

		// thisPacket := buffer[0:totalLength]
		actualData := buffer[4:4+dataLength]

		// our echos
		if slicesEqual(buffer[0:3], ecu3RequestHeader) {
			// fmt.Println("This looks like our own echo")

			if slicesEqual(actualData, ecu3InitCommand) {
				// fmt.Println("Got our init echo")
				buffer = buffer[totalLength:]
				continue
			}
			if slicesEqual(actualData, ecu3StartDiagnostic) {
				// fmt.Println("Got our start diag echo")
				buffer = buffer[totalLength:]
				continue
			}
			if slicesEqual(actualData, ecu3RequestSeed) {
				// fmt.Println("Got our seed req echo")
				buffer = buffer[totalLength:]
				continue
			}
			if slicesEqual(actualData[0:2], ecu3SendKey) {
				// fmt.Println("Got our key send echo")
				buffer = buffer[totalLength:]
				// buffer = buffer[(len(ecu3SendKey)+2+2):] // extra 2 for key
				continue
			}
			if slicesEqual(actualData, ecu3PingCommand) {
				// fmt.Println("Got our ping echo")
				buffer = buffer[totalLength:]
				continue
			}

			// echo not caught if we end up here
			// don't really care, if it's obviously our echo then it's fine
			// fmt.Println("*** Unknown echo caught here")
			// fmt.Printf("buffer: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
			buffer = buffer[totalLength:]
			continue

		} // end of our echos

		//
		// actual responses
		//
		if slicesEqual(actualData[0:2], ecu3InitAccepted) {
			fmt.Println("< ECU woke up")
			buffer = nil
			globalConnected = true
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3InitAccepted)
			continue
		}
		if slicesEqual(actualData, ecu3StartDiagResponse) {
			fmt.Println("< Diag mode accepted")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3StartDiagResponse)
			continue
		}
		if slicesEqual(actualData[0:2], ecu3SeedResponse) {
			fmt.Println("< seed")
			ecu3Seed = int(actualData[2]) << 8
			ecu3Seed += int(actualData[3])
			fmt.Println(ecu3Seed)
			if ecu3Seed == 0 { // auth not required/already done
				ecu3Key = 0
				buffer = nil
				time.Sleep(50 * time.Millisecond)
				ecu3SendNextCommand(sp, nil)
				continue
			} else {
				// do key generation
				ecu3Key = generateKey(ecu3Seed)
				buffer = nil
				time.Sleep(50 * time.Millisecond)
				ecu3SendNextCommand(sp, ecu3SeedResponse)
				continue
			}

		}
		if slicesEqual(actualData, ecu3KeyAcceptResponse) {
			fmt.Println("< Key accepted")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3KeyAcceptResponse)
			continue
		}
		if slicesEqual(actualData, ecu3PongResponse) {
			fmt.Println("< PONG")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3PongResponse)
			continue
		}
		if slicesEqual(actualData, ecu3FaultsClearedResponse) {
			fmt.Println("< FAULT CLEARED")
			globalAlert = "ECU reports faults cleared"
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3FaultsClearedResponse)
			continue
		}

		if slicesEqual(actualData[0:len(ecu3ResponseFaults)], ecu3ResponseFaults) {
			fmt.Println("< Faults")
			ecu3ParseFaults(actualData)
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3ResponseFaults)
			continue
		}

		if slicesEqual(actualData[0:2], ecu3ResponseData00) {
			fmt.Println("got data packet 00")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3ResponseData00)
			continue
		}
		if slicesEqual(actualData[0:2], ecu3ResponseData06) {
			fmt.Println("got data packet 06")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3ResponseData06)
			continue
		}
		if slicesEqual(actualData[0:2], ecu3ResponseData0A) {
			fmt.Println("got data packet 0A")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3ResponseData0A)
			continue
		}
		if slicesEqual(actualData[0:2], ecu3ResponseData0B) {
			fmt.Println("got data packet 0B")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3ResponseData0B)
			continue
		}
		if slicesEqual(actualData[0:2], ecu3ResponseData21) {
			fmt.Println("got data packet 21")
			buffer = nil
			time.Sleep(50 * time.Millisecond)
			ecu3SendNextCommand(sp, ecu3ResponseData21)
			continue
		}
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

		fmt.Printf("unknown command in buffer (burning it): got %d bytes \n%s", len(buffer), hex.Dump(buffer[0:totalLength]))
		fmt.Printf("actualData %d bytes \n%s", len(actualData), hex.Dump(actualData))
		buffer = buffer[totalLength:]

	}
	if readLoops >= readLoopsLimit {
		fmt.Printf("had buffer data: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		return nil, errors.New("readloop timed out")
	}
	fmt.Println("fell out of readloop")

	return nil, err
}

func ecu3ParseFaults(buffer []byte) {
	// fmt.Printf("ecu3ParseFaults got %d bytes \n%s", len(buffer), hex.Dump(buffer))

	faults := []string {}

  buffer = buffer[2:] // throw away command
  for len(buffer) >= 3 {
    thisfault := int(buffer[0]) << 8
    thisfault += int(buffer[1])

    faulttype := int(buffer[2])

		output_fault_type, ok := ecu3FaultTypes[int(faulttype)];
		if !ok {
			output_fault_type = "unknown ("+string(int(faulttype))+")"
		}

    output_fault, ok := ecu3Faults[int(thisfault)]
		if !ok {
			output_fault = "unknown ("+string(int(thisfault))+")"
		}

    full_output_text := "Fault - "+output_fault+" - "+output_fault_type

		faults = append(faults, full_output_text)

    if len(buffer) > 3 {
      buffer = buffer[3:]
    } else {
      buffer = nil
    }
  }

	globalFaults = faults
}
