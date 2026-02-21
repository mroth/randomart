# randomart

Visual fingerprint hash (e.g. "randomart") library for Go.

Implements the ["drunken bishop"][1] algorithm from OpenSSH, with added support for
arbitrary grid size and tile sets.

[1]: http://www.dirk-loss.de/sshvis/drunken_bishop.pdf


## Output formats
Examples of rendering the same data with different settings.

### OpenSSH compatible
```go
randomart.RenderOptions{Tiles: randomart.OpenSSHTiles, Border: true}
```
```
+-----------------+
|    .+.          |
|      o.         |
|     .. +        |
|      Eo =       |
|        S + .    |
|       o B . .   |
|        B o..    |
|         *...    |
|        .o+...   |
+-----------------+
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

## Rendering API

Render to any `io.Writer`:

```go
board, _ := randomart.NewBoard(17, 9)
_, _ = board.Write(fingerprint)

// defaults to randomart.OpenSSHTiles when Tiles is empty
_, _ = randomart.RenderTo(os.Stdout, board, randomart.RenderOptions{})

// choose tiles + border
_, _ = randomart.RenderTo(os.Stdout, board,
		randomart.RenderOptions{Tiles: randomart.GalaxyTiles, Border: true},
)
```
