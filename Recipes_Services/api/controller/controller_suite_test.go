package controller_test

import (
	"fmt"
	"io"
	"khanhanhtr/sample/api/controller"
	"khanhanhtr/sample/config"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/model/mocks"
	redismock "khanhanhtr/sample/redis/mocks"
	"khanhanhtr/sample/translator"
	"net/http/httptest"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var app *fiber.App
var usecase_mock *mocks.RecipeUseCase
var redis_mock *redismock.Client
var sample_recipe model.Recipe
var trans_mock translator.Translator

var _ = BeforeSuite(func() {
	app = fiber.New()
	logger := logStack.DefaultLogger()
	validator := validator.New()
	usecase := mocks.NewRecipeUseCase(GinkgoT())
	trans_mock = config.NewTrans()

	redis := redismock.NewClient(GinkgoT())
	redis_mock = redis

	sample_recipe = model.Recipe{
		ID:          primitive.NewObjectID(),
		Name:        "Com tam",
		Tags:        []string{"vietnamese food"},
		Ingredients: []string{"something1", "something2"},
		PublishedAt: time.Now(),
	}

	controller := &controller.RecipeController{
		Logger:        logger,
		Validator:     validator,
		Translator:    trans_mock,
		RecipeUseCase: usecase,
	}

	usecase_mock = usecase

	app.Get("/getbyid/:id?", controller.GetById)
	app.Post("/create", controller.Create)
	app.Get("/get", controller.Get)

	app.Post("/update", controller.Update)
	app.Post("/delete", controller.Delete)

})

var _ = Describe("Get Recipe By Id", func() {

	When("The recipe id is available", func() {
		Context("User did not submit the id via path", func() {

			It("Should return the status empty response", func() {

				req := httptest.NewRequest(
					fiber.MethodGet,
					"/getbyid",
					nil,
				)

				res, err := app.Test(req)

				usecase_mock.AssertNotCalled(
					GinkgoT(),
					"GetById",
					mock.Anything,
					mock.AnythingOfType("string"),
				)
				Expect(err).To(Succeed())
				Expect(res.Body).NotTo(BeNil())
				Expect(res.StatusCode).To(Equal(200))
				resp, err := io.ReadAll(res.Body)
				Expect(err).To(Succeed())
				got_status := gjson.GetBytes(resp, "status")
				Expect(got_status.String()).To(Equal("empty"))
			})

		})

		Context("And user submit the id via path", func() {

			It("should return the correct sample recipe", func() {
				req := httptest.NewRequest(
					fiber.MethodGet,
					fmt.Sprintf("/getbyid/%s", sample_recipe.ID.Hex()),
					nil,
				)
				usecase_mock.On("GetById", mock.Anything, mock.AnythingOfType("string")).
					Return(sample_recipe, nil).
					Once()

				res, err := app.Test(req)
				Expect(err).To(Succeed())
				Expect(res.Body).NotTo(BeNil())
				Expect(res.StatusCode).To(Equal(200))
				usecase_mock.AssertExpectations(GinkgoT())
				resp, err := io.ReadAll(res.Body)
				Expect(err).To(Succeed())
				got_id := gjson.GetBytes(resp, "id")
				Expect(got_id.String()).To(Equal(sample_recipe.ID.Hex()))
			})
		})
	})

})
