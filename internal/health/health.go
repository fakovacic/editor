package health

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	healthzStatus = http.StatusOK
	mu            sync.RWMutex
)

func StartServer() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		switch GetHealthStatus() {
		case http.StatusOK:
			SetHealthStatus(http.StatusOK)
		case http.StatusServiceUnavailable:
			SetHealthStatus(http.StatusServiceUnavailable)
		}

		return c.SendStatus(GetHealthStatus())
	})

	return app
}

func GetHealthStatus() int {
	mu.RLock()
	defer mu.RUnlock()

	return healthzStatus
}

func SetHealthStatus(status int) {
	mu.Lock()
	healthzStatus = status
	mu.Unlock()
}
