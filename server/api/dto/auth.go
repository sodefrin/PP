package dto

type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
