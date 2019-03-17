package api

type StoreRequest struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
}

type Response struct {
	StatusCode    int         `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Result        interface{} `json:"result,omitempty"`
	Errors        []string    `json:"errors,omitempty"`
}

type IdKeyPair struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}
