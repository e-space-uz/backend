package models

type FailureResponse struct {
	Success bool   `json:"success"`
	Error   error  `json:"error"`
	Message string `json:"message"`
}
