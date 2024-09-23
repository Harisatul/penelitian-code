package shared

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(StandardResponse{
		Message: err.Error(),
		Data:    nil,
	})
}

func WriteInternalServerErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(StandardResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	})
}

func WriteSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(StandardResponse{
		Data: data,
	})
}
