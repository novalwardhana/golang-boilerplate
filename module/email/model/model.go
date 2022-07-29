package model

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

type Request struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Message struct {
	From        string
	To          string
	Cc          string
	Subject     string
	MIMEVersion string
	ContentType string
	Text        string
}
