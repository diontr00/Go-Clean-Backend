package controller

import (
	"khanhanhtr/sample2/model"
	"khanhanhtr/sample2/translator"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type RssController struct {
	Logger     *logStack.Logger
	Validator  *validator.Validate
	Translator translator.Translator
	RssUsecase model.RssUsecase
}

func (r *RssController) ParseRss(c *fiber.Ctx) error {
	var rssRequest model.RssRequest
	if err := c.BodyParser(&rssRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": r.Translator.TranslateMessage(c, "bodyparse", nil, nil),
		})
	}

	if errors := r.Translator.ValidateRequest(c, r.Validator, rssRequest); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	err := r.RssUsecase.FetchAndInsert(c.Context(), &rssRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": r.Translator.TranslateMessage(
			c,
			"fetchreceive",
			translator.ParamConfig{"URL": rssRequest.URL},
			nil,
		),
	})

}
