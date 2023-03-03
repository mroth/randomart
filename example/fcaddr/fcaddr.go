package main

import (
	"encoding/base32"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mroth/randomart"
)

var (
	tiles  = flag.String("tiles", "galaxy", "galaxy|openssh")
	width  = flag.Uint("x", 10, "width (cols)")
	height = flag.Uint("y", 10, "height (rows)")
)

// Protocol 1 addresses represent secp256k1 public encryption keys. The payload
// field contains the Blake2b-160 hash of the uncompressed public key (65 bytes).
//
// +------------+----------+----------------------------------------------------+
// |  network   | protocol |             base32 encoded data [39 bytes]         |
// |------------|----------|----------------------------------------------------|
// | 'f' or 't' |    '1'   |       'neiyfto7ozo4xaamg35jig7xbbrdpl6s7u6uy4i'    |
// +------------+----------|--------------------------------+-------------------|
//                         |              payload           |     checksum      |
//                         |             [20 byte]          |     [4 byte]      |
//                         |--------------------------------|-------------------|
//                         |        blake2b-160(PubKey)     |  blake2b-32(xxx)  |
//                         +--------------------------------+-------------------+
//
// Diagram based on https://spec.filecoin.io/appendix/address/, but clarified.

// uses lowercase version of standard RFC-4648 base32 encoding
var low32 = base32.
	NewEncoding("abcdefghijklmnopqrstuvwxyz234567").
	WithPadding(base32.NoPadding)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "usage: fcaddr <addr>")
		os.Exit(1)
	}
	addr := flag.Arg(0)

	if !strings.HasPrefix(addr, "f1") {
		fmt.Fprintln(os.Stderr, "only f1 addresses supported for this demo")
		os.Exit(3)
	}

	payload := strings.TrimSpace(addr[2:])
	fmt.Printf("payload: [%v] (%v byte string)\n", payload, len(payload))
	bs, err := low32.DecodeString(payload)
	if err != nil {
		log.Fatal(err)
	}
	if len(bs) != 24 {
		log.Fatalf("expected 24 bytes of payload, got %v", len(bs))
	}

	hsh, checksum := bs[:20], bs[20:]
	fmt.Printf("decoded: 0x%x (%v bytes) / checksum: 0x%x (%v bytes)\n", hsh, len(hsh), checksum, len(checksum))

	board := randomart.NewBoard(int(*width), int(*height))
	board.Write((hsh))
	fmt.Print(board.RenderString(randomart.GalaxyTiles))
}
