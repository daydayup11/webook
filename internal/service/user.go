package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"webook/m/internal/domain"
	"webook/m/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

type UserService struct {
	repo *repository.UserRepository
}

// 初始化依赖注入
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		repo: userRepository,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	password := u.Password
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(fromPassword)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, err
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

// Edit 用户编译信息
func (svc *UserService) Edit(ctx *gin.Context, u domain.User) (domain.User, error) {
	return svc.repo.Update(ctx, u)
}
