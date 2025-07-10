package models
import "github.com/golang-jwt/jwt/v5"
type User struct {
	ID       uint    `json:"id"`
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}