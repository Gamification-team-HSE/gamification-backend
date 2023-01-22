package user

//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/brianvoe/gofakeit/v6"
//	"github.com/go-playground/validator/v10"
//	. "github.com/onsi/ginkgo/v2"
//	. "github.com/onsi/gomega"
//	"gitlab.com/krespix/gamification-api/internal/core/config"
//	"gitlab.com/krespix/gamification-api/internal/models"
//	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
//	"testing"
//)
//
//func TestUserService(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "User service suite")
//}
//
//var _ = Describe("test user service", func() {
//	fake := gofakeit.New(0)
//
//	ctx := context.Background()
//	validate := validator.New()
//	userRepo := user.NewMockRepository(&testing.T{})
//
//	userService := New(userRepo, validate)
//
//	Context("test get user", func() {
//		var userID int64
//		var result *models.User
//		BeforeEach(func() {
//			userID = fake.Int64()
//			result = &models.User{ID: userID}
//		})
//		It("OK", func() {
//			userRepo.On("Get", ctx, userID).Return(result, nil)
//
//			res, err := userService.Get(ctx, userID)
//
//			Ω(err).ShouldNot(HaveOccurred())
//			Ω(res.ID).Should(Equal(userID))
//		})
//		It("err", func() {
//			userRepo.On("Get", ctx, userID).Return(nil, fmt.Errorf("error"))
//			_, err := userService.Get(ctx, userID)
//			Ω(err).Should(HaveOccurred())
//		})
//	})
//
//	Context("init super admin", func() {
//		var superAdmin config.SuperAdmin
//		BeforeEach(func() {
//			superAdmin = config.SuperAdmin{
//				Email: fake.Email(),
//				Name:  fake.Name(),
//			}
//		})
//		It("ok, created new super admin", func() {
//			userRepo.On("ExistsByEmail", ctx, superAdmin.Email).Return(false, nil)
//			userRepo.On("Create", ctx, &models.User{
//				Email: superAdmin.Email,
//				Role:  "super_admin",
//				Name: sql.NullString{
//					Valid:  true,
//					String: superAdmin.Name,
//				},
//			}).Return(nil)
//
//			err := userService.InitSuperAdmin(ctx, superAdmin)
//
//			Ω(err).ShouldNot(HaveOccurred())
//		})
//		It("ok, super admin already created", func() {
//			userRepo.On("ExistsByEmail", ctx, superAdmin.Email).Return(true, nil)
//
//			err := userService.InitSuperAdmin(ctx, superAdmin)
//
//			Ω(err).ShouldNot(HaveOccurred())
//		})
//	})
//
//	Context("list users", func() {
//		var users []*models.User
//		BeforeEach(func() {
//			users = []*models.User{
//				{
//					ID: fake.Int64(),
//				},
//				{
//					ID: fake.Int64(),
//				},
//				{
//					ID: fake.Int64(),
//				},
//			}
//		})
//		It("ok", func() {
//			userRepo.On("List", ctx).Return(users, nil)
//
//			res, err := userService.List(ctx)
//
//			Ω(err).ShouldNot(HaveOccurred())
//			Ω(len(res)).Should(Equal(len(users)))
//			Ω(res[0].ID).Should(Equal(users[0].ID))
//		})
//	})
//
//	Context("create user", func() {
//		var user *models.User
//		BeforeEach(func() {
//			user = &models.User{
//				Email: fake.Email(),
//				Role:  "admin",
//			}
//		})
//		It("ok", func() {
//			userRepo.On("ExistsByEmail", ctx, user.Email).Return(false, nil)
//			userRepo.On("Create", ctx, user).Return(nil)
//
//			err := userService.Create(ctx, user)
//
//			Ω(err).ShouldNot(HaveOccurred())
//		})
//		It("err already exists", func() {
//			userRepo.On("ExistsByEmail", ctx, user.Email).Return(true, nil)
//
//			err := userService.Create(ctx, user)
//
//			Ω(err).Should(HaveOccurred())
//		})
//	})
//})
