package randomart

import "testing"

func TestTileSet_Index(t *testing.T) {
	tests := []struct {
		name    string
		tileset TileSet
		n       int
		want    rune
	}{
		{
			name:    "wraps within range",
			tileset: TileSet{Runes: []rune{'a', 'b', 'c'}},
			n:       1,
			want:    'b',
		},
		{
			name:    "wraps at boundary when overflow prevention disabled",
			tileset: TileSet{Runes: []rune{'a', 'b', 'c'}},
			n:       3,
			want:    'a',
		},
		{
			name:    "wraps past boundary when overflow prevention disabled",
			tileset: TileSet{Runes: []rune{'a', 'b', 'c'}},
			n:       8,
			want:    'c',
		},
		{
			name:    "clamps at boundary when overflow prevention enabled",
			tileset: TileSet{Runes: []rune{'a', 'b', 'c'}, PreventRuneOverflow: true},
			n:       3,
			want:    'c',
		},
		{
			name:    "clamps past boundary when overflow prevention enabled",
			tileset: TileSet{Runes: []rune{'a', 'b', 'c'}, PreventRuneOverflow: true},
			n:       99,
			want:    'c',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tileset.Index(tt.n); got != tt.want {
				t.Errorf("TileSet.Index(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestTileSets(t *testing.T) {
	t.Run("returns bundled tilesets", func(t *testing.T) {
		got := TileSets()
		if len(got) != 2 {
			t.Fatalf("TileSets() len = %d, want 2", len(got))
		}
	})
}
