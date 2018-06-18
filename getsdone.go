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
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}

type UserInfo struct {
	ID						string						`json:"id"`
	Email					string						`json:"email"`
	Name          string    				`json:"name"`
	Icon          sql.NullString    `json:"icon"`
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}

type Task struct {
	ID						string						`json:"id"`
	OwnerID       string            `json:"ownerId"`
	DelegateID    sql.NullString    `json:"delegateId"`
	OriginID      sql.NullString    `json:"originId"`
	StateID       sql.NullString    `json:"stateId"`
	PriorityID    sql.NullString    `json:"priorityId"`
	Task          string            `json:"task"`
	Visibility    bool              `json:"visibility"`
	Estimate			sql.NullString		`json:"estimate"`
	Actual				sql.NullString		`json:"actual"`
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}

type Hashtag struct {
	ID						string						`json:"id"`
  Tag       		string            `json:"tag"`
	Description   string            `json:"description"`
}
