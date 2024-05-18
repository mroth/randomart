package randomart

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewBoard(t *testing.T) {
	type args struct {
		x, y int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{x: 0, y: 0}, wantErr: true},
		{args: args{x: 1, y: 0}, wantErr: true},
		{args: args{x: 0, y: 1}, wantErr: true},
		{args: args{x: 1, y: 1}, wantErr: false},
		{args: args{x: -1, y: 1}, wantErr: true},
		{args: args{x: 10, y: 10}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBoard(tt.args.x, tt.args.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBoard(%v, %v) error = %v, wantErr %v", tt.args.x, tt.args.y, err, tt.wantErr)
				return
			}
		})
	}
}

func TestBoard_Write(t *testing.T) {
	type state struct {
		data  []uint8
		start position
		end   position
	}
	tests := []struct {
		fingerprint []byte
		dimX        int
		dimY        int
		wantErr     bool // should always be false with current setup, but allow override
		expectState state
	}{
		{
			fingerprint: []byte{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed},
			dimX:        17,
			dimY:        9,
			expectState: state{
				data: []uint8{
					0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 1, 1, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 1, 2, 0, 4, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 3, 0, 1, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 2, 0, 6, 0, 1, 0, 1, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 2, 1, 1, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 1, 1, 1, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 1, 1, 1, 0, 0, 0,
				},
				start: position{x: 8, y: 4},
				end:   position{x: 6, y: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(base64.RawStdEncoding.EncodeToString(tt.fingerprint), func(t *testing.T) {
			b, err := NewBoard(tt.dimX, tt.dimY)
			if err != nil {
				t.Fatal(err)
			}
			gotN, err := b.Write(tt.fingerprint)
			// check for error (always false currently, just there for io.Writer
			if (err != nil) != tt.wantErr {
				t.Errorf("Board.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// no error possibility, so n should always be number of bytes sent
			if wantN := len(tt.fingerprint); gotN != wantN {
				t.Errorf("Board.Write() = %v, want %v", gotN, wantN)
			}

			// compare internal state
			gotState := state{data: b.data, start: b.start, end: b.end}
			if !reflect.DeepEqual(tt.expectState, gotState) {
				t.Errorf("Board internal state: want %v, got %+v", tt.expectState, gotState)
			}
		})
	}
}

func TestBoard_Render(t *testing.T) {
	datacases := [][]byte{
		{},
		{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed},
	}

	rendercases := []struct {
		tiles      TileSet
		dimX, dimY int
		armor      bool
	}{
		{tiles: SSHTiles, dimX: 17, dimY: 9, armor: true},
		{tiles: GalaxyTiles, dimX: 10, dimY: 10, armor: false},
	}

	for _, dc := range datacases {
		slug := hex.EncodeToString(dc)
		if slug == "" {
			slug = "_empty"
		}
		t.Run(slug, func(t *testing.T) {

			for _, rc := range rendercases {
				specifier := fmt.Sprintf("%s-%dx%d%s",
					rc.tiles.ID,
					rc.dimX, rc.dimY,
					func() string {
						if rc.armor {
							return "-armored"
						}
						return ""
					}(),
				)
				t.Run(specifier, func(t *testing.T) {
					filename := strings.Join([]string{slug, specifier, "txt"}, ".")
					path := filepath.Join("testdata", filename)

					board, err := NewBoard(rc.dimX, rc.dimY)
					if err != nil {
						t.Fatal(err)
					}

					_, err = board.Write(dc)
					if err != nil {
						t.Fatal(err)
					}

					got := board.Render(rc.tiles)
					if rc.armor {
						got = Armor(got)
					}

					if *updateGolden {
						err := os.WriteFile(path, []byte(got), os.ModePerm)
						if err != nil {
							t.Fatal("error updating golden file: ", err)
						}
						t.Log("updated golden file", path)
					}

					want, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}

					if !bytes.Equal(got, want) {
						t.Errorf("got %v want %v", got, string(want))
					}
				})
			}
		})
	}
}

var (
	updateGolden    = flag.Bool("update", false, "update the golden files of this test")
	clobberTestdata = flag.Bool("clobber", false, "clobber generated testdata")
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *clobberTestdata {
		testdata, err := filepath.Glob("testdata/*.txt")
		if err != nil {
			log.Fatal(err)
		}
		for _, tf := range testdata {
			log.Println("deleting", tf)
			err := os.Remove(tf)
			if err != nil {
				log.Println("ERROR: ", err)
			}
		}
	}

	os.Exit(m.Run())
}

func BenchmarkNewBoard(b *testing.B) {
	b.ReportAllocs()
	b.Run("10x10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewBoard(10, 10)
		}
	})
	b.Run("17x9", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewBoard(17, 9)
		}
	})
}

func BenchmarkBoard_Write(b *testing.B) {
	data := []byte{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed}
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()

	board, err := NewBoard(17, 9)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		board.Write(data)
	}
}

func BenchmarkBoard_Render(b *testing.B) {
	board, err := NewBoard(17, 9)
	if err != nil {
		b.Fatal(err)
	}

	data := []byte{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed}
	board.Write(data)

	for _, ts := range BundledTileSets {
		b.ReportAllocs()
		b.Run(ts.ID, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = board.Render(ts)
			}
		})
	}
}

func BenchmarkArmor(b *testing.B) {
	board, err := NewBoard(17, 9)
	if err != nil {
		b.Fatal(err)
	}

	data := []byte{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed}
	board.Write(data)
	raw := board.Render(SSHTiles)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Armor(raw)
	}
}
