package tiles

var Registered = make([]Tile, 0, 10)

func Register(tile Tile) {
	Registered = append(Registered, tile)
}

func init() {
}