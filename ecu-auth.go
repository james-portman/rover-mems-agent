package main



func bit(bit_num int, thebyte int) int {
  // returns the requested bit, counting from 0-7, 0-15 for doubles
  return (thebyte >> bit_num) & 1;
}

func generateKey(seed int) int {
  key := 0
  loops := 1

  if (bit(15,seed) > 0) { loops += 8 }
  if (bit(7,seed) > 0) { loops += 4 }
  if (bit(4,seed) > 0) { loops += 2 }
  if (bit(0,seed) > 0) { loops += 1 }

  for loops > 0 {
    key = seed >> 1 // take the seed shifted right by 1 (each loop changes seed)

    if (bit(13,seed) > 0 && bit(3,seed) > 0) {
      key &= 0b11111111111111110 // unset LSB
    } else {
      key |= 0b0000000000000001 // set LSB
    }

    xors := bit(9,seed) ^ bit(8,seed) ^ bit(2,seed) ^ bit(1,seed)
    if (xors > 0) {
      key |= 0b1000000000000000 // set msb
    }

    seed = key
    loops--
  }

  return key
}


/*
different version from python someone made:
for idx in range(0, count):
        tap = ((seed >> 1) + (seed >> 2 ) + (seed >> 8 ) + (seed >> 9)) & 1
        tmp = (seed >> 1) | ( tap << 0xF)
        if (seed >> 0x3 & 1) and (seed >> 0xD & 1):
                seed = tmp & ~1
        else:
                seed = tmp | 1
*/
