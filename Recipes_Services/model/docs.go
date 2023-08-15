//go:build exclude

// This file is exclude when build and used only for swagger doc
package model

type RecipeCreateRequest struct {
	Name         string   `json:"name"         validate:"required" example:"Vietnamese Pho"`
	Tags         []string `json:"tags"         validate:"optional" example:"vietnamese,asian"                            extensions:"x-nullable"`
	Ingredients  []string `json:"ingredients"  validate:"required" example:"Beef,Noodle"`
	Instructions []string `json:"instructions" validate:"required" example:"Slow cook beef for 2 hours , prepare noodle"`
}

type RecipeCreateResponse struct {
	ID string `json:"id" example:"1234567abc"`
}
type RecipeParsingError struct {
	Error string `json:"error" example:"Payload parsing error reason"`
}

type RecipeNotFoundError struct {
	Error string `json:"error" example:"not found"`
}
type InternalServerError struct {
	Error string `json:"error" default:"internal server error"`
}

type RecipeGetResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

type UserReturnToken struct {
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}

type UserInvalid struct {
	error string `json:"error"`
}

type UserStatus struct {
	status string `json:"status"`
}

type UserRequest struct {
	username string `json:"username"`
	password string `json:"password"`
}
