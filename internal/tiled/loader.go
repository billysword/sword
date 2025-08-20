package tiled

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// --- TSX (Tileset) ---

type TSXTileset struct {
	XMLName    xml.Name     `xml:"tileset"`
	Name       string       `xml:"name,attr"`
	TileWidth  int          `xml:"tilewidth,attr"`
	TileHeight int          `xml:"tileheight,attr"`
	TileCount  int          `xml:"tilecount,attr"`
	Columns    int          `xml:"columns,attr"`
	Image      TSXImage     `xml:"image"`
	Tiles      []TSXTile    `xml:"tile"`
}

type TSXImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type TSXTile struct {
	ID         int                `xml:"id,attr"`
	Properties []TSXTileProperty  `xml:"properties>property"`
}

type TSXTileProperty struct {
	Name  string  `xml:"name,attr"`
	Type  string  `xml:"type,attr"`
	Value string  `xml:"value,attr"`
}

// --- TMJ (Map) ---

type TMJMap struct {
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	TileWidth  int           `json:"tilewidth"`
	TileHeight int           `json:"tileheight"`
	Layers     []TMJLayer    `json:"layers"`
	Tilesets   []TMJTileset  `json:"tilesets"`
}

type TMJTileset struct {
	FirstGID int    `json:"firstgid"`
	Source   string `json:"source"`
}

type TMJLayer struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Visible   bool             `json:"visible"`
	Opacity   float64          `json:"opacity"`
	Width     int              `json:"width"`
	Height    int              `json:"height"`
	Data      []uint32         `json:"data"`
	Objects   []TMJObject      `json:"objects"`
}

type TMJObject struct {
	ID         int                 `json:"id"`
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	X          float64             `json:"x"`
	Y          float64             `json:"y"`
	Width      float64             `json:"width"`
	Height     float64             `json:"height"`
	Properties []TMJObjectProperty `json:"properties"`
}

type TMJObjectProperty struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// --- Resolved structures ---

type TileProperties struct {
	Kind    string
	Solid   bool
	OneWay  string
	SlopeM  float64
	SlopeB  float64
}

type Portal struct {
	Name      string
	RectPx    [4]float64 // x, y, w, h
	ToZone    string
	ToRoom    string
	ToPortal  string
}

type LoadedTileset struct {
	FirstGID  int
	TSX       TSXTileset
	ByGID     map[uint32]TileProperties
}

type LoadedMap struct {
	TMJ            TMJMap
	Tilesets       []LoadedTileset
	RenderLayer    *TMJLayer
	CollisionLayer *TMJLayer
	Portals        []Portal
}

// LoadMap loads a .tmj map and all referenced .tsx tilesets.
func LoadMap(path string) (*LoadedMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m TMJMap
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}

	// Load tilesets and build gid map
	loaded := make([]LoadedTileset, 0, len(m.Tilesets))
	for _, ts := range m.Tilesets {
		tsPath := resolveRelative(path, ts.Source)
		tsFile, err := os.Open(tsPath)
		if err != nil {
			return nil, fmt.Errorf("open tileset %s: %w", tsPath, err)
		}
		var tsx TSXTileset
		if err := xml.NewDecoder(tsFile).Decode(&tsx); err != nil {
			tsFile.Close()
			return nil, fmt.Errorf("decode tileset %s: %w", tsPath, err)
		}
		tsFile.Close()

		byGID := make(map[uint32]TileProperties)
		for _, t := range tsx.Tiles {
			props := TileProperties{}
			for _, p := range t.Properties {
				switch p.Name {
				case "kind":
					props.Kind = p.Value
				case "solid":
					props.Solid = strings.EqualFold(p.Value, "true")
				case "one_way":
					props.OneWay = p.Value
				case "slopeM":
					props.SlopeM = parseFloat(p.Value)
				case "slopeB":
					props.SlopeB = parseFloat(p.Value)
				}
			}
			gid := uint32(ts.FirstGID + t.ID)
			byGID[gid] = props
		}
		loaded = append(loaded, LoadedTileset{FirstGID: ts.FirstGID, TSX: tsx, ByGID: byGID})
	}
	// Sort tilesets by FirstGID ascending to make gid lookup predictable
	sort.Slice(loaded, func(i, j int) bool { return loaded[i].FirstGID < loaded[j].FirstGID })

	lm := &LoadedMap{TMJ: m, Tilesets: loaded}
	for i := range m.Layers {
		layer := &m.Layers[i]
		switch layer.Name {
		case "render":
			lm.RenderLayer = layer
		case "collision":
			lm.CollisionLayer = layer
		case "portals":
			lm.Portals = extractPortals(layer)
		}
	}
	return lm, nil
}

