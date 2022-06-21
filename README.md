# bite - a library for parsing binaries with Go

```
type PNGChunk struct {
    Type    string
    Chunk   []byte
    CRC     []byte
}



ok, err := bite.Chew(
    `0x89, "PNG\r\n", 0x1A, "\n", Chunks[Len:32/integer, Type:32, Chunk:Len, CRC:32]:...`,
    reader,
    func(Chunks []PNGChunk) error {
        // do something
    }
)
```
