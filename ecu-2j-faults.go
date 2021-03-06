package main

func twojParseFaults(buffer []byte) {
	// fmt.Printf("got %d bytes \n%s", len(buffer), hex.Dump(buffer))
/*
MPI mini fault ram locations
9h





*/
	
	faults := []string {}

	// RAM 594h 4h
	if len(buffer) >= 5 {
	  if (buffer[4] & 0b01000000) > 0 { faults = append(faults, "Outside air temp (low voltage)") }
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
        // buffer[6]
	// buffer[7]

	
	// RAM 590h 4h
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
	if len(buffer) >= 11 {
	  // buffer[10];
	}
	if len(buffer) >= 12 {
	  // buffer[11];
	}
	
	
	// 14h 4h
	if len(buffer) >= 13 {
	  if ((buffer[12] >> 6) & 1) > 0 { faults = append(faults, "Outside temp sensor (present)") }
	  if ((buffer[12] >> 5) & 1) > 0 { faults = append(faults, "Power supply (present)") }
	  if ((buffer[12] >> 4) & 1) > 0 { faults = append(faults, "Oil temp (present)") }
	  if ((buffer[12] >> 2) & 1) > 0 { faults = append(faults, "Coolant temp (present)") } // one of these two is a typo
          if ((buffer[12] >> 2) & 1) > 0 { faults = append(faults, "System voltage (present)") } // one of these two is a typo
	}
	if len(buffer) >= 14 {
	  if ((buffer[13] >> 7) & 1) > 0 { faults = append(faults, "Battery voltage (present)") }
	  if ((buffer[13] >> 4) & 1) > 0 { faults = append(faults, "Lambda 1 bank 1 (present)") }
	  if ((buffer[13] >> 2) & 1) > 0 { faults = append(faults, "Throttle pot (present)") }
	  if ((buffer[13] >> 1) & 1) > 0 { faults = append(faults, "Intake air temp (present)") }
	  if ((buffer[13] >> 0) & 1) > 0 { faults = append(faults, "MAP sensor (present)") }
	}
	// 14
	// 15

	// 228h 4h
	// 16
	// 17
	// 18
	// 19
	
	// 598h 4h
	// 20
	// 21
	// 22
	if len(buffer) >= 24 {
	  if (buffer[23] & 0b1000) > 0 { faults = append(faults, "MAP sensor (present 2)") }
	  if (buffer[23] & 0b100) > 0 { faults = append(faults, "Oil temp (present 2)") }
	  if (buffer[23] & 0b10) > 0 { faults = append(faults, "Intake air temp (present 2)") }
	  if (buffer[23] & 0b1) > 0 { faults = append(faults, "Coolant temp (present 2)") }
	}

	// 374h 2h
	// 24 0x374
	if len(buffer) >= 26 {
	  if (buffer[25] & 0b1000) > 0 { faults = append(faults, "MAP sensor (historic)") }
	  if (buffer[25] & 0b100) > 0 { faults = append(faults, "Oil temp (historic)") }
	  if (buffer[25] & 0b10) > 0 { faults = append(faults, "Intake air temp (historic)") }
	  if (buffer[25] & 0b1) > 0 { faults = append(faults, "Coolant temp (historic)") }
	}

	// 5B0h 2h
	if len(buffer) >= 27 {
	  // MPI ram 513
	  if ((buffer[26] >> 0) & 1) > 0 { faults = append(faults, "Road speed sensor (present)") }
	  if ((buffer[26] >> 1) & 1) > 0 { faults = append(faults, "Comm. with AT (present)") }
	  // bit 2 - not set in code
	  // bit 3 - not set in code
          if ((buffer[26] >> 4) & 1) > 0 { faults = append(faults, "Bank 1 fuel feedback (present)") }
	  if ((buffer[26] >> 5) & 1) > 0 { faults = append(faults, "Bank 2 fuel feedback (present)") }
	  // bit 6 - not set in code
	  // bit 7 - not set in code
	}
	// 0x5b1 buffer[27] - this is a copy of 375 above
	
	// 513h 1h
	// buffer[28]
	if len(buffer) >= 29 {
	  if (buffer[28] & 0b00000001) > 0 { faults = append(faults, "Road speed sensor (historic)") }
	  if (buffer[28] & 0b00000010) > 0 { faults = append(faults, "Comm. with AT (historic)") }
	  if (buffer[28] & 0b00010000) > 0 { faults = append(faults, "Feedback (historic)") }
	}
	
	// 5B8h 1
	// buffer[29]
	// this is a copy of 513 above


	globalFaults = faults
}
