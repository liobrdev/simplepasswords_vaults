package controllers

import (
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

type Handler struct {
	DB   *gorm.DB
	Conf *config.AppConfig
}
