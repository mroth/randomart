package randomart

// A TileSet defines the symbols used when rendering a Board visually.
type TileSet struct {
	ID    string // A unique name to identify the tile set.
	Runes []rune // The set of symbols used to represent visit counts. The first rune is used for 0 visits, the second for 1 visit, etc.
	Start rune   // A special symbol used for the starting position (optional).
	End   rune   // A special symbol used for the ending position (optional).

	// If PreventRuneOverflow is true, the last value of Runes will be used for
	// any TileSet.Index(n) where n >= len(Runes). Otherwise, the default
	// behavior will be used which is to wrap around.
	PreventRuneOverflow bool
}

// Index returns the appropriate rune for a given visit count n.
func (t TileSet) Index(n int) rune {
	if t.PreventRuneOverflow && n >= len(t.Runes) {
		return t.Runes[len(t.Runes)-1]
	}

	return t.Runes[n%len(t.Runes)]
}

// This package includes a few bundled tile sets, but users are free to define their own as well.
var (
	// OpenSSHTiles is the classic OpenSSH randomart tile set using basic ASCII runes.
	OpenSSHTiles = TileSet{
		ID:    "openssh",
		Runes: []rune{' ', '.', 'o', '+', '=', '*', 'B', 'O', 'X', '@', '%', '&', '#', '/', '^'},
		Start: 'S',
		End:   'E',
	}

	// GalaxyTiles is a spacey emoji-based tile set.
	GalaxyTiles = TileSet{
		ID:    "galaxy",
		Runes: []rune{'🌑', '🌒', '🌓', '🌔', '🌕', '🪐', '🌖', '🌗', '🌘'},
		Start: '🌝',
		End:   '🌚',
	}
)

// TileSets returns the collection of bundled tile sets in this package.
func TileSets() []TileSet {
	return []TileSet{OpenSSHTiles, GalaxyTiles}
}
