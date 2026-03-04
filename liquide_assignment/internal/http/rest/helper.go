package rest

import (
	"encoding/json"
	"net/http"
)

func getRequestDetails(r *http.Request) *json.Decoder {
	return json.NewDecoder(r.Body)
}
