package main

import (
	"database/sql"
)


type User struct {
	ID						string						`json:"id"`
	Email					string						`json:"email"`
	Mobile				sql.NullString		`json:"mobile"`
	Password      string						`json:"password"`
	Salt          string            `json:"salt"`
	Registered  	bool              `json:"registered"`
	Icon          sql.NullString    `json:"icon"`
	RankID        string						`json:"rankdId"`
}
