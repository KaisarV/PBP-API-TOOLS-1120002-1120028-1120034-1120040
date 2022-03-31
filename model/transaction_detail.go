package model

type TransactionDetail struct {
	ID       int     `json:"id"`
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type TransactionDetailsResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []TransactionDetail `json:"data"`
}

type TransactionDetailResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    TransactionDetail `json:"data"`
}
