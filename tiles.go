package randomart

type TileSet struct {
	ID string // unique name

	Runes []rune

	Start rune // special rune for the starting position (optional)
	End   rune // special rune for the ending position (optional)

	// If PreventRuneOverflow is true, the last value of Runes will be used for
	// any TileSet.Index(n) where n >= len(Runes). Otherwise, the default
	// behavior will be used which is to wrap around.
	PreventRuneOverflow bool
}

func (t TileSet) Index(n int) rune {
	if t.PreventRuneOverflow && n >= len(t.Runes) {
		return t.Runes[len(t.Runes)-1]
	}

	return t.Runes[n%len(t.Runes)]
}

var (
	// All TileSets bundled with the package.
	BundledTileSets = []TileSet{SSHTiles, GalaxyTiles}

	// Classic OpenSSH randomart tileset, using basic ASCII.
	SSHTiles = TileSet{
		ID:    "openssh",
		Runes: []rune{' ', '.', 'o', '+', '=', '*', 'B', 'O', 'X', '@', '%', '&', '#', '/', '^'},
		Start: 'S',
		End:   'E',
	}

	// A spacey emoji based tileset.
	GalaxyTiles = TileSet{
		ID:    "galaxy",
		Runes: []rune{'🌑', '🌒', '🌓', '🌔', '🌕', '🪐', '🌖', '🌗', '🌘'},
		Start: '🌝',
		End:   '🌚',
	}
)
