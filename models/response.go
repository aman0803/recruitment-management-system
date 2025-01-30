package models

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Jwt     *string     `json:"jwt"`
	Data    interface{} `json:"data"`
}
