package main

// specifically for 1.2,1.3,1.6

import (
	"fmt"
	"github.com/distributed/sers"
)

func readFirstBytesFromPortEcu1x(fn string) ([]byte, error) {

	fmt.Println("Connecting to MEMS 1.x (1.2, 1.3, 1.6) ECU")
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

	// mode, err := sp.GetMode()
	// fmt.Println("Serial cable set to:")
	// fmt.Println(mode)

	return ecu1xLoop(sp, true)

}
