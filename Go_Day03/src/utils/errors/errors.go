package errors

import "net/http"

// ErrorResponse myError in format JSON
type ErrorResponse struct {
	Err string `json:"error"`
}

// getter поля Err структуры ErrorResponse
func (e ErrorResponse) Error() string {
	return e.Err
}

// HandleError Обработчик ошибок
func HandleError(w http.ResponseWriter, errorMsg string, statusCode int) {
	http.Error(w, errorMsg, statusCode)
}
