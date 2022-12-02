package main

import (
  "time"
)

func slicesEqual(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func xor_all_bytes(bytes []byte) byte {
  output := byte(0)
  for i:=0; i<len(bytes); i++ {
    output = output ^ bytes[i]
  }
  return output
}

func sleepUntil(start time.Time, plus int) {
  target := start.Add(time.Duration(plus) * time.Millisecond)
  sleepMs := target.Sub(time.Now()).Milliseconds()
  // fmt.Println("Sleeping for ms:")
  // fmt.Println(sleepMs)
  if sleepMs < 0 { return }
  time.Sleep(time.Duration(sleepMs) * time.Millisecond)
}

func timestampMs() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}