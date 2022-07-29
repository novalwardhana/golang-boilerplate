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

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

func (p *Person) TableName() string {
	return "persons"
}
