//go:build linux

package main

import (
	"os"
	"path/filepath"
	"fmt"
	"strings"
)

func nativeGetPortsList() ([]string, error) {
	output := []string {}

    dir, err := os.Open("/dev/")
    if err != nil {
        panic(err)
    }
    defer dir.Close()

    filepath.Walk(dir.Name(), func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if (strings.HasPrefix(info.Name(), "ttyUSB") || strings.HasPrefix(info.Name(), "ttyACM")) {
        	fmt.Println(info.Name())
        	output = append(output, "/dev/"+info.Name())
        }
        return nil
    })

	return output, nil
}
