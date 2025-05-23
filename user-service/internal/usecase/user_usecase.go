package usecase

import (
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserUsecase interface {
	Register(user *model.User) error
	Login(email, password string) (model.User, error)
	GetProfile(id int) (model.User, error)
	UpdateUser(user model.User) error
	DeleteUser(id int) error
	Verify(token string) error
}

// Mailer умеет отправлять письма с верификацией
type Mailer interface {
	SendVerification(email, token string) error
}

type userUsecase struct {
	repo   repository.UserRepository
	mailer Mailer
}

func NewUserUsecase(repo repository.UserRepository, mailer Mailer) UserUsecase {
	return &userUsecase{
		repo:   repo,
		mailer: mailer,
	}
}

func (u *userUsecase) Register(user *model.User) error {
	// 1. Хэшируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	// 2. Генерируем токен
	token := uuid.NewString()

	// 3. Сохраняем в pending_users
	if err := u.repo.CreatePendingUser(user, token); err != nil {
		return err
	}

	// 4. Отправляем письмо
	return u.mailer.SendVerification(user.Email, token)
}

func (u *userUsecase) Login(email, password string) (model.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return model.User{}, bcrypt.ErrMismatchedHashAndPassword
	}
	return user, nil
}

func (u *userUsecase) GetProfile(id int) (model.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUsecase) UpdateUser(user model.User) error {
	return u.repo.UpdateUser(user)
}

func (u *userUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}

func (u *userUsecase) Verify(token string) error {
	// 1. Достаём pending по токену
	pu, err := u.repo.GetPendingByToken(token)
	if err != nil {
		return status.Errorf(codes.NotFound, "invalid token")
	}
	// 2. Создаём реального пользователя
	full := &model.User{
		Email:    pu.Email,
		Name:     pu.Name,
		Password: pu.Password, // уже захэширован
	}
	if err := u.repo.CreateUser(full); err != nil {
		return err
	}
	// 3. Удаляем pending
	return u.repo.DeletePending(token)
}
