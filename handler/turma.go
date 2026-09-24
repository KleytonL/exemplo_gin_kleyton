package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-gin/model"
	"api-gin/store"
)

type TurmaHandler struct {
	Store *store.Store
}

func NewTurmaHandler(s *store.Store) *TurmaHandler {
	return &TurmaHandler{Store: s}
}

func (h *TurmaHandler) Criar(c *gin.Context) {
	turma := &model.Turma{}
	if err := c.ShouldBindJSON(turma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}
	if err := h.Store.CriarTurma(turma); err != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, turma)
}

func (h *TurmaHandler) Listar(c *gin.Context) {
	c.JSON(http.StatusOK, h.Store.ListarTurmas())
}

func (h *TurmaHandler) ListarAlunos(c *gin.Context) {
	turmaID := c.Param("id")
	alunos, err := h.Store.ListarAlunosDaTurma(turmaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alunos)
}

type matricularRequest struct {
	AlunoID string `json:"aluno_id"`
}

func (h *TurmaHandler) Matricular(c *gin.Context) {
	turmaID := c.Param("id")
	var req matricularRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	err := h.Store.MatricularAluno(turmaID, req.AlunoID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"mensagem": "aluno matriculado com sucesso"})
		return
	}

	switch err {
	case store.ErrTurmaNaoEncontrada, store.ErrAlunoNaoEncontrado:
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
	case store.ErrAlunoJaMatriculado, store.ErrConflitoDeAgenda:
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
	case store.ErrCapacidadeExcedida:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
	}
}

func (h *TurmaHandler) Alocar(c *gin.Context) {
	turmaID := c.Param("id")
	aloc := &model.Alocacao{}
	if err := c.ShouldBindJSON(aloc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	err := h.Store.AlocarSala(turmaID, aloc)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"mensagem": "sala alocada com sucesso"})
		return
	}

	switch err {
	case store.ErrTurmaNaoEncontrada, store.ErrSalaNaoEncontrada:
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
	case store.ErrConflitoDeSala:
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
	case store.ErrCapacidadeExcedida:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
	}
}
