package model

type Sala struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int `json:"capacidade" binding:"gt=0"`
	Recursos   []string `json:"recursos"`
	
}
