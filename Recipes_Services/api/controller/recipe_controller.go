package controller

import (
	"errors"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/translator"
	"khanhanhtr/sample/usecase"
	"strconv"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeController struct {
	Logger        *logStack.Logger
	Validator     *validator.Validate
	Translator    translator.Translator
	RecipeUseCase model.RecipeUseCase
}

// @Summary Create a new recipe
// @Tags Recipes
// @Accept json
// @Produce json
// @Param  recipe body model.RecipeCreateRequest true "Recipe to create"
// @Param locale  query string  false  "supported locale" Enums("vi" , "en")
// Success return code param type , data type , comment
// @Success  200 string model.RecipeCreateResponse "Id of the created recipe"
// @Failure  404 {object} model.RecipeNotFoundError "Recipe Cannot be found"
// @Failure  400 {object} model.RecipeParsingError "Parsing problem , either wrong type or validation constraint will be explain"
// @Failure  500 {object} model.InternalServerError "Problem with internal server"
// @securityDefinitions.apikey ApiKeyAuth
// @Router /recipes [post]
func (rc *RecipeController) Create(c *fiber.Ctx) error {
	if c.Locals("authenticated") == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": rc.Translator.TranslateMessage(c, "missingtoken", nil, nil),
		})
	}
	recipe := &model.Recipe{}

	if err := c.BodyParser(recipe); err != nil {
		rc.Logger.Error(err.Error(), logStack.String("Client", "Parsing Error"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": rc.Translator.TranslateMessage(c, "bodyparse", nil, nil),
		})
	}

	if errors := rc.Translator.ValidateRequest(c, rc.Validator, recipe); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})

	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	// TODO  : Proper error handling
	id, err := rc.RecipeUseCase.Create(c.Context(), recipe)

	if err != nil {
		rc.Logger.Error(err.Error(), logStack.String("Server", "Create Error"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot Create error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Success": id,
	})

}

// @Summary Get all the recipe
// @Tags Recipes
// @Accept json
// @Produce json
// @Param locale  query string  false "supported locale" Enums("vi" , "en") example("vi")
// @Success  200 {array} model.RecipeGetResponse "All of the recipe returned"
// @Failure  404 {object} model.RecipeNotFoundError "Recipe repo is currently empty"
// @Failure  500 {object} model.InternalServerError "Problem with internal server"
// @Router /recipes [get]
func (rc *RecipeController) Get(c *fiber.Ctx) error {

	page := c.Query("page", "0")
	limit := c.Query("limit", "0")

	page_num, err := strconv.Atoi(page)
	if err != nil {
		page_num = 0
	}
	limit_num, err := strconv.Atoi(limit)
	if err != nil {
		limit_num = 0
	}

	recipes, errs := rc.RecipeUseCase.Get(
		c.Context(),
		&model.FindOptions{Page: int64(page_num), Limit: int64(limit_num)},
	)
	if errs != nil {
		if errors.Is(errs, usecase.ErrInternal) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "try again",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errs,
		})
	}

	return c.Status(fiber.StatusOK).JSON(recipes)

}

// @Summary Get particular recipe
// @Tags Recipes
// @Accept json
// @Produce json
// @Param  recipe_id path string true  "The ID associate with the recipe"
// @Param locale  query string  false "supported locale" Enums("vi" , "en") example("vi")
// @Success  200 {object} model.RecipeGetResponse "All of the recipe returned"
// @Failure  404 {object} model.RecipeNotFoundError "Recipe id not found"
// @Failure  500 {object} model.InternalServerError "Problem with internal server"
// @Router /recipes/{:id} [get]
func (rc *RecipeController) GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "empty",
		})
	}
	recipe, err := rc.RecipeUseCase.GetById(c.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrInternal) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "try again",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})

	}

	return c.JSON(recipe)
}

// @Summary Update particular recipe
// @Tags Recipes
// @Accept json
// @Produce json
// @Param  id path string true  "The ID associate with the recipe"
// @Param  recipe body model.RecipeCreateRequest true "Recipe to update"
// @Param locale  query string  false "supported locale" Enums("vi" , "en") example("vi")
// @Success  200 string  "Id of the updated recipe"
// @Failure  404 {object} model.RecipeNotFoundError "Recipe id not found"
// @Failure  500 {object} model.InternalServerError "Problem with internal server"
// @securityDefinitions.apikey ApiKeyAuth
// @Router /private/recipes/{:id} [put]
func (rc *RecipeController) Update(c *fiber.Ctx) error {

	if c.Locals("authenticated") == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": rc.Translator.TranslateMessage(c, "missingtoken", nil, nil),
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusOK).JSON("empty")
	}

	var recipe model.Recipe

	err := c.BodyParser(&recipe)
	if err != nil {
		rc.Logger.Error(err.Error(), logStack.String("Client", "Parsing Error"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unknown message format",
		})
	}

	if errors := rc.Translator.ValidateRequest(c, rc.Validator, recipe); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	// TODO: Proper error handling
	recipe_count, err := rc.RecipeUseCase.UpdateOne(c.Context(), id, &recipe)
	if err != nil {

		rc.Logger.Error(
			"[UPDATE RECIPES]",
			logStack.Any("error", err.(*multierror.Error).WrappedErrors()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server error",
		})
	}
	return c.JSON(recipe_count)
}

// @Summary Delete particular recipe
// @Tags Recipes
// @Accept json
// @Produce json
// @Param  id path string true "The ID associate with the recipe"
// @Param locale  query string  false "supported locale" Enums("vi" , "en")
// @Success  200 {object} string "Id of the deleted recipe"
// @Failure  404 {object} model.RecipeNotFoundError "Recipe id not found"
// @Failure  500 {object} model.InternalServerError "Problem with internal server"
// @securityDefinitions.apikey ApiKeyAuth
// @Router /recipes/{:id} [delete]
func (rc *RecipeController) Delete(c *fiber.Ctx) error {
	if c.Locals("authenticated") == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": rc.Translator.TranslateMessage(c, "missingtoken", nil, nil),
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusOK).JSON("empty")
	}

	// TODO: Proper Error Handling
	recipe_count, err := rc.RecipeUseCase.Delete(c.Context(), id)
	if err != nil {

		rc.Logger.Error(
			"[DELETE RECIPES]",
			logStack.Any("error", err.(*multierror.Error).WrappedErrors()),
		)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server error",
		})
	}
	return c.JSON(recipe_count)
}
