package colors

import (
	"sync"

	"github.com/fakovacic/editor/internal/app/hub"
)

func New() hub.Colors {
	c := &colors{
		list: make(map[string]bool),
	}

	for i := range colorsList {
		c.list[colorsList[i]] = false
	}

	return c
}

type colors struct {
	list map[string]bool // false = available, true = taken
	lock sync.Mutex
}

func (c *colors) Lock() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k, v := range c.list {
		if !v {
			c.list[k] = true

			return k
		}
	}

	return ""
}

func (c *colors) Release(color string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k, v := range c.list {
		if k == color && v {
			c.list[k] = false

			return
		}
	}
}

var colorsList = []string{
	"blue",
	"red",
	"green",
	"orange",
	"purple",
	"brown",
	"gray",
	"black",
	"chocolate",
	"crimson",
	"violet",
	"darkgreen",
	"darkblue",
	"darkblue",
	"darkcyan",
}
