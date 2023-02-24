package randomart

type TileSet struct {
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
	SSHTiles = TileSet{
		Runes: []rune{' ', '.', 'o', '+', '=', '*', 'B', 'O', 'X', '@', '%', '&', '#', '/', '^'},
		Start: 'S',
		End:   'E',
	}

	GalaxyTiles = TileSet{
		Runes: []rune{'🌑', '🌒', '🌓', '🌔', '🌕', '🪐', '🌖', '🌗', '🌘'},
		Start: '🌝',
		End:   '🌚',
	}
)
