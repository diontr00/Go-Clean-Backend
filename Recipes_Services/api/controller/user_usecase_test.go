package controller_test

import (
	"io"
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
)

func TestSignin_Unit(t *testing.T) {
	Convey("When User is already authenticated", t, func() {

		app := fiber.New()
		defer app.Shutdown()

		app.Post("/signin", func(c *fiber.Ctx) error {
			c.Locals("locale", "en")
			c.Locals("authenticated", true)

			return c.Next()
		}, g.userController.Signin)

		Convey("Behavior: it should return Json already login", func() {
			req := newHttpRequest(g.sample_user_signin, "POST", "/signin")

			res, err := app.Test(req)

			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, 200)

			body, err := io.ReadAll(res.Body)
			So(err, ShouldBeNil)
			So(gjson.GetBytes(body, "status").String(), ShouldContainSubstring, "already")

			g.useruc_mock.AssertNotCalled(t,
				"GetIdTokenByPass",
				mock.AnythingOfType("context.Context"),
				&g.sample_user_signin,
			)

		})
	})

}
