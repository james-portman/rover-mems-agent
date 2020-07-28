package main

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
