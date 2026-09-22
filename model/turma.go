package model

type Turma struct {
	ID         string    `json:"id"`
	Nome       string    `json:"nome"`
	Professor  string    `json:"professor"`
	Disciplina string    `json:"disciplina"`
	Alunos     []string  `json:"alunos"`
	Alocacao   *Alocacao `json:"alocacao,omitempty"`
}

type Alocacao struct {
	SalaID     string `json:"sala_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}
