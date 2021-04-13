package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	env "github.com/tskippervold/golang-base-server/internal/app"
	"github.com/tskippervold/golang-base-server/internal/app/bambora"
	"github.com/tskippervold/golang-base-server/internal/app/model"
	"github.com/tskippervold/golang-base-server/internal/utils/handler"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
	"github.com/tskippervold/golang-base-server/internal/utils/request"
	"github.com/tskippervold/golang-base-server/internal/utils/respond"
)

func CheckoutHandlers(r *mux.Router, env *env.Env) {
	r.Handle("/order", checkoutOrder(env)).Methods("POST")
	r.Handle("/order/callback/{status}", handleOrderCallback(env)).Methods("GET")

	r.Handle("/init", initializeCheckout(env)).Methods("POST")
}

func initializeCheckout(env *env.Env) handler.Handler {
	type Request struct {
		model.CheckoutInit
	}

	type Response struct {
		model.CheckoutInit
		HTML string `json:"html"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		logger := log.ForRequest(r)

		var requestBody Request
		request.Decode(r.Body, &requestBody)

		if requestBody.Order.ID == "" {
			err := errors.New("Missing or invalid order.id")
			return respond.Error(err, http.StatusBadRequest, err.Error(), "invalid_orderid")
		}

		/*
			Prevent tampering of order when paying?
			Create a signed token of the digest of `requestBody`.
			Include that token as a query param in the source.

			Pass a token or some encoded string containing the order
			to the source. Source should then extract and populate this to the
			payment providers.

			Create the customer and order in our own database?
			Relate that to the merchant?
			Can we store this?
		*/

		checkoutRequest := bambora.CheckoutRequest{
			Order: bambora.Order{
				ID:        requestBody.Order.ID,
				Amount:    requestBody.Order.Amount,
				VATAmount: requestBody.Order.TaxAmount,
				Currency:  requestBody.Currency,
			},
			URL: bambora.URL{
				Accept:           "http://localhost:3000/api/checkout/order/callback/accept",
				Cancel:           "http://localhost:3000/api/checkout/order/callback/cancel",
				RedirectOnAccept: 0,
			},
		}

		logger.Debug("Initiating Bambora checkout request for order", requestBody.Order.ID)
		checkoutResponse, err := checkoutRequest.InitiateCheckout()
		if err != nil {
			return respond.GenericServerError(err)
		}

		if checkoutResponse.Meta.Result == false {
			err = errors.New(checkoutResponse.Meta.Message.Merchant)
			return respond.Error(err, http.StatusBadRequest, checkoutResponse.Meta.Message.EndUser, "bad_request")
		}

		logger.Debug("Created Bambora checkout token", *checkoutResponse.Token, *checkoutResponse.URL)
		checkoutURL, err := url.Parse("http://localhost:3001/checkout-frame")
		if err != nil {
			return respond.GenericServerError(err)
		}

		query := checkoutURL.Query()
		query.Set("bambora_token", *checkoutResponse.Token)
		checkoutURL.RawQuery = query.Encode()

		htmlHeight := "700"
		html := `<iframe id="zo-checkout-iframe" name="zo-checkout-iframe" scrolling="no" frameborder="0" style="height: ` + htmlHeight + `px; width: 1px; min-width: 100%;" src="` + checkoutURL.String() + `"></iframe>`

		response := Response{requestBody.CheckoutInit, html}

		return respond.Success(http.StatusCreated, response)
	})
}

func checkoutOrder(env *env.Env) handler.Handler {
	type Customer struct {
		ID          int64  `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}

	type Request struct {
		OrderID     int64    `json:"order_id"`
		TotalAmount int64    `json:"total_amount"`
		TotalVAT    int64    `json:"total_vat"`
		Currency    string   `json:"currency"`
		Customer    Customer `json:"customer"`
		CallbackURL string   `json:"callback_url"`
	}

	type Response struct {
		OrderID     string `json:"order_id"`
		CheckoutURI string `json:"checkout_uri"`
	}

	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		logger := log.ForRequest(r)

		id := uuid.NewV4().String()

		var requestBody Request
		request.Decode(r.Body, &requestBody)

		cbURLEncoded := base64.URLEncoding.EncodeToString([]byte(requestBody.CallbackURL))

		checkoutRequest := bambora.CheckoutRequest{
			Order: bambora.Order{
				ID:        fmt.Sprintf("%d", requestBody.OrderID),
				Amount:    requestBody.TotalAmount,
				VATAmount: requestBody.TotalVAT,
				Currency:  requestBody.Currency,
			},
			URL: bambora.URL{
				Accept:           fmt.Sprintf("http://localhost:3000/api/checkout/order/callback/accept?zo_cb=%s&order_id=%d", cbURLEncoded, requestBody.OrderID),
				Cancel:           fmt.Sprintf("http://localhost:3000/api/checkout/order/callback/cancel?zo_cb=%s&order_id=%d", cbURLEncoded, requestBody.OrderID),
				RedirectOnAccept: 1,
			},
		}

		checkoutResponse, err := checkoutRequest.InitiateCheckout()
		if err != nil {
			return respond.GenericServerError(err)
		}

		logger.Info("Checkout for order", requestBody.OrderID)

		return respond.Success(http.StatusOK, Response{
			OrderID:     id,
			CheckoutURI: *checkoutResponse.URL,
		})
	})
}

func handleOrderCallback(env *env.Env) handler.Handler {
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *respond.Response {
		//logger := log.ForRequest(r)

		/*
			?txnid=205888692085291008
			&orderid=125
			&reference=419080560448
			&amount=3751250
			&currency=NOK
			&date=20200820
			&time=1352
			&feeid=309812
			&txnfee=0
			&paymenttype=5
			&cardno=415421XXXXXX0001
			&eci=7
			&issuercountry=DNK
			&hash=b73b96c5d27a419971e4fb11f9e02cd3
		*/

		vars := mux.Vars(r)
		query := r.URL.Query()

		status := strings.ToUpper(vars["status"])
		orderID := query.Get("order_id")
		cbEncoded := query.Get("zo_cb")

		if len(status) <= 0 || len(orderID) <= 0 || len(cbEncoded) <= 0 {
			err := errors.New("Invalid callback. Missing required query params.")
			return respond.Error(err, http.StatusBadRequest, "Needs `status`, `order_id` and `zo_cb`", "error")
		}

		cbDecoded, err := base64.URLEncoding.DecodeString(cbEncoded)
		if err != nil {
			return respond.GenericServerError(err)
		}

		callbackURL, err := url.Parse(string(cbDecoded))
		if err != nil {
			return respond.GenericServerError(err)
		}

		callbackQuery := callbackURL.Query()
		callbackQuery.Add("zo_order_id", orderID)

		switch status {
		case "ACCEPT":
			callbackQuery.Add("zo_status", "ok")
			callbackQuery.Add("txnid", query.Get("txnid"))

		default: // CANCEL or other
			callbackQuery.Add("zo_status", "fail")
		}

		callbackURL.RawQuery = callbackQuery.Encode()
		http.Redirect(w, r, callbackURL.String(), http.StatusFound)

		return nil
	})
}
