function parseDataBufferSlowInit() {
  switch (dataBuffer[0]) {
    case 0x00:
      debug("got a stray 0x00, ignoring");
      dataBuffer = dataBuffer.slice(1);
      break;

    case 0x55:
      if (dataBuffer.length < 3) {
        debug("got some of slow init initial reply but need more");
      } else {
        if (dataBuffer[1] == 0x76 && dataBuffer[2] == 0x83) {
          debug("got slow init initial reply, sending 7C back (after a sleep)");
          sleep(25);
          dataBuffer = [];

          // todo: invert (xor) byte 2 (x83) and send back to ecu
          // 0x83, 1000 0011 -> 0x7C 0111 1100
          // doing manually for now
          sendToEcu([0x7C])
        }
      }
      break;

    case 0x7C:
      if (dataBuffer.length < 2) {
        debug("got our 0x7C echo, waiting for E9 back still");
      } else {
        if (dataBuffer[1] == 0xE9) {
          debug("got inverted ECU address back, looks good, triggering normal init now");
          log("slow init successful");
          doingSlowInit19 = false;
          dataBuffer = [];
          resetTimeout(2000);
          initEcu131619();
        } else {
          debug("got wrong ECU address back, not good");
          dataBuffer = [];
        }
      }
      break;

    default:
      debug("unexpected byte back: "+dataBuffer[0]);
      dataBuffer = [];
      doingSlowInit19 = false;
      break;
  }
}

function parseDataBuffer() {
  switch (dataBuffer[0]) {
    case 0xCA:
      debug("got CA back");
      dataBuffer = [];
      sendToEcu([0x75]);
      break;

    case 0xCC:
      debug("got CC back, need 00 as well");
      if (dataBuffer.length >= 2 && dataBuffer[1] == 0x00) {
        log("ECU reported faults cleared");
        commandsAlert("ECU reports faults cleared");
        dataBuffer = [];
        sendNextCommand([0xCC]);
      }
      break;

    case 0x75:
      debug("got 75 back");
      dataBuffer = [];
      sendToEcu([0xF4]);
      break;

    case 0xF0:
      debug("got F0 back, checking for 50 as well");
      if (dataBuffer.length >= 2 && dataBuffer[1] == 0x50) {
        debug("got F0 50 back, looks good");
        dataBuffer = [];
        sendNextCommand(0xF0);
      } else {
        debug("need to wait for 00 still")
        break;
      }
      break;

    case 0xF4:
      debug("got F4 back, checking for 00 as well");
      if (dataBuffer.length >= 2 && dataBuffer[1] == 0x00) {
        debug("got F4 00 back, looks good");
        dataBuffer = [];
        if (!ecuConnected) {
          sendToEcu([0xD0]); // not connected so part of normal init, nearly finished
        } else {
          sendNextCommand(0xF4);
        }
      } else {
        debug("need to wait for 00 still")
        break;
      }
      break;

    case 0xD0:
      debug("got D0 back, checking for extra bytes too");
      if (dataBuffer.length >= 5) {
        debug("got D0 XX XX XX XX back, looks good");
        // ecuId = dataBuffer[1].toString(16)+","+dataBuffer[2].toString(16)+","+dataBuffer[3].toString(16)+","+dataBuffer[4].toString(16);
        // ecuId = ecuId.toUpperCase();

        ecuId1 = dataBuffer[1].toString(16);
        if (ecuId1.length == 1) { ecuId1 = "0"+ecuId1; }
        ecuId2 = dataBuffer[2].toString(16);
        if (ecuId2.length == 1) { ecuId2 = "0"+ecuId2; }
        ecuId3 = dataBuffer[3].toString(16);
        if (ecuId3.length == 1) { ecuId3 = "0"+ecuId3; }
        ecuId4 = dataBuffer[4].toString(16);
        if (ecuId4.length == 1) { ecuId4 = "0"+ecuId4; }
        ecuId = ecuId1+","+ecuId2+","+ecuId3+","+ecuId4;
        ecuId = ecuId.toUpperCase();
        updateEcuId();

        dataBuffer = [];
        setEcuConnected(true);
        log("Connected to ECU")
        sendNextCommand(0xD0);
      } else {
        debug("need to wait for more D0 bytes still")
        break;
      }
      break;

    case 0x80:
      debug("got 80 back, checking for extra bytes too");
      if (dataBuffer.length >= 2) {
        debug(dataBuffer[1]+" byte packet size");
        if (dataBuffer.length < dataBuffer[1]+1) { // +1 for command reply
          debug("Waiting for more 80 bytes");
          break;
        } else {
          debug("Got full 80 packet!");
          outputData(parseData80(dataBuffer));
          spin();
          dataBuffer = [];
          sendNextCommand(0x80);
          break;
        }
      } else {
        debug("Waiting for more 80 bytes");
      }
      break;

    case 0x7D:
      debug("got 7D back, checking for extra bytes too");
      if (dataBuffer.length >= 2) {
        debug(dataBuffer[1]+" byte packet size");
        if (dataBuffer.length < dataBuffer[1]+1) { // +1 for command reply
          debug("Waiting for more 7D bytes");
          break;
        } else {
          debug("Got full 7D packet!");
          outputData(parseData7d(dataBuffer));
          spin();
          dataBuffer = [];
          sendNextCommand(0x7D);
          break;
        }
      } else {
        debug("Waiting for more 7D bytes");
      }
      break;

    default:
      debug("Unknown data back");
      debug(dataBuffer[0]);
      dataBuffer = [];
      sendNextCommand(0xFF);
      break;
  }
}

