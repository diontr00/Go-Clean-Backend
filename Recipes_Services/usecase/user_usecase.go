package usecase

import (
	"context"
	"khanhanhtr/sample/model"
	"time"
)

type userUseCase struct {
	repo           model.UserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(
	timeout time.Duration,
	repo model.UserRepository,
) *userUseCase {
	return &userUseCase{
		contextTimeout: timeout,
		repo:           repo,
	}
}

// Get ID Token by oauth password grant type
func (l *userUseCase) GetIDTokenByPass(
	c context.Context,
	user *model.UserSigninRequest,
) (*model.UserSigninResponse, error) {
	return l.repo.GetIDTokenByPass(c, user)
}

func (l *userUseCase) CreateUserByPass(
	c context.Context,
	user *model.UserSignUpRequest,
) (*model.UserSignUpResponse, error) {

	return l.repo.CreateUserByPass(c, user)
}
