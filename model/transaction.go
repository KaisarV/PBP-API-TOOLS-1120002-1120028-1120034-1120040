package model

type Transaction struct {
	ID        int `json:"id"`
	UserID    int `json:"userid"`
	ProductId int `json:"productid"`
	Quantity  int `json:"qty"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type TransactionResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    Transaction `json:"data"`
}
