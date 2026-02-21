package randomart

import (
	"bytes"
	"errors"
	"io"
	"strings"
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

	// default to OpenSSHTiles if no tiles provided
	tileset := opts.Tiles
	if len(tileset.Runes) == 0 {
		tileset = OpenSSHTiles
	}

	// max buffer size is 4 bytes per cell (worst case with
	// multi-byte runes and border), pre-allocate to avoid
	// resizing during render.
	var buf bytes.Buffer
	buf.Grow((b.dimX + 2) * (b.dimY + 2) * 4)

	if opts.Border {
		buf.WriteString(borderLine(b.dimX))
	}

	for y := range b.dimY {
		if opts.Border {
			buf.WriteByte('|')
		}

		for x := range b.dimX {
			pos := position{x: x, y: y}
			switch {
			case pos == b.start && tileset.Start != 0:
				buf.WriteRune(tileset.Start)
			case pos == b.end && tileset.End != 0:
				buf.WriteRune(tileset.End)
			default:
				buf.WriteRune(tileset.Index(int(b.getValue(x, y))))
			}
		}

		if opts.Border {
			buf.WriteByte('|')
		}
		buf.WriteByte('\n')
	}

	if opts.Border {
		buf.WriteString(borderLine(b.dimX))
	}

	return io.Copy(w, &buf)
}

func borderLine(width int) string {
	return "+" + strings.Repeat("-", width) + "+\n"
}
