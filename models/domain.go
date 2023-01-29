package models


type Domain struct{
	Account string `json:"account"`
	URL string `json:"url" binding:"required"`
	Kind string `json:"kind" default:"Native"`
	Name string `json:"name"`
	Type string `json:"type" default:"Zone"`
	
} 

