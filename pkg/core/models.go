package core

type SQLEntry struct {
	UUID  string      `json:"uuid"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	Flag  int         `json:"flag"`
}
