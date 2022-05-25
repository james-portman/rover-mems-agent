package main

import (
	"fmt"
	"time"
	"errors"
  "encoding/hex"
	"github.com/distributed/sers"
)

var (
  ecu1xGotKlineEcho = false
  ecu1xLastKlineByte = byte(0x00)

	ecu1xRequestClearFaults = byte(0xCC)
	ecu1xStartTestRpmGauge = byte(0x6B)
	ecu1xStartTestLambdaHeater = byte(0x19)
	ecu1xStopTestLambdaHeater = byte(0x09)
	ecu1xStartTestACClutch = byte(0x13)
	ecu1xStopTestACClutch = byte(0x03)
	ecu1xStartTestFuelPump = byte(0x11)
	ecu1xStopTestFuelPump = byte(0x01)
	ecu1xStartTestFan1 = byte(0x1D)
	ecu1xStopTestFan1 = byte(0x0D)
	ecu1xStartTestPurgeValve = byte(0x18)
	ecu1xStopTestPurgeValve = byte(0x08)
	ecu1xIncreaseIdleSpeed = byte(0x91)
	ecu1xDecreaseIdleSpeed = byte(0x92)
	ecu1xIncreaseIgnitionAdvanceOffset = byte(0x93)
	ecu1xDecreaseIgnitionAdvanceOffset = byte(0x94)

	ecu1xUserCommands = map[string] byte{
		"clearfaults": ecu1xRequestClearFaults,
		"startTestRpmGauge": ecu1xStartTestRpmGauge,
		"startTestLambdaHeater": ecu1xStartTestLambdaHeater,
		"stopTestLambdaHeater": ecu1xStopTestLambdaHeater,
		"startTestACClutch": ecu1xStartTestACClutch,
		"stopTestACClutch": ecu1xStopTestACClutch,
		"startTestFuelPump": ecu1xStartTestFuelPump,
		"stopTestFuelPump": ecu1xStopTestFuelPump,
		"startTestFan1": ecu1xStartTestFan1,
		"stopTestFan1": ecu1xStopTestFan1,
		"startTestPurgeValve": ecu1xStartTestPurgeValve,
		"stopTestPurgeValve": ecu1xStopTestPurgeValve,
		"increaseIdleSpeed": ecu1xIncreaseIdleSpeed,
		"decreaseIdleSpeed": ecu1xDecreaseIdleSpeed,
		"increaseIgnitionAdvanceOffset": ecu1xIncreaseIgnitionAdvanceOffset,
		"decreaseIgnitionAdvanceOffset": ecu1xDecreaseIgnitionAdvanceOffset,
	}
)

func ecu1xNextCommand(previousResponse byte) byte {
	if globalUserCommand != "" {
		command, ok := ecu1xUserCommands[globalUserCommand];
		if ok {
			globalUserCommand = ""
			fmt.Println("> "+globalUserCommand)
			return command
		} else {
			fmt.Println("Asked to perform a user command but don't understand it")
		}
	}

	switch previousResponse {
		// go back to data 80 after clearing faults
		case ecu1xRequestClearFaults: return 0x80; break

		// init sequence then data 80
		case 0xCA: return 0x75; break
		case 0x75: return 0xF4; break
		case 0xF4: return 0xD0; break
		case 0xD0: return 0x80; break

		// toggle between data packets (1.2 ECU can only do 80 I think?)
		case 0x80: return 0x7D; break
		case 0x7D: return 0x80; break
	}

	return 0x80; // data 80 if we aren't sure
}

func ecu1xSend(sp sers.SerialPort, data byte) {
  sp.Write([]byte{data})
  ecu1xGotKlineEcho = false
  ecu1xLastKlineByte = data
}

