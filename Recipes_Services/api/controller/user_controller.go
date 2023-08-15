package controller

import (
	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/translator"
	"reflect"
	"time"
)

type UserController struct {
	Logger      *logStack.Logger
	Validator   *validator.Validate
	Translator  translator.Translator
	UserUseCase model.UserUseCase
	EncryptKey  string
	Signature   []byte
}

// @Summary Login to make private request
// @Tags User
// @Accept json
// @Produce json
// @Param  recipe body model.UserRequest true "User login request"
// @Param locale  query string  false  "supported locale" Enums("vi" , "en")
// Success return code param type , data type , comment
// @Success  200 string model.UserReturnToken "Access token and expiry"
// @Failure  400 {object} model.UserInvalid "Either user not found or invalid password and username"
// @Router /signin [post]
func (u *UserController) Signin(c *fiber.Ctx) error {
	if c.Locals("authenticated") == true {
		return c.JSON(fiber.Map{
			"status": u.Translator.TranslateMessage(c, "alreadylogin", nil, nil),
		})
	}
	user := &model.UserSigninRequest{}
	if err := c.BodyParser(user); err != nil {

		u.Logger.Error(err.Error(), logStack.String("Client", "Parsing Error"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": u.Translator.TranslateMessage(c, "bodyparse", nil, nil),
		})
	}

	if errors := u.Translator.ValidateRequest(c, u.Validator, user); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})

	}

	userToken, err := u.UserUseCase.GetIDTokenByPass(c.Context(), user)

	if err != nil {
		u.Logger.Error(err.Error(), logStack.String("[Error]", "Auth server"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": u.Translator.TranslateMessage(
				c,
				"tryagain",
				nil,
				nil,
			),
		})
	}
	val := reflect.ValueOf(*userToken)
	error_field := val.FieldByName("Error")
	if !error_field.IsZero() {
		return c.Status(fiber.StatusUnauthorized).JSON(userToken)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    userToken.Id_Token,
		Expires:  time.Now().Add(time.Duration(userToken.Expires_in) * time.Second),
		HTTPOnly: true,
		SameSite: "strict",
	})

	return c.JSON(fiber.Map{
		"Status": u.Translator.TranslateMessage(c, "loggedin", nil, nil),
	})
}

// @Summary  Signup a new user
// @Tags User
// @Accept json
// @Produce jsons
// @Param  recipe body model.UserRequest true "User login request"
// @Param locale  query string  false  "supported locale" Enums("vi" , "en")
// Success return code param type , data type , comment
// @Success  200 string model.UserStatus "Access token and expiry"
// @Failure  400 {object} model.UserInvalid "User Exist or bad request"
// @Router /signup [post]
func (u *UserController) Signup(c *fiber.Ctx) error {
	if c.Locals("authenticated") == true {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": u.Translator.TranslateMessage(c, "loggedin", nil, nil),
		})
	}

	user := &model.UserSignUpRequest{}
	if err := c.BodyParser(user); err != nil {

		u.Logger.Error(err.Error(), logStack.String("Client", "Parsing Error"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": u.Translator.TranslateMessage(c, "bodyparse", nil, nil),
		})
	}

	if errors := u.Translator.ValidateRequest(c, u.Validator, user); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})

	}

	user_response, err := u.UserUseCase.CreateUserByPass(c.Context(), user)
	if err != nil {

		u.Logger.Error(err.Error(), logStack.String("[Error]", "Auth server"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": u.Translator.TranslateMessage(
				c,
				"tryagain",
				nil,
				nil,
			),
		})
	}

	val := reflect.ValueOf(*user_response)
	error_field := val.FieldByName("Error")

	if !error_field.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(user_response)
	}

	return c.Status(fiber.StatusOK).JSON(user_response)
}

// @Summary  Sigout and clear cookie token
// @Tags User
// @Accept json
// @Produce json
// @Param  recipe body model.UserRequest true "User login request"
// @Param locale  query string  false  "supported locale" Enums("vi" , "en")
// Success return code param type , data type , comment
// @Success  200 string model.UserStatusOk "Signout of user"
// @Router /signout [post]
func (u *UserController) SignOut(c *fiber.Ctx) error {
	if c.Locals("authenticated") == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": u.Translator.TranslateMessage(c, "notlogin", nil, nil),
		})
	}

	// NOT work  unless expire
	c.ClearCookie()
	c.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Value:   "",
		Expires: time.Now(),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Ok",
	})
}
