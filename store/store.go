package store

import (
	"errors"
	"sync"
	"time"

	"api-gin/model"
)

var (
	ErrSalaJaExiste  = errors.New("já existe uma sala com este id")
	ErrAlunoJaExiste = errors.New("já existe um aluno com esta matrícula")
	ErrTurmaJaExiste = errors.New("já existe uma turma com este id")

	ErrSalaNaoEncontrada  = errors.New("sala não encontrada")
	ErrAlunoNaoEncontrado = errors.New("aluno não encontrado")
	ErrTurmaNaoEncontrada = errors.New("turma não encontrada")

	ErrAlunoJaMatriculado = errors.New("aluno já matriculado nesta turma")
	ErrCapacidadeExcedida = errors.New("capacidade da sala excedida")
	ErrConflitoDeSala     = errors.New("sala já alocada para outra turma neste horário")
	ErrConflitoDeAgenda   = errors.New("aluno já matriculado em turma com horário conflitante")
	ErrHorarioInvalido    = errors.New("horário inválido")
)

type Store struct {
	mu     sync.Mutex
	salas  map[string]*model.Sala
	alunos map[string]*model.Aluno
	turmas map[string]*model.Turma
}

func NewStore() *Store {
	return &Store{
		salas:  make(map[string]*model.Sala),
		alunos: make(map[string]*model.Aluno),
		turmas: make(map[string]*model.Turma),
	}
}

func (s *Store) CriarSala(sala *model.Sala) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.salas[sala.ID]; exists {
		return ErrSalaJaExiste
	}
	s.salas[sala.ID] = sala
	return nil
}

func (s *Store) ListarSalas() []*model.Sala {
	s.mu.Lock()
	defer s.mu.Unlock()
	var salas []*model.Sala
	for _, sala := range s.salas {
		salas = append(salas, sala)
	}
	return salas
}

func (s *Store) CriarAluno(aluno *model.Aluno) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.alunos[aluno.Matricula]; exists {
		return ErrAlunoJaExiste
	}
	s.alunos[aluno.Matricula] = aluno
	return nil
}

func (s *Store) ListarAlunos() []*model.Aluno {
	s.mu.Lock()
	defer s.mu.Unlock()
	var alunos []*model.Aluno
	for _, aluno := range s.alunos {
		alunos = append(alunos, aluno)
	}
	return alunos
}

func (s *Store) BuscarAluno(matricula string) (*model.Aluno, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	aluno, ok := s.alunos[matricula]
	if !ok {
		return nil, ErrAlunoNaoEncontrado
	}
	return aluno, nil
}

func (s *Store) CriarTurma(turma *model.Turma) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.turmas[turma.ID]; exists {
		return ErrTurmaJaExiste
	}
	turma.Alunos = []string{}
	s.turmas[turma.ID] = turma
	return nil
}

func (s *Store) ListarTurmas() []model.TurmaResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	turmas := make([]model.TurmaResponse, 0, len(s.turmas))
	for _, turma := range s.turmas {
		turmas = append(turmas, turma.ToResponse())
	}
	return turmas
}

func (s *Store) ListarAlunosDaTurma(turmaID string) ([]*model.Aluno, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	turma, ok := s.turmas[turmaID]
	if !ok {
		return nil, ErrTurmaNaoEncontrada
	}
	alunos := make([]*model.Aluno, 0, len(turma.Alunos))
	for _, matricula := range turma.Alunos {
		if aluno, ok := s.alunos[matricula]; ok {
			alunos = append(alunos, aluno)
		}
	}
	return alunos, nil
}

func (s *Store) MatricularAluno(turmaID, alunoID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	turma, ok := s.turmas[turmaID]
	if !ok {
		return ErrTurmaNaoEncontrada
	}
	if _, ok := s.alunos[alunoID]; !ok {
		return ErrAlunoNaoEncontrado
	}
	for _, matricula := range turma.Alunos {
		if matricula == alunoID {
			return ErrAlunoJaMatriculado
		}
	}

	if turma.Alocacao != nil {
		if sala, ok := s.salas[turma.Alocacao.SalaID]; ok {
			if len(turma.Alunos)+1 > sala.Capacidade {
				return ErrCapacidadeExcedida
			}
		}
		for _, outra := range s.turmas {
			if outra.ID == turma.ID || outra.Alocacao == nil {
				continue
			}
			if !contem(outra.Alunos, alunoID) {
				continue
			}
			if conflita(turma.Alocacao, outra.Alocacao) {
				return ErrConflitoDeAgenda
			}
		}
	}

	turma.Alunos = append(turma.Alunos, alunoID)
	return nil
}

func (s *Store) AlocarSala(turmaID string, aloc *model.Alocacao) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	turma, ok := s.turmas[turmaID]
	if !ok {
		return ErrTurmaNaoEncontrada
	}
	sala, ok := s.salas[aloc.SalaID]
	if !ok {
		return ErrSalaNaoEncontrada
	}
	if _, _, err := paraIntervalo(aloc.HoraInicio, aloc.HoraFim); err != nil {
		return err
	}
	if len(turma.Alunos) > sala.Capacidade {
		return ErrCapacidadeExcedida
	}
	for _, outra := range s.turmas {
		if outra.ID == turma.ID || outra.Alocacao == nil {
			continue
		}
		if outra.Alocacao.SalaID != aloc.SalaID {
			continue
		}
		if conflita(aloc, outra.Alocacao) {
			return ErrConflitoDeSala
		}
	}

	turma.Alocacao = aloc
	return nil
}

func contem(lista []string, id string) bool {
	for _, v := range lista {
		if v == id {
			return true
		}
	}
	return false
}

func conflita(a, b *model.Alocacao) bool {
	if a.DiaSemana != b.DiaSemana {
		return false
	}
	iniA, fimA, e1 := paraIntervalo(a.HoraInicio, a.HoraFim)
	iniB, fimB, e2 := paraIntervalo(b.HoraInicio, b.HoraFim)
	if e1 != nil || e2 != nil {
		return false
	}
	return iniA < fimB && fimA > iniB
}

func paraIntervalo(inicio, fim string) (int, int, error) {
	i, err := paraMinutos(inicio)
	if err != nil {
		return 0, 0, err
	}
	f, err := paraMinutos(fim)
	if err != nil {
		return 0, 0, err
	}
	if i >= f {
		return 0, 0, ErrHorarioInvalido
	}
	return i, f, nil
}

func paraMinutos(hhmm string) (int, error) {
	t, err := time.Parse("15:04", hhmm)
	if err != nil {
		return 0, ErrHorarioInvalido
	}
	return t.Hour()*60 + t.Minute(), nil
}
