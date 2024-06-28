package storage

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type WrapperStorage interface {
	GetWrapperByType(string) ([]models.Wrapper, error)
}
