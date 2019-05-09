package types

type (
	SearchResult struct {
		TextSearch   string `json:"text_search"`
		TookTime     int64  `json:"took_time(milliseconds)"`
		NumberResult int64  `json:"number_result"`
		Result       []Film `json:"result"`
	}
)
