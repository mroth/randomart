package randomart

import (
	"bytes"
	"errors"
	"strings"
	"unicode/utf8"
)

type Board struct {
	data       []uint8
	dimX, dimY int
	start, end position
	pos        position
}

type position struct{ x, y int }

// NewBoard initializes a Board of dimensions x, y with a starting position of
// the center point.
//
// x 0→n, y 0↓n.
func NewBoard(x, y int) (*Board, error) {
	if x <= 0 || y <= 0 {
		return nil, errors.New("invalid dimensions")
	}

	data := make([]uint8, x*y)

	b := Board{
		data: data,
		dimX: x,
		dimY: y,
	}

	err := b.setStartPos(x/2, y/2)
	return &b, err
}

func (b *Board) setStartPos(x, y int) error {
	if x >= b.dimX || x < 0 || y >= b.dimY || y < 0 {
		return errors.New("invalid start position")
	}
	b.pos = position{x: x, y: y}
	b.start = position{x: x, y: y}
	return nil
}

// Write writes len(p) bytes to the underlying Board. The provided fingerprint
// will be used to explore the board using the drunken bishop algorithm.
//
// Implements the io.Writer interface. The returned number of bytes will always
// equal len(fingerprint), and the error will always be nil
func (b *Board) Write(fingerprint []byte) (n int, err error) {
	// leave breadcrumb at start position
	b.increment(b.pos.x, b.pos.y)

	for _, fingerByte := range fingerprint { // for each byte of fingerprint
		for s := uint(0); s < 8; s += 2 { // stride byte in 2 bit chunks
			// bitwise right shift the byte by s, which results in last 2 bits
			// being the ones we want to analyze. then take the arithmetic AND
			// of int 3 (0b11), in order to effectively extract the value of
			// those last two bits.
			//
			// elegant bitshift solution found via:
			// https://github.com/calmh/randomart/blob/master/randomart.go
			d := (fingerByte >> s) & 0b11

			// move in direction
			switch d {
			case 0: // 0b00 ↖
				b.moveLeft()
				b.moveUp()
			case 1: // 0b01 ↗
				b.moveRight()
				b.moveUp()
			case 2: // 0b10 ↙
				b.moveLeft()
				b.moveDown()
			case 3: // 0b11 ↘
				b.moveRight()
				b.moveDown()
			}

			// mark breadcrumb after move
			b.increment(b.pos.x, b.pos.y)
		}

	}
	b.end = b.pos

	return len(fingerprint), nil
}

// move left if possible
func (b *Board) moveLeft() {
	if b.pos.x > 0 {
		b.pos.x--
	}
}

// move right if possible
func (b *Board) moveRight() {
	if b.pos.x < b.dimX-1 {
		b.pos.x++
	}
}

// move up if possible
func (b *Board) moveUp() {
	if b.pos.y > 0 {
		b.pos.y--
	}
}

// move down if possible
func (b *Board) moveDown() {
	if b.pos.y < b.dimY-1 {
		b.pos.y++
	}
}

// increment the value at the given position
func (b *Board) increment(x, y int) {
	b.data[y*b.dimX+x]++
}

// get the value at the given position
func (b *Board) getValue(x, y int) uint8 {
	return b.data[y*b.dimX+x]
}

// Renders output from the current state of Board b using TileSet t.
func (b *Board) Render(t TileSet) []byte {
	var buf bytes.Buffer
	runeLen := utf8.RuneLen(t.Runes[0]) // assume first rune is avg length (not always accurate)
	buf.Grow(((b.dimX * runeLen) + 1) * b.dimY)
	for y := 0; y < b.dimY; y++ {
		for x := 0; x < b.dimX; x++ {
			pos := position{x: x, y: y}
			switch {
			case pos == b.start && t.Start != 0:
				buf.WriteRune(t.Start)
			case pos == b.end && t.End != 0:
				buf.WriteRune(t.End)
			default:
				buf.WriteRune(t.Index(int(b.getValue(x, y))))
			}
		}
		buf.WriteRune('\n')
	}
	return buf.Bytes()
}

// Armor wraps the lines of a rendered output b in a simple ASCII box.
func Armor(b []byte) []byte {
	// This could be done much more efficiently with a Scanner, but since we're
	// working on very small data and it's a proof of concept, optimize for
	// simplicity and understandability.
	lines := bytes.Split(b, []byte("\n"))
	nDataCols := len(lines[0])

	var buf bytes.Buffer
	buf.WriteRune('+')
	buf.WriteString(strings.Repeat("-", nDataCols))
	buf.WriteRune('+')
	buf.WriteRune('\n')

	for _, row := range lines {
		if len(row) == nDataCols {
			buf.WriteRune('|')
			buf.Write(row)
			buf.WriteRune('|')
			buf.WriteRune('\n')
		}
	}

	buf.WriteRune('+')
	buf.WriteString(strings.Repeat("-", nDataCols))
	buf.WriteRune('+')
	buf.WriteRune('\n')

	return buf.Bytes()
}
