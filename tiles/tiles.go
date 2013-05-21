package tiles

import (
	"log"
	"text/template"
)

type TileRenderFunc func(renderer *TileRenderer, arg ...interface{}) string

type Tile interface {
	Name() string
	RenderFunc(renderer *TileRenderer, arg ...interface{}) string
}

type TileRenderer struct {
	Templ   *template.Template
	Tiles   []*Tile
	Globals map[string]string
	tname string
}

func NewRenderer(filename string) *TileRenderer {
	renderer := &TileRenderer{Templ: template.New(filename), tname: filename}
	
	for t := range Registered {
		renderer.UsesTile(Registered[t])
	}

	log.Printf("Creating renderer")
	
	renderer.Parse()

	return renderer
}

func (renderer *TileRenderer) UsesTile(tile Tile) *TileRenderer {
	renderer.Tiles = append(renderer.Tiles, &tile)

	funcMap := template.FuncMap{
		tile.Name(): func(args ...interface{}) string {
			return tile.RenderFunc(renderer, args...)
		}}

	log.Printf("Adding template render functions %v", funcMap)
	renderer.Templ.Funcs(funcMap)
	
	return renderer
}

func (renderer *TileRenderer) Parse() {
	template.Must(renderer.Templ.ParseGlob(renderer.tname))
}
