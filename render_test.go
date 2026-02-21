package randomart

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestRenderTo_Golden(t *testing.T) {
	datacases := [][]byte{
		{},
		{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed},
	}

	rendercases := []struct {
		name       string
		tiles      TileSet
		dimX, dimY int
		opts       RenderOptions
	}{
		{
			name:  "openssh-17x9-border",
			tiles: OpenSSHTiles,
			dimX:  17,
			dimY:  9,
			opts:  RenderOptions{Border: true},
		},
		{
			name:  "openssh-17x9-border-header-footer",
			tiles: OpenSSHTiles,
			dimX:  17,
			dimY:  9,
			opts: RenderOptions{
				Border: true,
				Header: "ED25519 256",
				Footer: "SHA256",
			},
		},
		{
			name:  "openssh-17x9-border-header-trunc",
			tiles: OpenSSHTiles,
			dimX:  17,
			dimY:  9,
			opts: RenderOptions{
				Border: true,
				Header: "THIS HEADER IS INTENTIONALLY TOO LONG",
			},
		},
		{
			name:  "galaxy-10x10",
			tiles: GalaxyTiles,
			dimX:  10,
			dimY:  10,
			opts:  RenderOptions{},
		},
		{
			name:  "galaxy-10x10-header-footer",
			tiles: GalaxyTiles,
			dimX:  10,
			dimY:  10,
			opts: RenderOptions{
				Header: "GALAXY",
				Footer: "END",
			},
		},
		{
			name:  "galaxy-10x10-border-header-footer-trunc",
			tiles: GalaxyTiles,
			dimX:  10,
			dimY:  10,
			opts: RenderOptions{
				Border: true,
				Header: "GALAXY HEADER TOO LONG",
				Footer: "GALAXY FOOTER TOO LONG",
			},
		},
	}

	for _, dc := range datacases {
		slug := hex.EncodeToString(dc)
		if slug == "" {
			slug = "_empty"
		}
		t.Run(slug, func(t *testing.T) {

			for _, rc := range rendercases {
				specifier := rc.name
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

					var out bytes.Buffer
					renderOpts := rc.opts
					renderOpts.Tiles = rc.tiles
					_, err = RenderTo(&out, board, renderOpts)
					if err != nil {
						t.Fatal(err)
					}
					got := out.Bytes()

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
						t.Errorf("\ngot:\n%v\nwant:\n%v", string(got), string(want))
					}
				})
			}
		})
	}
}

func TestRenderTo_Errors(t *testing.T) {
	board, err := NewBoard(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("nil writer", func(t *testing.T) {
		n, err := RenderTo(nil, board, RenderOptions{})
		if !errors.Is(err, errNilWriter) {
			t.Fatalf("RenderTo() error = %v, want %v", err, errNilWriter)
		}
		if n != 0 {
			t.Fatalf("RenderTo() bytes = %v, want 0", n)
		}
	})
}

func Test_borderLineWithLabel(t *testing.T) {
	tests := []struct {
		name  string
		width int
		label string
		want  string
	}{
		{name: "empty label falls back to plain border", width: 17, label: "", want: "+-----------------+\n"},
		{name: "short centered label", width: 17, label: "SHA256", want: "+----[SHA256]-----+\n"},
		{name: "minimal width label", width: 3, label: "A", want: "+[A]+\n"},
		{name: "truncate single rune uses ellipsis", width: 3, label: "HELLO", want: "+[…]+\n"},
		{name: "truncate multi rune uses ellipsis suffix", width: 10, label: "123456789", want: "+[1234567…]+\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := borderLineWithLabel(tt.width, tt.label); got != tt.want {
				t.Fatalf("borderLineWithLabel(%d, %q) = %q, want %q", tt.width, tt.label, got, tt.want)
			}
		})
	}
}

func BenchmarkRenderTo(b *testing.B) {
	board, err := NewBoard(17, 9)
	if err != nil {
		b.Fatal(err)
	}

	data := []byte{0x9b, 0x4c, 0x7b, 0xce, 0x7a, 0xbd, 0x0a, 0x13, 0x61, 0xfb, 0x17, 0xc2, 0x06, 0x12, 0x0c, 0xed}
	board.Write(data)

	cases := []struct {
		name string
		opts RenderOptions
	}{
		{name: "default", opts: RenderOptions{}},
		{name: OpenSSHTiles.ID, opts: RenderOptions{Tiles: OpenSSHTiles}},
		{name: GalaxyTiles.ID, opts: RenderOptions{Tiles: GalaxyTiles}},
		{name: OpenSSHTiles.ID + "-border", opts: RenderOptions{Tiles: OpenSSHTiles, Border: true}},
	}

	for _, bc := range cases {
		b.Run(bc.name, func(b *testing.B) {
			for b.Loop() {
				_, err := RenderTo(io.Discard, board, bc.opts)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