func ecu1xLoop(sp sers.SerialPort, kline bool) ([]byte, error) {
  // start of init
  ecu1xSend(sp, 0xCA)

	buffer := make([]byte, 0)

	readLoops := 0
	readLoopsLimit := 200

	READLOOP:
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

    if kline && !ecu1xGotKlineEcho {
      if buffer[0] == ecu1xLastKlineByte {
        ecu1xGotKlineEcho = true
        // fmt.Println("got our kline echo")
        buffer = buffer[1:]
        continue
      }
    }

    if len(buffer) == 0 { continue }

		// check through the user commands (if we got another byte back as well)
		if len(buffer) >= 2 {
			for key := range ecu1xUserCommands {
				if buffer[0] == ecu1xUserCommands[key] {

					fmt.Println("< "+key)
					globalAlert = "ECU accepted "+key
					ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
					buffer = nil
					continue READLOOP // need to jump out twice
				}
			}
		}


		switch buffer[0] {
			// case ecu1xRequestClearFaults:
			// 	if len(buffer) >= 2 && buffer[1] == 0x00 {
			// 		fmt.Println("< FAULTS CLEARED")
			// 		globalAlert = "ECU reports faults cleared"
			// 		ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
			// 		buffer = nil
			// 		continue
			// 	}
			// 	break

			case 0xCA:
	      fmt.Println("Got CA")
				ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
				buffer = nil
	      continue
				break

			case 0x75:
	      fmt.Println("Got 75")
	      ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
				buffer = nil
	      continue
				break

	    case 0xF4:
	      if len(buffer) >= 2 && buffer[1] == 0x00 {
	        fmt.Println("Got F4 00")
	        ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
					buffer = nil
	        continue
	      }
				break


	    case 0xD0:
	      if len(buffer) >= 5 {
	        globalConnected = true
	        fmt.Println("Got D0 and ECU ID")
	        fmt.Printf("ECU ID:\n%s", hex.Dump(buffer[1:5]))
	        ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
					buffer = nil
	        continue
	      }
				break


	    case 0x80:
	      if len(buffer) >= 2 {
	        fullLength := int(buffer[1]) + 1
	        if len(buffer) >= fullLength {
	          fmt.Println("Got data 80")
	          ecu1xParseData80(buffer)
	          ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
						buffer = nil
	        }
	      }
	      // not got the full packet yet
	      continue
				break

	    case 0x7D:
	      if len(buffer) >= 2 {
	        fullLength := int(buffer[1]) + 1
	        if len(buffer) >= fullLength {
	          fmt.Println("Got data 7D")
	          ecu1xParseData7D(buffer)
	          ecu1xSend(sp, ecu1xNextCommand(buffer[0]))
						buffer = nil
	        }
	      }
	      // not got the full packet yet
	      continue
				break


		} // end switch

		// unknown command?
		// could be one of the normal commands waiting for their 2nd byte so don't do anything here

  }
  if readLoops >= readLoopsLimit {
		fmt.Printf("had buffer data: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		return nil, errors.New("MEMS 1.x timed out")
	}
	fmt.Println("fell out of readloop")

	return nil, nil

}

