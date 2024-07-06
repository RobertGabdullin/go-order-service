package integration

import (
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

func Normalize(ords ...*models.Order) {
	for i := range ords {
		ords[i].Expire = ords[i].Expire.UTC().Truncate(time.Second)
		ords[i].DeliveredAt = ords[i].DeliveredAt.UTC().Truncate(time.Second)
		ords[i].ReturnedAt = ords[i].ReturnedAt.UTC().Truncate(time.Second)
	}
}
