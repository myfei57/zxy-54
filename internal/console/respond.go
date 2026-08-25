package console

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, code int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func readJSON(r *http.Request, value interface{}) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(value)
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}
