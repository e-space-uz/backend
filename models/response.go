package models

type FailureResponse struct {
	Success bool   `json:"success"`
	Error   error  `json:"error"`
	Message string `json:"message"`
}

type GetAllResponse struct {
	Data  []interface{} `json:"data"`
	Count int64         `json:"count"`
}
