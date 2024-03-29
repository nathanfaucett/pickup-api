package model

type TokenST struct {
	AccessToken string `json:"access_token" validate:"required"`
	TokenType   string `json:"token_type" validate:"required"`
} // @name Token
