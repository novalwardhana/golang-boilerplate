package repository

type repository struct {
}

type Repository interface {
}

func NewRepository() Repository {
	return &repository{}
}
