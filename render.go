package randomart

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

var (
	errNilWriter = errors.New("nil writer")
)

// RenderOptions customizes output behavior for RenderTo.
type RenderOptions struct {
	Tiles  TileSet
	Border bool
}

// RenderTo writes output from the current state of Board b to w.
//
// If opts.Tiles has no runes, OpenSSHTiles is used by default.
func RenderTo(w io.Writer, b *Board, opts RenderOptions) (int64, error) {
	if w == nil {
		return 0, errNilWriter
	}

	tileset := opts.Tiles
	if len(tileset.Runes) == 0 {
		tileset = OpenSSHTiles
	}

	raw := renderBoard(b, tileset)
	if opts.Border {
		raw = borderBytes(raw)
	}

	n, err := w.Write(raw)
	return int64(n), err
}

func renderBoard(b *Board, t TileSet) []byte {
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

func borderBytes(b []byte) []byte {
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
