package model

type Person struct {
	ID      int    `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	Age     int    `json:"age" gorm:"age"`
	Address string `json:"address" gorm:"address"`
}

func (p *Person) TableName() string {
	return "persons"
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

type Pagination struct {
	Page         int      `json:"page"`
	Limit        int      `json:"limit"`
	TotalData    int      `json:"total_data"`
	NumberOfPage int      `json:"number_of_page"`
	Data         []Person `json:"data"`
}
