package respond

import (
	"encoding/json"
	"net/http"
)

// Externals

func InternalError(w http.ResponseWriter) {
	Error(w, "server_error", http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter) {
	Error(w, "Bad request", http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter) {
	Error(w, "unauthorized", http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter) {
	Error(w, "forbidden", http.StatusForbidden)
}

func Error(w http.ResponseWriter, error string, code int) {
	writeHeader(w, code)

	if err := writeJson(w, map[string]string{errorFieldName: error}); err != nil {
		panic(err)
	}
}

func Ok(w http.ResponseWriter, v interface{}) {
	writeHeader(w, http.StatusOK)

	if err := writeJson(w, v); err != nil {
		panic(err)
	}
}

func Created(w http.ResponseWriter, v interface{}) {
	writeHeader(w, http.StatusCreated)

	if err := writeJson(w, v); err != nil {
		panic(err)
	}
}

// Internals

var (
	errorFieldName = "error"
)

func writeHeader(w http.ResponseWriter, code int) {
	// Per spec, UTF-8 is the default, and the charset parameter should not
	// be necessary. But some clients (eg: Chrome) think otherwise.
	// Since json.Marshal produces UTF-8, setting the charset parameter is a
	// safe option.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)
}

func encodeJson(w http.ResponseWriter, v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Encode the object in JSON and call Write.
func writeJson(w http.ResponseWriter, v interface{}) error {
	if v == nil {
		return nil
	}

	b, err := encodeJson(w, v)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// Provided in order to implement the http.ResponseWriter interface.
func write(w http.ResponseWriter, b []byte) (int, error) {
	return w.Write(b)
}
