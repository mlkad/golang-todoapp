package users_service

import (
	"context"

	"github.com/mlkad/golang-todoapp/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

//тут будут методы(что может делать)
type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

/*
┌────────────────────────────┐
│   HTTP слой (handler)      │
│  (CreateUser HTTP)         │
└────────────┬───────────────┘
             ↓
┌────────────────────────────┐
│   Service слой (бизнес)    │  ← ВЫ ЗДЕСЬ
│  (CreateUser логика)       │  (users_service)
│  - Валидация данных        │
│  - Бизнес-правила          │
│  - Использует репозиторий  │
└────────────┬───────────────┘
             ↓
┌────────────────────────────┐
│   Repository слой (БД)     │
│  (CreateUser в БД)         │
│  - SQL запросы             │
│  - Работа с БД             │
└────────────────────────────┘
*/