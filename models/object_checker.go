package models

type Text struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

type RequestTask struct {
	Props []string `json:"texts"`
}
