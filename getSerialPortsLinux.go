//go:build linux

package main

import (
)

func nativeGetPortsList() ([]string, error) {
	return []string{"/dev/ttyUSB0"}, nil
}
