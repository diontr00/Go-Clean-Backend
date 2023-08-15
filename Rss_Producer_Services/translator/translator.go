package translator

import (
	"embed"
	_ "embed"
	"io/fs"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

type Translator interface {
	translateValidationError(
		ctx *fiber.Ctx,
		fields validator.FieldError,
		plurals interface{},
	) string
	TranslateMessage(c *fiber.Ctx, key string, param ParamConfig, plurals interface{}) string
	ValidateRequest(
		ctx *fiber.Ctx,
		validator *validator.Validate,
		validateStruct interface{},
	) []string
}

type ParamConfig map[string]interface{}

// UniversalTrans concrete
type UniversalTrans struct {
	fs.FileInfo
	bundle *i18n.Bundle
}

func NewUtTrans(fs embed.FS, transFolderName string) (*UniversalTrans, error) {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	dirs, err := fs.ReadDir(transFolderName)
	if err != nil {

		return nil, err
	}

	for _, file := range dirs {
		file_path := filepath.Join(transFolderName, file.Name())

		_, err = bundle.LoadMessageFileFS(fs, file_path)
		if err != nil {
			return nil, err

		}
	}

	return &UniversalTrans{bundle: bundle}, err
}

// Implementing Translator

// Use as the middleware to extract and stock locale in url query for translation services

// Internal usage to transalate validation error to more user friendly mmessage
func (u *UniversalTrans) translateValidationError(
	c *fiber.Ctx,
	fe validator.FieldError,
	plurals interface{},
) string {
	locale := c.Locals("locale")

	localizer := i18n.NewLocalizer(u.bundle, locale.(string))

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: fe.Tag(),
		TemplateData: map[string]interface{}{
			"Field": fe.Field(),
			"Param": fe.Param(),
		},
		PluralCount: plurals,
	})
	if err != nil {
		return "Error Translate Message"
	}

	return message
}

// Translate message with associated key , plurals define the plurals variable that determine more than one form
func (u *UniversalTrans) TranslateMessage(c *fiber.Ctx,
	key string,
	para ParamConfig,
	plurals interface{},
) string {
	locale := c.Locals("locale")
	localizer := i18n.NewLocalizer(u.bundle, locale.(string))
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: para,
		PluralCount:  plurals,
	})
	if err != nil {
		return "Error TranslateMessage"
	}
	return message
}

// Translate Validate Error when parsing request body to a more user friendly form or other lang depend on locale store on url query
func (u *UniversalTrans) ValidateRequest(
	c *fiber.Ctx,
	v *validator.Validate,
	validateStruct interface{},
) []string {
	var return_err []string

	validate_errs := v.Struct(validateStruct)

	if validate_errs != nil {
		for _, err := range validate_errs.(validator.ValidationErrors) {
			err_msg := u.translateValidationError(c, err, nil)
			return_err = append(return_err, err_msg)

		}
		return return_err
	}
	return nil
}
