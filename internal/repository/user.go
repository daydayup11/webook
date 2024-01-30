package repository

import (
	"context"
	"webook/m/internal/domain"
	"webook/m/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrDataNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(d *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: d,
	}
}

// 调用dao层接口
func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	//操作缓存。。。。
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (ur *UserRepository) FindById(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := ur.dao.FindById(ctx, id)
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (ur *UserRepository) Update(ctx context.Context, u domain.User) (domain.User, error) {
	daoUser := dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
	updatedUser, err := ur.dao.UpdateById(ctx, daoUser)
	return domain.User{
		Id:       updatedUser.Id,
		Email:    updatedUser.Email,
		Password: updatedUser.Password,
	}, err
}
