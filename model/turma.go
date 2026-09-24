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

type TurmaResponse struct {
	ID            string    `json:"id"`
	Nome          string    `json:"nome"`
	Disciplina    string    `json:"disciplina"`
	Professor     string    `json:"professor"`
	QtdMatriculos int       `json:"qtd_alunos_matriculados"`
	Alocada       bool      `json:"alocada"`
	Alocacao      *Alocacao `json:"alocacao,omitempty"`
}

func (t *Turma) ToResponse() TurmaResponse {
	return TurmaResponse{
		ID:            t.ID,
		Nome:          t.Nome,
		Disciplina:    t.Disciplina,
		Professor:     t.Professor,
		QtdMatriculos: len(t.Alunos),
		Alocada:       t.Alocacao != nil,
		Alocacao:      t.Alocacao,
	}
}
