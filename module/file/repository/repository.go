package repository

import (
	"gorm.io/gorm"
)

type repository struct {
	dbMaster *gorm.DB
}

type Repository interface {
}

func NewRepository(dbMaster *gorm.DB) Repository {
	return &repository{
		dbMaster: dbMaster,
	}
}
