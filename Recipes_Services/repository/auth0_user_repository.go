package repository

import (
	"context"
	"fmt"
	"khanhanhtr/sample/config"
	"khanhanhtr/sample/model"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type auth0UserRepo struct {
	mu          sync.RWMutex
	env         *config.Auth0Env
	apiEndpoint map[string]string
}

func (a *auth0UserRepo) GetIDTokenByPass(
	ctx context.Context,
	user *model.UserSigninRequest,
) (*model.UserSigninResponse, error) {

	request_body := model.Auth0SigninRequest{
		ClientId:     a.env.Client_ID,
		ClientSecret: a.env.Client_Secret,
		Username:     user.Username,
		Password:     user.Password,
		GrantType:    "password",
	}

	a.mu.RLock()
	defer a.mu.Unlock()
	agent := fiber.Post(a.apiEndpoint["idtoken"]).JSON(request_body)
	defer fiber.ReleaseAgent(agent)

	if err := agent.Parse(); err != nil {
		return nil, err
	}
	var response model.UserSigninResponse

	_, _, _ = agent.Struct(&response)
	return &response, nil
}

func (a *auth0UserRepo) CreateUserByPass(
	ctx context.Context,
	user *model.UserSignUpRequest,
) (*model.UserSignUpResponse, error) {

	request_body := model.Auth0SignupRequest{
		ClientId:   a.env.Client_ID,
		Connection: a.env.AUTH_DB_NAME,
		Email:      user.Email,
		Username:   user.Username,
		Password:   user.Password,
	}

	a.mu.RLock()
	defer a.mu.Unlock()
	agent := fiber.Post(a.apiEndpoint[""]).JSON(request_body)

	if err := agent.Parse(); err != nil {
		return nil, err
	}

	var response model.UserSignUpResponse
	_, _, _ = agent.Struct(&response)

	return &response, nil
}

func NewAuth0UserRepo(env *config.Auth0Env) model.UserRepository {
	// Signup endpoint
	signup := fmt.Sprintf("https://%s/dbconnections/signup", env.Domain)
	// Signin endpoint
	idtoken := fmt.Sprintf("https://%s/oauth/token", env.Domain)

	return &auth0UserRepo{
		env: env,
		apiEndpoint: map[string]string{
			"signup":  signup,
			"idtoken": idtoken,
		},
	}

}
