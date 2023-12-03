package util

type Msg struct {
	SourceFile string
	DstFile    string
	Status     string
}

type SpineItem struct {
	Id     int    `json:"id"`
	Length int    `json:"length"`
	Src    string `json:"src"`
}

type Spine []SpineItem

type Directory struct {
	MakeTime string                 `json:"make-time"`
	Metadata map[string]interface{} `json:"metadata"`
	Images   []string               `json:"images"`
	Css      string                 `json:"css"`
	Js       string                 `json:"js"`
	Catalog  []interface{}          `json:"catalog"`
	Spine    Spine                  `json:"spine"`
}