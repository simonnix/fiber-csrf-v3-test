package csrfv3_test

import (
	"net/url"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	csrfv3 "github.com/simonnix/fiber-csrf-v3-test"
)

var _ = Describe("Csrf", func() {
	It("can login", func() {
		client := TestClient{}.New(csrfv3.GetApp())

		resp := client.GET("/login")
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

		form := url.Values{}
		form.Set("username", "user")
		form.Set("_csrf", client.GetCookieValue(csrfCookieName))

		resp = client.POST("/login", form)
		Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))

		form.Set("password", "user")
		resp = client.POST("/login", form)
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
	})
})