function parseData80(data) {
  // data[0] is the command (0x80)
  data = data.slice(1) // remove it ^

  packet_size = data[0];
  debug("data 80 packet size: "+packet_size);
  // 14 bytes length for mems 1.3
  // ? bytes length for mems 1.6

  output = {};

  // byte 1-2(16 bit) - engine speed in RPM
  output['rpm'] = (data[1] << 8) + data[2];

  // byte 3 - coolant temp (+55 offset and 8 bit wrap)
  output['coolant_temp_celcius'] = data[3] - 55;
  if (output['coolant_temp_celcius'] == 59 || output['coolant_temp_celcius'] == 60) {
    output['coolant_temp_celcius'] = output['coolant_temp_celcius']+" * default value - check *";
  }
  output['coolant_temp_celcius'] = {"name": "Coolant temperature", "data": output['coolant_temp_celcius']};

  //TODO:  byte 4 - (computed) ambient temp (+55 offset, 8 bit wrap) - doesn't work, might on MPI?
  ambient_temp = data[4] - 55;
  if (ambient_temp == 200) {
    ambient_temp = ambient_temp + " * default value - check *";
  }
  output['ambient_temperature'] = {"name": "Ambient temperature", "data": ambient_temp};

  // byte 5 - intake air temp (+55 offset, 8 bit wrap)
  output['intake_air_temp_celcius'] = data[5] - 55;
  if (output['intake_air_temp_celcius'] == 35) {
    output['intake_air_temp_celcius'] = output['intake_air_temp_celcius']+" * default value - check *";
  }
  output['intake_air_temp_celcius'] = {"name": "Intake air temperature", "data": output['intake_air_temp_celcius']};

  // byte 6 - fuel temp - doesn't work on SPI, do for MPI? # defaults to FF
  output['fuel_rail_temp'] = {"name": "Fuel rail temperature", "data": data[6] - 55};

  // byte 7 - map sensor kpa
  output['map_sensor_kpa'] = {"name": "MAP sensor (kPa)", "data": data[7]};

  // byte 8 - battery voltage
  output['battery_voltage'] = {"name": "Battery voltage", "data": data[8]/10};

  // byte 9 - throttle pot voltage, WOT should be about 5v, should closed be near 0v? 0.02V per LSB. WOT should probably be close to 0xFA or 5.0V.
  output['throttle_pot_voltage'] = {"name": "Throttle voltage", "data": data[9]/200};

  // byte 10(A) - idle switch (bit 4 set if throttle closed)
  output['idle_switch'] = {"name": "Idle switch", "data": (data[10] & 0x00001000) >> 3};

  // byte 11(B) - unknown, Probably a bitfield. Observed as 0x24 with engine off, and 0x20 with engine running. A single sample during a fifteen minute test drive showed a value of 0x30.

  // byte 12(C) - Park/neutral switch. Zero is closed, nonzero is open
  output['park_or_neutral_switch'] = {"name": "Park/Neutral switch", "data": data[12]};

  // byte 13(D) - faults on mini spi
  output['fault_1_bits'] = data[13].toString(2);
  while (output['fault_1_bits'].length < 8) {
    output['fault_1_bits'] = "0"+output['fault_1_bits'];
  }
   output['fault_1_bits'] = {"name": "Fault byte 1 bits", "data": output['fault_1_bits']};

  output['fault_coolant_temp_sensor'] = {"name": "Coolant temperature fault", "data": ((data[13] >> 0) & 1) > 0};
  output['fault_inlet_air_temp_sensor'] = {"name": "Intake air temperature fault", "data": ((data[13] >> 1) & 1) > 0};
  output['fault_turbo_overboost'] = {"name": "Turbo overboosted", "data": ((data[13] >> 3) & 1) > 0};
  output['fault_ambient_temp_sensor'] = {"name": "Ambient air temperature fault", "data": ((data[13] >> 4) & 1) > 0};
  output['fault_fuel_rail_temp_sensor'] = {"name": "Fuel rail temperature fault", "data": ((data[13] >> 5) & 1) > 0};
  output['fault_knock_detected'] = {"name": "Knock detected", "data": ((data[13] >> 6) & 1) > 0};



  // byte 14(E) - fault codes
  output['fault_2_bits'] = data[14].toString(2);
  while (output['fault_2_bits'].length < 8) {
    output['fault_2_bits'] = "0"+output['fault_2_bits'];
  }
  output['fault_2_bits'] = {"name": "Fault byte 2 bits", "data": output['fault_2_bits']};

  output['fault_coolant_temp_gauge'] = {"name": "Coolant temperature gauge fault", "data": ( (data[14] >> 0) & 1 ) > 0};
  output['fault_fuel_pump_circuit'] = {"name": "Fuel pump circuit fault", "data": ( (data[14] >> 1) & 1 ) > 0};
  output['fault_air_con_clutch'] = {"name": "Air con clutch fault", "data": ( (data[14] >> 3) & 1 ) > 0};
  output['fault_purge_valve'] = {"name": "Purge valve fault", "data": ( (data[14] >> 4) & 1 ) > 0};
  output['fault_map_sensor'] = {"name": "MAP sensor fault", "data": ( (data[14] >> 5) & 1 ) > 0};
  output['fault_boost_valve'] = {"name": "Boost valve fault", "data": ( (data[14] >> 6) & 1 ) > 0};

  output['fault_throttle_pot_circuit'] = {"name": "Throttle fault", "data": ( (data[14] >> 7) & 1 ) > 0};

  // 15(F) idle setting - x6.1
  if (packet_size > 15) {
    output['idle_setpoint'] = {"name": "Idle setpoint", "data": data[15] * 6.1};
  }

  // 16 (10) unknown
  if (packet_size > 16) {
    output['idle_hotdb'] = {"name": "Idle hotdb steps", "data": data[16]};
  }

  // 17 (11) unknown

  // 18 (x12) - idle air control motor position - 0 closed, 180 fully open
  if (packet_size > 0x12) {
    output['idle_valve_position'] = {"name": "Idle valve position", "data": data[0x12]};
  }

  // 19-20 (x13-14) - idle speed deviation (16 bits)
  if (packet_size > 0x14) {
    idle_deviation = data[0x13] << 8;
    idle_deviation += data[0x14];
    output['idle_speed_deviation'] = {"name": "Idle speed deviation", "data": idle_deviation};
  }

  // 21 (x15) unknown
  if (packet_size > 0x15) {
    output['ignition_advance_offset'] = {"name": "Ignition advance offset", "data": data[0x15]};
  }

  // TODO: 22 (x16) - ignition advance 0.5 degrees per lsb, range -24 deg (00) to 103.5 deg (0xFF)
  if (packet_size > 0x16) {
    advance = data[0x16];
    // advance /= 2;
    // advance -= 24;
    output['ignition_advance_raw'] = {"name": "Ignition advance raw", "data": advance};
  }

  // TODO: 23-24 (x17-18) - coil time 0.002ms per lsb (16 bit)
  if (packet_size > 0x18) {
    coil_time = data[0x17] << 8;
    coil_time += data[0x18];
    coil_time *= 2;
    output['coil_time_microseconds'] = {"name": "Coil time microseconds", "data": coil_time};
  }
  // 25 (x19) unknown
  // 26 (x1a) unknown
  // 27 (x1B) unknown


  return(output);
}

