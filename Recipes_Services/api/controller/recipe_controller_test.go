package controller_test

import (
	"fmt"
	"io"
	"testing"

	"khanhanhtr/sample/redis"

	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
)

func TestCreate_Unit(t *testing.T) {
	Convey("When : User create a recipe", t, func() {
		Convey("Context: User is authenticated", func() {

			// TODO: why only work when new fiber instance get create ?
			app := fiber.New()
			defer app.Shutdown()
			app.Post("/create-authenticated", func(c *fiber.Ctx) error {
				c.Locals("locale", "en")
				c.Locals("authenticated", true)

				return c.Next()
			}, g.recipecontroller.Create)
			Convey("Behavior : User submit correct recipe struct", func() {

				g.repo_mock.On("Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Recipe")).
					Return(g.sample_recipe.ID.Hex(), nil).
					Once()

				g.redis_mock.On("Rpush", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(&redis.RedisIntegerStatus{}).
					Once()

				req := newHttpRequest(g.sample_recipe, "POST", "/create-authenticated")

				res, err := app.Test(req)

				Convey("Then : User should return success message with associated id", func() {

					So(err, ShouldBeNil)

					So(res.StatusCode, ShouldEqual, 200)

					body, err := io.ReadAll(res.Body)
					So(err, ShouldBeNil)
					So(
						gjson.GetBytes(body, "Success").String(),
						ShouldEqual,
						g.sample_recipe.ID.Hex(),
					)
					g.repo_mock.AssertExpectations(t)
					g.redis_mock.AssertExpectations(t)
				})
			})

		})
		Convey("Context: User is not authenticated", func() {

			// TODO: why only work when new fiber instance get create ?
			app := fiber.New()
			defer app.Shutdown()

			app.Post("/create-unauthenticated", func(c *fiber.Ctx) error {
				fmt.Println("here")
				c.Locals("locale", "en")
				c.Locals("authenticated", false)
				return c.Next()
			}, g.recipecontroller.Create)

			Convey("User should be enforced login with unauthorized message", func() {

				g.repo_mock.On("Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Recipe")).
					Return(g.sample_recipe.ID.Hex(), nil).
					Times(0)

				g.redis_mock.On("Rpush", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string"), mock.Anything).
					Return(&redis.RedisIntegerStatus{}).
					Times(0)

				req := newHttpRequest(g.sample_recipe, "POST", "/create-unauthenticated")

				res, err := app.Test(req)

				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, fiber.StatusUnauthorized)

				body, err := io.ReadAll(res.Body)
				So(err, ShouldBeNil)
				So(
					gjson.GetBytes(body, "error").String(),
					ShouldNotBeEmpty,
				)
				g.repo_mock.AssertExpectations(t)
				g.redis_mock.AssertExpectations(t)

			})

		})

	})

}
