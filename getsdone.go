package main

import (
	"database/sql"
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
