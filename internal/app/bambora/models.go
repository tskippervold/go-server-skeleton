package bambora

type MetaMessage struct {
	EndUser  string `json:"enduser"`
	Merchant string `json:"merchant"`
}

type Meta struct {
	Result  bool        `json:"result"`
	Message MetaMessage `json:"message"`
}

type URL struct {
	Accept           string `json:"accept"`
	Cancel           string `json:"cancel"`
	RedirectOnAccept int    `json:"immediateredirecttoaccept"` // If set to 1 the Checkout will redirect to the defined accept URL immediately after the payment is completed. If set to a value higher than 1 it will delay the redirect for the defined value in seconds.
}

type Order struct {
	ID        string `json:"id"`
	Amount    int64  `json:"amount"`
	VATAmount int64  `json:"vatamount"`
	Currency  string `json:"currency"`
}

type CheckoutRequest struct {
	Order Order `json:"order"`
	URL   URL   `json:"url"`
}

type CheckoutRequestResponse struct {
	Token *string `json:"token"`
	URL   *string `json:"url"`
	Meta  Meta    `json:"meta"`
}
