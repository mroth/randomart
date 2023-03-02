package randomart

import (
	"errors"
	"strings"
)

type Board struct {
	data       [][]uint8
	dimX, dimY int
	start, end position
	pos        position
}

type position struct{ x, y int }

// NewBoard intializes a Board of dimensions x, y with a starting position of
// the center point.
//
// x 0→n, y 0↓n.
func NewBoard(x, y int) *Board {
	data := make([][]uint8, y)
	for i := range data {
		data[i] = make([]uint8, x)
	}

	b := Board{
		data: data,
		dimX: x,
		dimY: y,
	}

	err := b.setStartPos(x/2, y/2)
	if err != nil {
		panic(err) // should be unreachable in current usage
	}
	return &b
}

func (b *Board) setStartPos(x, y int) error {
	if x >= b.dimX || x < 0 || y >= b.dimY || y < 0 {
		return errors.New("invalid start position")
	}
	b.pos = position{x: x, y: y}
	b.start = position{x: x, y: y}
	return nil
}

// Write writes len(p) bytes to the underlying Board.  The provided fingerprint
// will be used to explore the board using the drunken bishop algorithm.
//
// Implements the io.Writer interface. The returned number of bytes will always
// equal len(fingerprint), and the error will always be nil
func (b *Board) Write(fingerprint []byte) (n int, err error) {
	// leave breadcrumb at start position
	b.data[b.pos.y][b.pos.x]++

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
			b.data[b.pos.y][b.pos.x]++
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

func (b *Board) RenderString(t TileSet) string {
	var sb strings.Builder
	for y := range b.data {
		for x := range b.data[y] {
			pos := position{x: x, y: y}
			switch {
			case pos == b.start && t.Start != 0:
				sb.WriteRune(t.Start)
			case pos == b.end && t.End != 0:
				sb.WriteRune(t.End)
			default:
				sb.WriteRune(t.Index(int(b.data[y][x])))
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func Armor(b string) string {
	// This could be done much more efficiently with a Scanner, but since we're
	// working on very small data and it's a proof of concept, optimize for
	// simplicity and understandability.
	lines := strings.Split(b, "\n")
	nDataCols := len(lines[0])

	var sb strings.Builder
	sb.WriteRune('+')
	sb.WriteString(strings.Repeat("-", nDataCols))
	sb.WriteRune('+')
	sb.WriteRune('\n')

	for _, row := range lines {
		if len(row) == nDataCols {
			sb.WriteRune('|')
			sb.WriteString(row)
			sb.WriteRune('|')
			sb.WriteRune('\n')
		}
	}

	sb.WriteRune('+')
	sb.WriteString(strings.Repeat("-", nDataCols))
	sb.WriteRune('+')
	sb.WriteRune('\n')

	return sb.String()
}
