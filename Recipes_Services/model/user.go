package model

import "context"

type UserSigninRequest struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type UserSignUpRequest struct {
	Email     string `json:"email"                validate:"required,email"`
	Password  string `json:"password"             validate:"required,containsany=@!*&#$,min=8"`
	Username  string `json:"username"             validate:"required,alpha"`
	Givenname string `json:"given_name,omitempty"`
}

type UserSignUpResponse struct {
	Username          string `json:"username,omitempty"`
	Email             string `json:"email,omitempty"`
	Error             string `json:"code,omitempty"`
	Error_Description string `json:"description,omitempty"`
}

type UserSigninResponse struct {
	Access_Token      string `json:"access_token,omitempty"`
	Expires_in        int    `json:"expires_in,omitempty"`
	Id_Token          string `json:"id_token,omitempty"`
	Error             string `json:"error,omitempty"`
	Error_Description string `json:"error_description,omitempty"`
}

type Auth0SigninRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	GrantType    string `json:"grant_type"`
}

type Auth0SignupRequest struct {
	ClientId   string `json:"client_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Connection string `json:"connection"`
	Username   string `json:"username,omitempty"`
}

//go:generate mockery --name UserRepository
type UserRepository interface {
	// Get User ID Token by password type
	GetIDTokenByPass(ctx context.Context, user *UserSigninRequest) (*UserSigninResponse, error)
	// Create User by password type
	CreateUserByPass(ctx context.Context, user *UserSignUpRequest) (*UserSignUpResponse, error)
}

//go:generate mockery --name UserUseCase
type UserUseCase interface {
	// Get User ID Token by password type
	GetIDTokenByPass(ctx context.Context, user *UserSigninRequest) (*UserSigninResponse, error)

	// Create User by password type
	CreateUserByPass(ctx context.Context, user *UserSignUpRequest) (*UserSignUpResponse, error)
}
