package csrfv3_test

import (
	"log/slog"
	"net/url"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	csrfv3 "github.com/simonnix/fiber-csrf-v3-test"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	options := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, &options)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	client := TestClient{}.New(csrfv3.GetApp())

	resp := client.GET("/login")
	assert.Equal(resp.StatusCode, fiber.StatusOK)

	form := url.Values{}
	form.Set("username", "user")
	form.Set("_csrf", client.GetCookieValue(csrfCookieName))

	resp = client.POST("/login", form)
	assert.Equal(resp.StatusCode, fiber.StatusUnauthorized)

	form.Set("password", "user")
	resp = client.POST("/login", form)
	assert.Equal(resp.StatusCode, fiber.StatusOK)
}
