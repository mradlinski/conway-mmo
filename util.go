package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func convertIntsToBytes(ints []int) []byte {
	b := new(bytes.Buffer)

	fixedSizeInts := make([]int32, len(ints))
	for idx, i := range ints {
		fixedSizeInts[idx] = int32(i)
	}

	err := binary.Write(b, binary.LittleEndian, fixedSizeInts)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}

func colorToRGBInt(c colorful.Color) int {
	r, g, b := c.RGB255()
	return (int(r) << 16) + (int(g) << 8) + int(b)
}

func initializeColorPicker() func() int {
	palette, err := colorful.SoftPalette(300)
	if err != nil {
		panic(err)
	}

	return func() int {
		c := palette[rand.Intn(300)]

		return colorToRGBInt(c)
	}
}
