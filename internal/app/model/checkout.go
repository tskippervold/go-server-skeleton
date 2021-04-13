package model

type CheckoutOrderLine struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int64  `json:"price"`
	Tax      int64  `json:"tax"`
}

type CheckoutOrder struct {
	ID         string               `json:"id"`
	Amount     int64                `json:"amount"`
	TaxAmount  int64                `json:"tax_amount"`
	OrderLines *[]CheckoutOrderLine `json:"order_lines"`
}

type CheckoutInit struct {
	Currency string        `json:"currency"`
	Country  string        `json:"country"`
	Locale   string        `json:"locale"`
	Order    CheckoutOrder `json:"order"`
}
