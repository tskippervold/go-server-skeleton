package bambora

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://api.v1.checkout.bambora.com"
)

func (r CheckoutRequest) InitiateCheckout() (CheckoutRequestResponse, error) {
	url := urlWithPath("sessions")
	body, err := json.Marshal(r)
	if err != nil {
		return CheckoutRequestResponse{}, err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return CheckoutRequestResponse{}, err
	}

	res, err := executeRequest(req)
	if err != nil {
		return CheckoutRequestResponse{}, err
	}

	var response CheckoutRequestResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return CheckoutRequestResponse{}, err
	}

	return response, nil
}

func urlWithPath(path string) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseURL, path))
	return u
}

func executeRequest(req *http.Request) ([]byte, error) {
	// accesstoken@merchantnumber:secrettoken

	/*
	   ENCODED API KEY: NVIzQk1OQkM2R0JkVGdBY3dNbTk=
	   SECRET KEY: 5wR3rGPAmH4JARhou6zDcs3dtYWmz9K5Ho9iaauf
	   ACCESS KEY: 5R3BMNBC6GBdTgAcwMm9
	   MERCHANT NUMBER: T105856101
	   MD5 KEY: H8xKXwOrVh*/

	accessToken := "5R3BMNBC6GBdTgAcwMm9"
	merchantNumber := "T105856101"
	secretToken := "5wR3rGPAmH4JARhou6zDcs3dtYWmz9K5Ho9iaauf"

	token := fmt.Sprintf("%s@%s:%s", accessToken, merchantNumber, secretToken)
	token = base64.StdEncoding.EncodeToString([]byte(token))
	fmt.Println(token)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	/*if res.StatusCode > 299 {
		var meta Meta
		if err := json.Unmarshal(body, &meta); err != nil {
			return nil, err
		}

		return nil, errors.New(meta.Message.Merchant)
	}*/

	return body, nil
}
