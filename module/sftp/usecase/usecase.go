package usecase

import (
	"github.com/novalwardhana/golang-boilerplate/module/sftp/repository"
)

type usecase struct {
	repo repository.Repository
}

type Usecase interface {
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}
