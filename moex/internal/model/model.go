package model

type ISSResponse struct {
	Securities Block `json:"securities"`
	Marketdata Block `json:"marketdata"`
	//Dataversion Block `json:"dataversion"`
}

type Block struct {
	//Metadata map[string]ColumnMetadata `json:"metadata"`
	Columns []string        `json:"columns"`
	Data    [][]interface{} `json:"data"`
}

type ColumnMetadata struct {
	Type    string `json:"type"`
	Bytes   int    `json:"bytes"`
	MaxSize int    `json:"max_size"`
}
