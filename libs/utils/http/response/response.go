package response

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")

	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	return json.NewEncoder(w).Encode(v)
}
