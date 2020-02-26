package request

import (
	"encoding/json"
	"io"
)

func Decode(r io.ReadCloser, t interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&t)
	r.Close()
	return err
}