func ecu1xParseData80(data []byte) {
	globalDataOutputLock.Lock()
	defer globalDataOutputLock.Unlock()

  faults := []string {}
  // fmt.Printf("data 80 %d bytes \n%s", len(data), hex.Dump(data))

  // data[0] is the command (0x80)
  data = data[1:]

  packet_size := int(data[0])
  // 14 bytes length for mems 1.3
  // ? bytes length for mems 1.6

  // byte 1-2(16 bit) - engine speed in RPM
  globalDataOutput["rpm"] = float32( (int(data[1]) << 8) + int(data[2]) );

  // // byte 3 - coolant temp (+55 offset and 8 bit wrap)
  globalDataOutput["coolant_temp"] = float32(data[3]) - 55
  // TODO: fault if 59

  // //TODO:  byte 4 - (computed) ambient temp (+55 offset, 8 bit wrap) - doesn't work, might on MPI?
  globalDataOutput["ambient_temp"] = float32(data[4]) - 55
  // TODO: fault if 200

  // // byte 5 - intake air temp (+55 offset, 8 bit wrap)
  globalDataOutput["intake_air_temp"] = float32(data[5]) - 55
  // TODO: fault if 35

  // // byte 6 - fuel temp - doesn't work on SPI, do for MPI? # defaults to FF
  globalDataOutput["fuel_rail_temp"] = float32(data[6]) - 55

  // // byte 7 - map sensor kpa
  globalDataOutput["map_sensor_kpa"] = float32(data[7])

  // // byte 8 - battery voltage
  globalDataOutput["battery_voltage"] = float32(data[8])/10

  // // byte 9 - throttle pot voltage, WOT should be about 5v, should closed be near 0v? 0.02V per LSB. WOT should probably be close to 0xFA or 5.0V.
  globalDataOutput["throttle_pot_voltage"] = float32(data[9])/200
  //
  // // byte 10(A) - idle switch (bit 4 set if throttle closed)
  globalDataOutput["idle_switch"] = float32( (int(data[10]) & 0x00001000) >> 3 );
  //
  // // byte 11(B) - unknown, Probably a bitfield. Observed as 0x24 with engine off, and 0x20 with engine running. A single sample during a fifteen minute test drive showed a value of 0x30.
  //
  // // byte 12(C) - Park/neutral switch. Zero is closed, nonzero is open
  globalDataOutput["park_or_neutral_switch"] = float32(data[12])
  //
  // byte 13(D) - faults on mini spi
  // output['fault_1_bits'] = data[13].toString(2);
  // while (output['fault_1_bits'].length < 8) {
  //   output['fault_1_bits'] = "0"+output['fault_1_bits'];
  // }
  //  output['fault_1_bits'] = {"name": "Fault byte 1 bits", "data": output['fault_1_bits']};
  //
  if ((int(data[13]) >> 0) & 1) > 0 { faults = append(faults, "fault_coolant_temp_sensor") }
  if ((int(data[13]) >> 1) & 1) > 0 { faults = append(faults, "fault_inlet_air_temp_sensor") }
  if ((int(data[13]) >> 3) & 1) > 0 { faults = append(faults, "fault_turbo_overboost") }
  if ((int(data[13]) >> 4) & 1) > 0 { faults = append(faults, "fault_ambient_temp_sensor") }
  if ((int(data[13]) >> 5) & 1) > 0 { faults = append(faults, "fault_fuel_rail_temp_sensor") }
  if ((int(data[13]) >> 6) & 1) > 0 { faults = append(faults, "fault_knock_detected") }

  // // byte 14(E) - fault codes
  // output['fault_2_bits'] = data[14].toString(2);
  // while (output['fault_2_bits'].length < 8) {
  //   output['fault_2_bits'] = "0"+output['fault_2_bits'];
  // }
  // output['fault_2_bits'] = {"name": "Fault byte 2 bits", "data": output['fault_2_bits']};
  //
  if ((int(data[14]) >> 0) & 1 ) > 0 { faults = append(faults, "fault_coolant_temp_gauge") }
  if ((int(data[14]) >> 1) & 1 ) > 0 { faults = append(faults, "fault_fuel_pump_circuit") }
  if ((int(data[14]) >> 3) & 1 ) > 0 { faults = append(faults, "fault_air_con_clutch") }
  if ((int(data[14]) >> 4) & 1 ) > 0 { faults = append(faults, "fault_purge_valve") }
  if ((int(data[14]) >> 5) & 1 ) > 0 { faults = append(faults, "fault_map_sensor") }
  if ((int(data[14]) >> 6) & 1 ) > 0 { faults = append(faults, "fault_boost_valve") }
  if ((int(data[14]) >> 7) & 1 ) > 0 { faults = append(faults, "fault_throttle_pot_circuit") }

  // // 15(F) idle setting - x6.1
  if packet_size > 15 {
    globalDataOutput["idle_setpoint"] = float32(data[15]) * 6.1
  }
  //
  // // 16 (10) unknown
  if packet_size > 16 {
    globalDataOutput["idle_hotdb"] = float32(data[16])
  }
  //
  // // 17 (11) unknown
  //
  // // 18 (x12) - idle air control motor position - 0 closed, 180 fully open
  if packet_size > 0x12 {
    globalDataOutput["idle_valve_position"] = float32(data[0x12])
  }
  //
  // // 19-20 (x13-14) - idle speed deviation (16 bits)
  if packet_size > 0x14 {
    idle_deviation := int(data[0x13]) << 8;
    idle_deviation += int(data[0x14]);
    globalDataOutput["idle_speed_deviation"] = float32(idle_deviation)
  }
  //
  // // 21 (x15) unknown
  if (packet_size > 0x15) {
    globalDataOutput["ignition_advance_offset"] = float32(data[0x15])
  }
  //
  // // TODO: 22 (x16) - ignition advance 0.5 degrees per lsb, range -24 deg (00) to 103.5 deg (0xFF)
  if (packet_size > 0x16) {
  //   // advance /= 2;
  //   // advance -= 24;
    globalDataOutput["ignition_advance_raw"] = float32(data[0x16])
    globalDataOutput["ignition_advance"] = float32(data[0x16] / 2)
  }
  //
  // // TODO: 23-24 (x17-18) - coil time 0.002ms per lsb (16 bit)
  if (packet_size > 0x18) {
    coil_time := int(data[0x17]) << 8
    coil_time += int(data[0x18])
    globalDataOutput["coil_time_microseconds"] = float32(coil_time) * 2
  }
  // // 25 (x19) unknown
  // // 26 (x1a) unknown
  // // 27 (x1B) unknown


  globalFaults = faults
}



