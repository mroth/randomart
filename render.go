package randomart

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

var (
	errNilWriter = errors.New("nil writer")
)

// RenderOptions customizes output behavior for RenderTo.
type RenderOptions struct {
	Tiles  TileSet
	Border bool
	Header string
	Footer string
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
	borderWidth := b.dimX * maxRuneCellWidth(tileset)

	if opts.Border {
		buf.WriteString(borderLineWithLabel(borderWidth, opts.Header))
	} else if opts.Header != "" {
		buf.WriteString(opts.Header)
		buf.WriteByte('\n')
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
		buf.WriteString(borderLineWithLabel(borderWidth, opts.Footer))
	} else if opts.Footer != "" {
		buf.WriteString(opts.Footer)
		buf.WriteByte('\n')
	}

	return io.Copy(w, &buf)
}

func borderLine(width int) string {
	return "+" + strings.Repeat("-", width) + "+\n"
}

func borderLineWithLabel(width int, label string) string {
	if label == "" {
		return borderLine(width)
	}

	maxLabelWidth := width - 2 // account for surrounding brackets: [label]
	if maxLabelWidth <= 0 {
		return borderLine(width)
	}

	truncatedLabel := label
	if runewidth.StringWidth(label) > maxLabelWidth {
		truncatedLabel = runewidth.Truncate(label, maxLabelWidth, "…")
	}

	content := "[" + truncatedLabel + "]"
	remaining := width - runewidth.StringWidth(content)
	leftDashes := remaining / 2
	rightDashes := remaining - leftDashes

	return "+" + strings.Repeat("-", leftDashes) + content + strings.Repeat("-", rightDashes) + "+\n"
}

func maxRuneCellWidth(t TileSet) int {
	maxWidth := 1
	for _, r := range t.Runes {
		maxWidth = max(maxWidth, runewidth.RuneWidth(r))
	}
	maxWidth = max(maxWidth, runewidth.RuneWidth(t.Start))
	maxWidth = max(maxWidth, runewidth.RuneWidth(t.End))
	return maxWidth
}
