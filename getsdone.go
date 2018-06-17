package main

import (
	"database/sql"
)

const (
	APP_NAME				= "getsdone"
	GETSDONE        = APP_NAME
	HASH_LENGTH     = 32
	HMAC_KEY        = "spain this summer"
	LOCAL_HOST      = "127.0.0.1"
	PEPPER          = "getsdone is the bomb"
	SALT_LENGTH     = 24
	TOKEN_LENGTH    = 32
	VERSION					= "0.1"
)

type User struct {
	ID						string						`json:"id"`
	Email					string						`json:"email"`
	Name          string    				`json:"name"`
	Mobile				sql.NullString		`json:"mobile"`
	Password      string						`json:"password"`
	Salt          string            `json:"salt"`
	Registered  	bool              `json:"registered"`
	Icon          sql.NullString    `json:"icon"`
	Token         sql.NullString    `json:"token"`
	RankID        string						`json:"rankdId"`
}
