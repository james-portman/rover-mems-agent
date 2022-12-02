package main

import (
	"github.com/distributed/sers"
)
// go routines for serial read
// mostly because reads erroneously block in Linux, even with a timeout set

func serialReadRoutine(sp sers.SerialPort) {
    for {
    	rb := make([]byte, 256)
		n, _ := sp.Read(rb[:])
		rb = rb[0:n] // chop down to actual data size
		for i := 0; i < n; i++ {
        	serialReadChannel <- rb[i]
        }
    }
}

// reads all currently available data from the channel
func nonBlockingSerialRead() ([]byte) {
	buffer := make([]byte, 0)
	outer:
	for {
		select {
		    case msg := <-serialReadChannel:
		        buffer = append(buffer, msg)
		    default:
		        break outer
	    }
	}
	return buffer
}
