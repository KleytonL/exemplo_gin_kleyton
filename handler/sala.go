package handler

import (
	"api-gin/model"
	"api-gin/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalaHandler struct {
	Store *store.Store
}

func NewSalaHandler(store *store.Store) *SalaHandler {
	return &SalaHandler{
		Store: store,
	}
}

func (h *SalaHandler) Criar(c *gin.Context) {
	sala := &model.Sala{}

	if err := c.ShouldBindJSON(sala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	if err := h.Store.CriarSala(sala); err != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": "Sala já existe"})
		return
	}

	c.JSON(http.StatusCreated, sala)
}

func (h *SalaHandler) Listar(c *gin.Context) {
	salas := h.Store.ListarSalas()

	c.JSON(http.StatusOK, salas)
}
