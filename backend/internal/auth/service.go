package auth

import (
	"context"
	"errors"
	"time"

	"stockstalk/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	users  *user.Service
	secret []byte
}

func NewService(users *user.Service, secret string) *Service {
	return &Service{
		users:  users,
		secret: []byte(secret),
	}
}

// Register.
func (s *Service) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*user.User, string, error) {

	u, err := s.users.Register(ctx, name, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := s.token(u.ID)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

// Login.
func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (*user.User, string, error) {

	u, err := s.users.Authenticate(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := s.token(u.ID)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *Service) token(userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	return t.SignedString(s.secret)
}

type userIDContextKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// Current user.
func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDContextKey{}).(string)
	return id, ok
}

func (s *Service) Parse(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", errors.New("invalid subject")
	}

	return sub, nil
}
