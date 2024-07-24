package cache

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type Cache interface {
	GetOrder(id int) (models.Order, bool)
	SetOrder(id int, order models.Order)
	InvalidateOrder(id int)
}
