package tiles

import (
	"fmt"
	"bytes"
)

func init() {
	Register(new(HelloTile))
}

type HelloTile struct {};

func (h *HelloTile) Name() string {
	return "hello"
}

func (h *HelloTile) RenderFunc(renderer *TileRenderer, arg ... interface{}) string {
	str := bytes.NewBufferString("Welcome aboard")
	
	if arg[0] != nil {
		fmt.Fprintf(str, ", %s", arg[0])
	}
		
	return str.String() 
}
