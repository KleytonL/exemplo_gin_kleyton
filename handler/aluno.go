package handler

import (
	"api-gin/model"
	"api-gin/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlunoHandler struct {
	Store *store.Store
}

func NewAlunoHandler(store *store.Store) *AlunoHandler {
	return &AlunoHandler{
		Store: store,
	}
}

func (h *AlunoHandler) Criar(c *gin.Context) {
	aluno := &model.Aluno{}

	if err := c.ShouldBindJSON(aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	if err := h.Store.CriarAluno(aluno); err != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": "Aluno já existe"})
		return
	}

	c.JSON(http.StatusCreated, aluno)
}

func (h *AlunoHandler) Listar(c *gin.Context) {
	alunos := h.Store.ListarAlunos()

	c.JSON(http.StatusOK, alunos)
}