function parseData7d(data) {
  // data[0] is the command (0x7D)
  data = data.slice(1) // remove it ^

  packet_size = data[0];
  debug("data 7D packet size: "+packet_size);

  output = {};

  output['ignition_switch'] = {"name": "Ignition switch", "data": data[1]};

  output['throttle_angle'] = {"name": "Throttle angle", "data": data[2]/2};

  // 0x03  Unknown

  output['air_fuel_ratio'] = {"name": "Air/fuel ratio", "data": data[4]}; // "A/F ratio? might just be 0xFF (unknown)" ## if it's FF then don't output?

  // 0x05 Unknown
  dtc_byte = data[5];
  output['lambda_heater_relay'] = {"name": "Lambda heater relay", "data": ((dtc_byte >> 3) & 1) > 0 };
  output['crank_shaft_sync'] = {"name": "Crank shaft sync", "data": ((dtc_byte >> 4) & 1) > 0 };
  output['fan_1_control'] = {"name": "Fan 1 control", "data": ((dtc_byte >> 5) & 1) > 0 };
  output['fan_2_control'] = {"name": "Fan 2 control", "data": ((dtc_byte >> 7) & 1) > 0 };

  // output['raw_lambda'] = data[6];
  output['lambda_mv'] = data[6] * 5;
  if (output['lambda_mv'] == 440 || output['lambda_mv'] == 435) {
    output['lambda_mv'] = output['lambda_mv']+" * default value - check *";
  }
  output['lambda_mv'] = {"name": "Lambda mV", "data": output['lambda_mv']};

  output['lambda_sensor_frequency'] = data[7];
  if (output['lambda_sensor_frequency'] == 0) {
  } else if (output['lambda_sensor_frequency'] == 255) {
    output['lambda_sensor_frequency'] = output['lambda_sensor_frequency']+" (not reporting)";
  }
  output['lambda_sensor_frequency'] = {"name": "Lambda sensor frequency", "data": output['lambda_sensor_frequency']};

  output['lambda_sensor_duty_cycle'] = {"name": "Lambda sensor duty cycle", "data": data[8]};
  output['lambda_sensor_status'] = {"name": "Lambda sensor status", "data": data[9]}; // "Lambda sensor status? 0x01 for good, any other value for no good"

  output['closed_loop'] = {"name": "Closed loop", "data": data[10] > 0}; // "Loop indicator, 0 for open loop and nonzero for closed loop"

  output['long_term_trim'] = {"name": "Long term trim", "data": data[11]};
  output['short_term_trim_percent'] = {"name": "Short term trim", "data": data[12]};
  output['carbon_can_purge_valve_duty_cycle'] = {"name": "Purge valve duty cycle", "data": data[13]}; // "Carbon canister purge valve duty cycle?"

  // 0x0E  Unknown
  dtc2 = data[0xE];
  output['cam_shaft_sync'] = {"name": "Cam shaft sync", "data": ((dtc2 >> 1) & 1) > 0 };

  if (packet_size >= 16) {
    output['idle_base_position'] = {"name": "Idle base position", "data": data[15]};
  }

  // 0x10  Unknown
  // 0x11  Unknown
  // 0x12  Unknown
  // 0x13  Unknown
  if (packet_size >= 21) {
    output['idle_error'] = {"name": "Idle error", "data": data[20]};
  }
  // 0x15  Unknown

  if (packet_size >= 0x16) {
    dtc3 = data[0x16];
    output['injector_1_4_driver'] = {"name": "Injector 1/4 driver", "data": ((dtc_byte >> 1) & 1) > 0 };
    output['injector_2_3_driver'] = {"name": "Injector 2/3 driver", "data": ((dtc_byte >> 2) & 1) > 0 };
    output['fault_engine_bay_vent_warning'] = {"name": "Engine bay ventilation warning", "data": ((dtc_byte >> 3) & 1) > 0 };
    output['engine_bay_vent_relay'] = {"name": "Engine bay ventilation relay", "data": ((dtc_byte >> 4) & 1) > 0 };
    output['hill_assist'] = {"name": "Hill assist", "data": ((dtc_byte >> 5) & 1) > 0 };
    output['cruise_control'] = {"name": "Cruise control", "data": ((dtc_byte >> 6) & 1) > 0 };
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
    output['crank_counter'] = {"name": "Crank counter", "data": data[0x1F]};
  }
  return output;

}
