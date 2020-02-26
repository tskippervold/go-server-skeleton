package respond

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	Status  int
	Failure *Failure
	Success interface{}
}

type Failure struct {
	Cause   error  `json:"-"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (r *Response) Write(w http.ResponseWriter) {
	body, err := json.Marshal(r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Per spec, UTF-8 is the default, and the charset parameter should not
	// be necessary. But some clients (eg: Chrome) think otherwise.
	// Since json.Marshal produces UTF-8, setting the charset parameter is a
	// safe option.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.Status)
	w.Write(body)
}

func (r *Response) MarshalJSON() ([]byte, error) {
	if f := r.Failure; f != nil {
		// Maps the failure to a json object:
		// { "error": { <Failure> + debug string } }
		return json.Marshal(&struct {
			Error interface{} `json:"error"`
		}{
			Error: struct {
				*Failure
				Debug string `json:"debug"`
			}{
				Failure: f,
				Debug:   f.Cause.Error(),
			},
		})
	}

	if s := r.Success; s != nil {
		return json.Marshal(s)
	}

	return json.Marshal(struct{}{})
}

func Success(status int, i interface{}) *Response {
	return &Response{
		Status:  status,
		Failure: nil,
		Success: i,
	}
}

func Error(err error, status int, detail string, code string) *Response {
	f := Failure{
		Cause:   err,
		Message: detail,
		Code:    code,
	}

	return &Response{
		Status:  status,
		Failure: &f,
		Success: nil,
	}
}

func GenericServerError(err error) *Response {
	status := http.StatusInternalServerError

	f := Failure{
		Cause:   err,
		Message: http.StatusText(status),
		Code:    "error",
	}

	return &Response{
		Status:  status,
		Failure: &f,
		Success: nil,
	}
}
