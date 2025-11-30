package api

import (
	"puyo-server/server/db"
)

type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var Queries *db.Queries
