package model

import "github.com/google/uuid"

type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refresh_token"`
}

type IDToken struct {
	SS string `json:"id_token"`
}

type TokenPair struct {
	IDToken
	RefreshToken
}
