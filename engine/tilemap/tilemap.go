package tilemap

type TileType uint8

type Tile struct {
	Type TileType
}

type Tilemap struct {
	TileSize int // In Pixels
	tiles    [][]Tile
}

func New(tiles [][]Tile, tileSize int) *Tilemap {
	return &Tilemap{tileSize, tiles}
}

func (tilemap *Tilemap) Width() int {
	return len(tilemap.tiles)
}

func (tilemap *Tilemap) Height() int {
	return len(tilemap.tiles[0])
}

func (Tilemap *Tilemap) Get(x int, y int) (Tile, bool) {
	if x < 0 || x >= len(Tilemap.tiles) || y < 0 || y >= len(Tilemap.tiles[x]) {
		return Tile{}, false
	}

	return Tilemap.tiles[x][y], true
}
