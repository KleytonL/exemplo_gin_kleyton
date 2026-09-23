package store

import (
	"errors"
	"sync"

	"api-gin/model"
)

var (
	ErrSalaJaExiste  = errors.New("já existe uma sala com este id")
	ErrAlunoJaExiste = errors.New("já existe um aluno com esta matrícula")
	ErrTurmaJaExiste = errors.New("já existe uma turma com este id")
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

func (s *Store) CriarTurma(turma *model.Turma) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.turmas[turma.ID]; exists {
		return ErrTurmaJaExiste
	}

	s.turmas[turma.ID] = turma
	return nil
}

func (s *Store) ListarTurmas() []*model.Turma {
	s.mu.Lock()
	defer s.mu.Unlock()
	var turmas []*model.Turma
	for _, turma := range s.turmas {
		turmas = append(turmas, turma)
	}
	return turmas
}