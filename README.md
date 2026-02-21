# randomart

Visual fingerprint hash (e.g. "randomart") library for Go.

Implements the ["drunken bishop"][1] algorithm from OpenSSH, with added support for
arbitrary grid size and tile sets.

[1]: http://www.dirk-loss.de/sshvis/drunken_bishop.pdf


## Output formats
Examples of rendering the same data with different settings.

### OpenSSH compatible
```go
randomart.RenderOptions{
	Tiles:  randomart.OpenSSHTiles,
	Border: true,
	Header: "ED25519 256",
	Footer: "SHA256",
}
```
```
+--[ED25519 256]--+
|    .+.          |
|      o.         |
|     .. +        |
|      Eo =       |
|        S + .    |
|       o B . .   |
|        B o..    |
|         *...    |
|        .o+...   |
+----[SHA256]-----+
```

### Spacey emoji
```go
randomart.RenderOptions{Tiles: randomart.GalaxyTiles}
```
```
🌒🌔🌑🌑🌑🌑🌑🌑🌑🌑
🌑🌑🌔🌑🌑🌑🌑🌑🌑🌑
🌑🌒🌑🌒🌑🌑🌑🌑🌑🌑
🌑🌑🌚🌑🌔🌑🌑🌑🌑🌑
🌑🌑🌑🌓🌑🌕🌑🌑🌑🌑
🌑🌑🌑🌑🌓🌝🌔🌑🌒🌑
🌑🌑🌑🌑🌓🌔🌔🌒🌑🌒
🌑🌑🌑🌑🌓🌕🌒🌒🌓🌑
🌑🌑🌑🌑🌑🌓🌔🌓🌑🌒
🌑🌑🌑🌑🌑🌒🌔🌔🌒🌒
```
 
 ## Examples
 
 * [fcaddr](./example/fcaddr/): Fingerprint Filecoin f1 addresses
