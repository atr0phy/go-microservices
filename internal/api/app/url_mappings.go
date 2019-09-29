package app

import (
	"github.com/atr0phy/go-microservices/internal/api/controllers/polo"
	"github.com/atr0phy/go-microservices/internal/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", polo.Marco)
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
