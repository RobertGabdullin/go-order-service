package cache

import (
	"sync"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type InMemoryCache struct {
	mu     sync.RWMutex
	orders map[int]models.Order
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		orders: make(map[int]models.Order),
	}
}

func (c *InMemoryCache) GetOrder(id int) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, found := c.orders[id]
	return order, found
}

func (c *InMemoryCache) SetOrder(id int, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[id] = order
}

func (c *InMemoryCache) InvalidateOrder(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.orders, id)
}
