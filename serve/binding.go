package serve

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type J map[string]any

// bindJSON binds the request body (`r.Body`) to the object passed.
func bindJSON(r *http.Request, obj any) error {
	if r == nil || r.Body == nil {
		return errors.New("malformed body")
	}

	return DecodeJSON(r.Body, obj)
}

// writeJSON writes JSON to the response with specified status.
func writeJSON(w http.ResponseWriter, status int, data J) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DecodeJSON takes input from r and decodes it into the obj reference, or returns an error
// Helper function for json.Decoder
// TODO(RM): Validation: return validate(obj)
func DecodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return nil
}
