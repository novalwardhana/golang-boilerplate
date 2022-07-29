package model

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}
