package auth

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/krespix/gamification-api/internal/clients/smtp"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/cache/auth"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
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

type Code struct {
	Code int
}

type Mail struct {
	Sender  string
	To      string
	Subject string
	Body    string
}

type service struct {
	smtpClient smtp.Client

	userRepo user.Repository
	authRepo auth.Repository

	validate *validator.Validate

	jwtSecret         string
	defaultExpiration time.Duration
}

func buildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", mail.To)
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
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

	err = s.authRepo.CreateCode(ctx, email, code)
	if err != nil {
		return err
	}

	t := template.Must(template.New("template_data").Parse(emailTemplate))
	var body bytes.Buffer
	cd := Code{Code: code}
	err = t.Execute(&body, cd)
	if err != nil {
		return err
	}

	mail := Mail{
		Sender:  "gamification-noreply@mail.ru",
		To:      email,
		Subject: "Код авторизации",
		Body:    body.String(),
	}
	res := buildMessage(mail)

	return s.smtpClient.Send(email, []byte(res))
}

func (s *service) VerifyCode(ctx context.Context, email string, code int) (string, error) {
	savedCode, err := s.authRepo.GetCode(ctx, email)
	if err != nil {
		return "", err
	}
	if code != savedCode {
		return "", fmt.Errorf("received code does not match with saved")
	}
	err = s.authRepo.DeleteCode(ctx, email)
	if err != nil {
		return "", err
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
