package models

type Response struct{
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type NSResponse struct{
	CurrentNS1 string `json:"current_ns1"`
	CurrentNS2 string `json:"current_ns2"`
	CorrectNS1 string `json:"correct_ns1"`
	CorrectNS2 string `json:"correct_ns2"`
}

