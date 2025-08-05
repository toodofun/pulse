package server

import (
	"github.com/gin-gonic/gin"
	"pulse/internal/infra"
)

type Controller interface {
	RegisterRoute(group *gin.RouterGroup)
}

type Service interface {
	Initialize(db *infra.Database) error
}
