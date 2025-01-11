package utils

import (
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}