func extractPortals(layer *TMJLayer) []Portal {
	var out []Portal
	for _, obj := range layer.Objects {
		if strings.ToLower(obj.Type) != "portal" {
			continue
		}
		p := Portal{Name: obj.Name, RectPx: [4]float64{obj.X, obj.Y, obj.Width, obj.Height}}
		for _, prop := range obj.Properties {
			switch prop.Name {
			case "toZone":
				p.ToZone, _ = prop.Value.(string)
			case "toRoom":
				p.ToRoom, _ = prop.Value.(string)
			case "toPortal":
				p.ToPortal, _ = prop.Value.(string)
			}
		}
		out = append(out, p)
	}
	return out
}

// Tiled stores flip flags in the top 3 bits of the GID. NormalizeGID clears those flags.
const (
	gidMask uint32 = 0x1FFFFFFF
)

// NormalizeGID removes flip flags and returns the raw tile id used for tileset indexing.
func NormalizeGID(gid uint32) uint32 {
	return gid & gidMask
}

// Resolve tile properties from a global id. Zero gid returns zero-value props.
func (lm *LoadedMap) PropertiesForGID(gid uint32) (TileProperties, bool) {
	gid = NormalizeGID(gid)
	if gid == 0 {
		return TileProperties{}, false
	}
	for _, ts := range lm.Tilesets {
		if gid < uint32(ts.FirstGID) {
			break
		}
		if props, ok := ts.ByGID[gid]; ok {
			return props, true
		}
	}
	return TileProperties{}, false
}

func resolveRelative(basePath, rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Clean(filepath.Join(filepath.Dir(basePath), rel))
}

func parseFloat(s string) float64 {
	f, err := strconvParseFloat(s)
	if err != nil {
		return math.NaN()
	}
	return f
}

// local wrapper to keep imports minimal if strconv is not preferred in this snippet
func strconvParseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscan(s, &f)
	return f, err
}

// Utility to iterate the collision grid as booleans using either explicit collision layer or tile properties.
func (lm *LoadedMap) IsSolidAt(index int) bool {
	if lm.CollisionLayer != nil && index >= 0 && index < len(lm.CollisionLayer.Data) {
		// Treat value 1 as definitely solid
		// Value 2 needs context - it's used for boundaries that might be one-way or special
		// For now, treat both 1 and 2 as solid to be safe, but this could be refined
		return lm.CollisionLayer.Data[index] == 1 || lm.CollisionLayer.Data[index] == 2
	}
	if lm.RenderLayer != nil && index >= 0 && index < len(lm.RenderLayer.Data) {
		gid := NormalizeGID(lm.RenderLayer.Data[index])
		if props, ok := lm.PropertiesForGID(gid); ok {
			return props.Solid
		}
	}
	return false
}

// GetCollisionValue returns the raw collision value at the given index.
// This allows for more nuanced collision handling when needed.
func (lm *LoadedMap) GetCollisionValue(index int) uint32 {
	if lm.CollisionLayer != nil && index >= 0 && index < len(lm.CollisionLayer.Data) {
		return lm.CollisionLayer.Data[index]
	}
	return 0
}