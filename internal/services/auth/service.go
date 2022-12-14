package auth

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/krespix/gamification-api/internal/clients/smtp"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/cache/auth"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
	"math/rand"
	"time"
)

// const for generate auth codes
const (
	low  int = 1000
	high int = 9999
)

type Service interface {
	SendCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email string, code int) (string, error)
	ValidateToken(ctx context.Context, token string) (*models.Claims, error)
}

type service struct {
	smtpClient smtp.Client

	userRepo user.Repository
	authRepo auth.Repository

	validate *validator.Validate

	jwtSecret         string
	defaultExpiration time.Duration
}

func (s *service) SendCode(ctx context.Context, email string) error {
	err := s.validate.Var(email, "email")
	if err != nil {
		return err
	}
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user with email %s does not exists", email)
	}
	usr, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if usr.DeletedAt.Valid {
		return fmt.Errorf("user deleted")
	}

	code := generateCode()

	text := fmt.Sprintf(
		"From: gamification-noreply@mail.ru\r\nTo: %s\r\nSubject: Authorization code\r\n\r\nCode: %d \r\n",
		email, code)
	msg := []byte(text)

	err = s.authRepo.CreateCode(ctx, email, code)
	if err != nil {
		return err
	}

	return s.smtpClient.Send(email, msg)
}

func (s *service) VerifyCode(ctx context.Context, email string, code int) (string, error) {
	savedCode, err := s.authRepo.GetCode(ctx, email)
	if err != nil {
		return "", err
	}
	err = s.authRepo.DeleteCode(ctx, email)
	if err != nil {
		return "", err
	}
	if code != savedCode {
		return "", fmt.Errorf("received code does not match with saved")
	}
	usr, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	token, err := generateToken(usr.ID, usr.Role, s.jwtSecret, s.defaultExpiration)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateCode() int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return low + r.Intn(high-low)
}

func generateToken(id int64, role models.Role, jwtSecret string, defaultExpiration time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		ID:        id,
		Role:      role,
		ExpiresAt: time.Now().Add(defaultExpiration),
	})
	token, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) (*models.Claims, error) {
	claims := &models.Claims{}
	res, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	cl, ok := res.Claims.(*models.Claims)
	if !ok || !res.Valid {
		return nil, fmt.Errorf("unauthorized")
	}
	//проверить что юзер не удален
	repoUser, err := s.userRepo.Get(ctx, cl.ID)
	if err != nil {
		return nil, err
	}
	if repoUser.DeletedAt.Valid {
		return nil, fmt.Errorf("user deleted")
	}
	return cl, nil
}

func New(smtpClient smtp.Client,
	userRepo user.Repository,
	authRepo auth.Repository,
	validate *validator.Validate,
	jwtSecret string,
	defaultExpiration time.Duration,
) Service {
	return &service{
		smtpClient:        smtpClient,
		userRepo:          userRepo,
		authRepo:          authRepo,
		validate:          validate,
		jwtSecret:         jwtSecret,
		defaultExpiration: defaultExpiration,
	}
}
