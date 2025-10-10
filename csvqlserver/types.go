package csvqlserver

type SQLRequest struct {
	Query  string `json:"query"`
	Params []any  `json:"params,omitempty"`
}

type SQLResponse struct {
	Data  *QueryData `json:"data,omitempty"`
	Meta  *QueryMeta `json:"meta"`
	Error string     `json:"error,omitempty"`
}

type QueryData struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
}

type QueryMeta struct {
	RowCount int    `json:"rowCount"`
	Query    string `json:"query"`
	Params   []any  `json:"params"`
	Duration string `json:"duration,omitempty"`
}
