package main

import (
	"bytes"
	"encoding/binary"
)

func convertIntsToBytes(ints []int) []byte {
	b := new(bytes.Buffer)

	fixedSizeInts := make([]int16, len(ints))
	for idx, i := range ints {
		fixedSizeInts[idx] = int16(i)
	}

	err := binary.Write(b, binary.LittleEndian, fixedSizeInts)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}
