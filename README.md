# randomart

Visual fingerprint hash (e.g. "randomart") library for Go.

Implements the ["drunken bishop"][1] algorithm from OpenSSH, with added support for
arbitrary grid size and tilesets.

[1]: http://www.dirk-loss.de/sshvis/drunken_bishop.pdf


## Output formats
Examples of rendering the same data with different settings.

### OpenSSH compatible
Dimensions: `17x9`, Tileset: `randomart.SSHTiles`, Armor: `true`
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
Dimensions: `10x10`, Tileset: `randomart.Galaxy`, Armor: `false`
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
