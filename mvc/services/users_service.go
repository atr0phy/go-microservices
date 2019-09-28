package services

import (
	"github.com/atr0phy/go-microservices/mvc/domain"
	"github.com/atr0phy/go-microservices/mvc/utils"
)

func GetUser(userId int64) (*domain.User, *utils.ApplicationError) {
	return domain.GetUser(userId)
}
