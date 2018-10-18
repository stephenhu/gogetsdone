package main

import (
	"database/sql"
)

const (
	APP_NAME				= "getsdone"
	BLOCK_KEY				= "abcdefabcdefabcd"
	GETSDONE        = APP_NAME
	HASH_LENGTH     = 32
	HMAC_KEY        = "spain this summer"
	IV							= "this is ricky bo"
	LOCAL_HOST      = "127.0.0.1"
	PEPPER          = "getsdone is the bomb"
	SALT_LENGTH     = 24
	TOKEN_LENGTH    = 32
	VERSION					= "0.1"
)

const (
	ICON_MAX_BYTES	= 1024 * 100
)

const (
	PASSWD_RULE_LENGTH		= 8
)

const (
	TASK_OPEN				= "open"
	TASK_COMPLETED	= "completed"
	TASK_ASSIGNED  	= "assigned"
	TASK_DEFERRED   = "deferred"
	TASK_ALL        = "all"
)

const (
	ACTION_COMPLETED		= "completed"
	ACTION_DEFERRED			= "deferred"
	ACTION_UNDEFERRED   = "undeferred"
)

const (
	CONTACT_ACCEPTED		= "accepted"
	CONTACT_DECLINED    = "declined"
	CONTACT_PENDING     = "pending"
	CONTACT_REQUESTED   = "requested"
)

const (
	UPDATE_TASK_DEFERRED			= 1
	UPDATE_TASK_UNDEFERRED		= 0
)

type VersionInfo struct {
	Version 			string						`json:"version"`
}

type Rank struct {
	ID						string						`json:"id"`
	Name          string            `json:"name"`
	Count         int               `json:"count"`
}

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
	RankName      string            `json:"rankName"`
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}

type UserInfo struct {
	ID						string						`json:"id"`
	Email					string						`json:"email"`
	Name          string    				`json:"name"`
	Icon          sql.NullString    `json:"icon"`
	RankName      string            `json:"rankName"`
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}


type Contact struct {
	ID						string						`json:"id"`
	ContactID			string						`json:"contactId"`
	UserID				string            `json:"userId"`
	ContactName		string						`json:"contactName"`
	ContactIcon   sql.NullString		`json:"contactIcon"`
	State			    string						`json:"state"`
}


type Comment struct {
	ID						string						`json:"id"`
	UserID				string						`json:"userId"`
	UserName      string            `json:"userName"`
	UserIcon      sql.NullString     `json:"userIcon"`
	TaskID				string						`json:"taskId"`
	Comment       string						`json:"comment"`
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}

type Task struct {
	ID						string						`json:"id"`
	OwnerID       string            `json:"ownerId"`
	OwnerName     string            `json:"ownerName"`
	OwnerIcon     sql.NullString    `json:"ownerIcon"`
	DelegateID    sql.NullString    `json:"delegateId"`
	OriginID      sql.NullString    `json:"originId"`
	StateID       sql.NullString    `json:"stateId"`
	PriorityID    sql.NullString    `json:"priorityId"`
	Task          string            `json:"task"`
	Deferred      bool            	`json:"deferred"`
	Visibility    bool              `json:"visibility"`
	Actual				sql.NullString		`json:"actual"`
	Comments      []Comment         `json:"comments"`
	Created  			string						`json:"created"`
	Updated			  string						`json:"updated"`
}

type Hashtag struct {
	ID						string						`json:"id"`
  Tag       		string            `json:"tag"`
}

type CookieData struct {
	ID				string				`json:"id"`
	Token			string				`json:"token"`
	Icon			string				`json:"icon"`
}
