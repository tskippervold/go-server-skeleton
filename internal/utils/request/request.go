package request

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, t interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&t)
}
