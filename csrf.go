package csrfv3

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func GetApp() *fiber.App {
	app := fiber.New(
		fiber.Config{
			PassLocalsToViews: true,
			ErrorHandler: func(c fiber.Ctx, err error) error {
				slog.Error("FiberErrorHandler", "error", err.Error())
				code := fiber.StatusInternalServerError
				var e *fiber.Error
				if errors.As(err, &e) {
					code = e.Code
				}
				return c.SendStatus(code)
			},
		},
	)

	loggerHandler := logger.New(logger.Config{
		Next: func(c fiber.Ctx) bool {
			return true
		},
		Format:        logger.DefaultFormat,
		TimeFormat:    "15:04:05",
		TimeZone:      "Local",
		TimeInterval:  500 * time.Millisecond,
		DisableColors: false,
	})

	app.Use(loggerHandler)

	// --- Session middleware (matches user config)
	sessHandler, store := session.NewWithStore(session.Config{
		AbsoluteTimeout:   0 * time.Second,
		CookieDomain:      "",
		CookieHTTPOnly:    false,
		CookiePath:        "",
		CookieSameSite:    "Lax",
		CookieSecure:      false, // user has this false
		CookieSessionOnly: false,
		IdleTimeout:       30 * time.Minute,
	})
	app.Use(sessHandler)

	// --- CSRF middleware (matches user config)
	csrfHandler := csrf.New(csrf.Config{
		CookieDomain:      "",
		CookieHTTPOnly:    false,
		CookieName:        "csrf_",
		CookiePath:        "",
		CookieSameSite:    "Lax",
		CookieSecure:      false,
		CookieSessionOnly: false,
		IdleTimeout:       30 * time.Minute,
		SingleUseToken:    false,
		TrustedOrigins:    []string{},
		Session:           store,
		Extractor: extractors.FromCustom(
			"CsrfExtractor", func(c fiber.Ctx) (string, error) {
				formKeyLookup := "_csrf"
				headerKeyLookup := "X-Csrf-Token"
				if token := c.FormValue(formKeyLookup); token != "" {
					slog.Debug("CsrfExtractor() Form", "token", token, "cookie", c.Cookies("csrf_"))
					return token, nil
				} else if token := c.Get(headerKeyLookup); token != "" {
					slog.Debug("CsrfExtractor() Header", "token", token, "cookie", c.Cookies("csrf_"))
					return token, nil
				} else {
					return "", errors.New("CsrfExtractor() missing csrf token from form and header")
				}
			},
		),
		ErrorHandler: func(_ fiber.Ctx, err error) error {
			slog.Error("CsrfErrorHandler()", "error", err)
			return fiber.ErrForbidden
		},
	})
	app.Use(csrfHandler)

	// --- GET handler to get CSRF cookie
	app.Get("/login", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		if sess == nil {
			return fiber.ErrInternalServerError
		}
		if sess.Get("key") != nil {
			return c.SendString("already set")
		} else {
			sess.Set("key", "value")
		}
		return c.SendString("login page")
	})

	// --- POST handler
	app.Post("/login", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		if sess == nil {
			return fiber.ErrInternalServerError
		}
		val, ok := sess.Get("key").(string)
		if !ok || val != "value" {
			return fiber.ErrNotFound
		}
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username != "user" || password != "user" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		sess.Set("user", username)
		return c.SendStatus(fiber.StatusOK)
	})

	return app
}
