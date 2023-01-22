package auth

//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/brianvoe/gofakeit/v6"
//	"github.com/go-playground/validator/v10"
//	. "github.com/onsi/ginkgo/v2"
//	. "github.com/onsi/gomega"
//	"gitlab.com/krespix/gamification-api/internal/clients/smtp"
//	"gitlab.com/krespix/gamification-api/internal/models"
//	"gitlab.com/krespix/gamification-api/internal/repositories/cache/auth"
//	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
//	"testing"
//	"time"
//)
//
//func TestAuthService(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "Auth service suite")
//}
//
//var _ = Describe("test auth service", func() {
//	ctx := context.Background()
//	t := &testing.T{}
//
//	fake := gofakeit.New(0)
//	jwtSecret := fake.MonthString()
//	validate := validator.New()
//	userRepo := user.NewMockRepository(t)
//	authCache := auth.NewMockRepository(t)
//	smtpClient := smtp.NewMockClient(t)
//	defExpiration := time.Second * 120
//
//	authService := New(
//		smtpClient,
//		userRepo,
//		authCache,
//		validate,
//		jwtSecret,
//		defExpiration,
//	)
//
//	Context("send code", func() {
//		var email string
//		var repoUser *models.User
//		var code int
//		var text string
//		var msg []byte
//		BeforeEach(func() {
//			email = fake.Email()
//			repoUser = &models.User{
//				ID:        fake.Int64(),
//				Email:     email,
//				CreatedAt: time.Now(),
//				DeletedAt: sql.NullTime{
//					Valid: false,
//				},
//				Role: "user",
//			}
//			code = generateCode()
//			text = fmt.Sprintf(
//				"From: gamification-noreply@mail.ru\r\nTo: %s\r\nSubject: Authorization code\r\n\r\nCode: %d \r\n",
//				email, code)
//			msg = []byte(text)
//		})
//		It("ok", func() {
//			userRepo.On("ExistsByEmail", ctx, email).Return(true, nil)
//			userRepo.On("GetByEmail", ctx, email).Return(repoUser, nil)
//			authCache.On("CreateCode", ctx, email, code).Return(nil)
//			smtpClient.On("Send", email, msg).Return(nil)
//
//			err := authService.SendCode(ctx, email)
//
//			Ω(err).ShouldNot(HaveOccurred())
//		})
//		It("err user not found by email", func() {
//			userRepo.On("ExistsByEmail", ctx, email).Return(false, nil)
//
//			err := authService.SendCode(ctx, email)
//
//			Ω(err).Should(HaveOccurred())
//		})
//		It("user deleted", func() {
//			repoUser.DeletedAt.Valid = true
//			userRepo.On("ExistsByEmail", ctx, email).Return(true, nil)
//			userRepo.On("GetByEmail", ctx, email).Return(repoUser, nil)
//
//			err := authService.SendCode(ctx, email)
//
//			Ω(err).Should(HaveOccurred())
//		})
//	})
//
//	Context("verify code", func() {
//		var email string
//		var code int
//		var repoUser *models.User
//		BeforeEach(func() {
//			email = fake.Email()
//			repoUser = &models.User{
//				ID:        fake.Int64(),
//				Email:     email,
//				CreatedAt: time.Now(),
//				DeletedAt: sql.NullTime{
//					Valid: false,
//				},
//				Role: "user",
//			}
//			code = generateCode()
//		})
//		It("ok", func() {
//			authCache.On("GetCode", ctx, email).Return(code, nil)
//			authCache.On("DeleteCode", ctx, email).Return(nil)
//			userRepo.On("GetByEmail", ctx, email).Return(repoUser, nil)
//
//			token, err := authService.VerifyCode(ctx, email, code)
//
//			Ω(err).ShouldNot(HaveOccurred())
//			Ω(token).ShouldNot(Equal(""))
//		})
//	})
//
//	Context("validate token", func() {
//		var token string
//		var repoUser *models.User
//		BeforeEach(func() {
//			repoUser = &models.User{
//				ID:        fake.Int64(),
//				Email:     fake.Email(),
//				CreatedAt: time.Now(),
//				DeletedAt: sql.NullTime{
//					Valid: false,
//				},
//				Role: "user",
//			}
//			token, _ = generateToken(repoUser.ID, repoUser.Role, jwtSecret, defExpiration)
//		})
//		It("ok", func() {
//			userRepo.On("Get", ctx, repoUser.ID).Return(repoUser, nil)
//
//			res, err := authService.ValidateToken(ctx, token)
//
//			Ω(err).ShouldNot(HaveOccurred())
//			Ω(res.ID).Should(Equal(repoUser.ID))
//			Ω(res.Role).Should(Equal(repoUser.Role))
//		})
//		It("error user deleted", func() {
//			repoUser.DeletedAt.Valid = true
//			userRepo.On("Get", ctx, repoUser.ID).Return(repoUser, nil)
//
//			res, err := authService.ValidateToken(ctx, token)
//
//			Ω(err).Should(HaveOccurred())
//			Ω(res).Should(BeNil())
//		})
//		It("invalid token", func() {
//			token = fake.Name()
//			userRepo.On("Get", ctx, repoUser.ID).Return(repoUser, nil)
//
//			res, err := authService.ValidateToken(ctx, token)
//
//			Ω(err).Should(HaveOccurred())
//			Ω(res).Should(BeNil())
//		})
//	})
//})