func ecu1xParseData7D(data []byte) {
	globalDataOutputLock.Lock()
	defer globalDataOutputLock.Unlock()
  // fmt.Printf("data 7D %d bytes \n%s", len(data), hex.Dump(data))

  // data[0] is the command (0x7D)
  data = data[1:]
  packet_size := int(data[0])

  globalDataOutput["ignition_switch"] = float32(data[1])
  globalDataOutput["throttle_angle"] = float32(data[2])/2
  // 0x03  Unknown
  globalDataOutput["air_fuel_ratio"] = float32(data[4]) // "A/F ratio? might just be 0xFF (unknown)" ## if it's FF then don't output?

  dtc_byte := int(data[5])
  globalDataOutput["lambda_heater_relay"] = float32((dtc_byte >> 3) & 1)
  globalDataOutput["secondary_trigger_sync"] = float32((dtc_byte >> 4) & 1)
  globalDataOutput["fan_1_control"] = float32((dtc_byte >> 5) & 1)
  globalDataOutput["fan_2_control"] = float32((dtc_byte >> 7) & 1)

  globalDataOutput["lambda_mv"] = float32(data[6]) * 5
  // TODO: error on 435/440?

  globalDataOutput["lambda_sensor_frequency"] = float32(data[7]);

  globalDataOutput["lambda_sensor_duty_cycle"] = float32(data[8])
  globalDataOutput["lambda_sensor_status"] = float32(data[9]) // "Lambda sensor status? 0x01 for good, any other value for no good"

  if int(data[10]) > 0 {
    globalDataOutput["closed_loop"] = 1 // "Loop indicator, 0 for open loop and nonzero for closed loop"
  } else {
    globalDataOutput["closed_loop"] = 0
  }

  globalDataOutput["long_term_trim"] = float32(data[11])
  globalDataOutput["short_term_trim_percent"] = float32(data[12])
  globalDataOutput["carbon_can_purge_valve_duty_cycle"] = float32(data[13]) // "Carbon canister purge valve duty cycle?"

  dtc2 := int(data[0xE])
  globalDataOutput["primary_trigger_sync"] = float32( (dtc2 >> 1) & 1 )

  if (packet_size >= 16) {
    globalDataOutput["idle_base_position"] = float32(data[15])
  }

  // 0x10  Unknown
  // 0x11  Unknown
  // 0x12  Unknown
  // 0x13  Unknown
  if (packet_size >= 21) {
    globalDataOutput["idle_error"] = float32(data[20])
  }
  // 0x15  Unknown

  if (packet_size >= 0x16) {
    dtc3 := int(data[0x16]);
    globalDataOutput["injector_1_4_driver"] = float32((dtc3 >> 1) & 1)
    globalDataOutput["injector_2_3_driver"] = float32((dtc3 >> 2) & 1)
    globalDataOutput["fault_engine_bay_vent_warning"] = float32((dtc3 >> 3) & 1)
    globalDataOutput["engine_bay_vent_relay"] = float32((dtc3 >> 4) & 1)
    globalDataOutput["hill_assist"] = float32((dtc3 >> 5) & 1)
    globalDataOutput["cruise_control"] = float32((dtc3 >> 6) & 1)
  }

  // 0x17  Unknown
  // 0x18  Unknown
  // 0x19  Unknown
  // 0x1A  Unknown
  // 0x1B  Unknown
  // 0x1C  Unknown
  // 0x1D  Unknown
  // 0x1E  Unknown
  if (packet_size >= 0x1F) {
    globalDataOutput["crank_counter"] = float32(data[0x1F])
  }

}
