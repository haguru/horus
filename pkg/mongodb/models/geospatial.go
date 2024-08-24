package models

type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type LineStringAndMultipoint struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type Polygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type GeoJSON struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}
