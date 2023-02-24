package main

import (
	"fmt"

	"github.com/mroth/randomart"
)

func main() {
	fingerprint := []byte{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13,
		0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed}

	b := randomart.NewBoard(17, 9)
	b.Write(fingerprint)
	fmt.Println(b.RenderString(randomart.SSHTiles))

	b = randomart.NewBoard(8, 8)
	b.Write(fingerprint)
	fmt.Println(b.RenderString(randomart.GalaxyTiles))
}
