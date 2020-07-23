package main

import (
	"fmt"
	"time"
	"errors"
	"encoding/hex"
	"github.com/distributed/sers"
)

var (
	ecu19SpecificInitCommand = []byte {0x7C}

	ecu19WokeResponse = []byte {0x55, 0x76, 0x83}
  ecu19SpecificInitResponse = []byte {ecu19SpecificInitCommand[0], 0xE9} // includes our echo
)

func sleepUntil(start time.Time, plus int) {
  target := start.Add(time.Duration(plus) * time.Millisecond)
  sleepMs := target.Sub(time.Now()).Milliseconds()
  // fmt.Println("Sleeping for ms:")
  // fmt.Println(sleepMs)
  if sleepMs < 0 { return }
  time.Sleep(time.Duration(sleepMs) * time.Millisecond)
}

func readFirstBytesFromPortEcu19(fn string) ([]byte, error) {

	fmt.Println("Connecting to MEMS 1.9 ECU")
	globalConnected = false

	sp, err := sers.Open(fn)
	if err != nil {
		return nil, err
	}
	defer sp.Close()

	err = sp.SetMode(9600, 8, sers.N, 1, sers.NO_HANDSHAKE)
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

	// try the normal method first
	ecu1xLoop(sp, true)

  // clear the line
	sp.SetBreak(false)
	time.Sleep(2000 * time.Millisecond)

  start := time.Now()

  // start bit
	sp.SetBreak(true)
  sleepUntil(start, 200)

  // send the byte
  ecuAddress := 0x16
  for i:=0; i<8; i++ {

    bit := (ecuAddress >> i) & 1;
    if (bit > 0) {
        sp.SetBreak(false)
    } else {
        sp.SetBreak(true)
    }

    sleepUntil(start, 200 + ((i+1)*200))

  }
  // stop bit
	sp.SetBreak(false)
  sleepUntil(start, 200 + (8*200) + 200)

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

		if slicesEqual(buffer, ecu19WokeResponse) {
      fmt.Println("1.9 ECU woke up - init stage 1")
			buffer = nil
			time.Sleep(50 * time.Millisecond) // TODO: is this the right sleep?
      // todo: invert (xor) byte 2 (x83) and send back to ecu
      // 0x83, 1000 0011 -> 0x7C 0111 1100
      // doing manually for now (doesn't hurt)
      sp.Write(ecu19SpecificInitCommand)
			continue
		}

    if slicesEqual(buffer, ecu19SpecificInitResponse) {
      fmt.Println("1.9 ECU init stage 2")
      buffer = nil
      ecu1xLoop(sp, true)
      continue
    }

	}
	if readLoops >= readLoopsLimit {
		fmt.Printf("1.9 had buffer data: got %d bytes \n%s", len(buffer), hex.Dump(buffer))
		return nil, errors.New("1.9 readloop timed out")
	}
	fmt.Println("fell out of 1.9 readloop")

	return nil, err
}
