package catomic

func uint2Bool(i uint32) bool {
	if i == 0 {
		return false
	}

	return true
}

func bool2Uint(b bool) uint32 {
	if b {
		return 1
	}

	return 0
}
