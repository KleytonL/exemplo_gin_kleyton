package main

import (
	"api-gin/handler"
	"api-gin/store"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	st := store.NewStore()
	salaHandler := handler.NewSalaHandler(st)
	alunoHandler := handler.NewAlunoHandler(st)
	turmaHandler := handler.NewTurmaHandler(st)

	r := gin.New()

	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})
		v1.POST("/salas", salaHandler.Criar)
		v1.GET("/salas", salaHandler.Listar)
		v1.POST("/alunos", alunoHandler.Criar)
		v1.GET("/alunos", alunoHandler.Listar)
		v1.GET("/alunos/:id", alunoHandler.Buscar)
		v1.POST("/turmas", turmaHandler.Criar)
		v1.GET("/turmas", turmaHandler.Listar)
		v1.GET("/turmas/:id/alunos", turmaHandler.ListarAlunos)
		v1.POST("/turmas/:id/alunos", turmaHandler.Matricular)
		v1.POST("/turmas/:id/alocar", turmaHandler.Alocar)
	}

	r.Run(":8080")
}
