package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/mlkad/golang-todoapp/internal/core/errors"
)

type User struct {
	ID 			int
	Version int
	
	FullName string
	PhoneNumber *string
} 

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID: id,
		Version: version,
		FullName: fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)
}

//это валидация перед сохранением в БД
func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 1 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}
		
		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	return nil
}

/*
HTTP слой (CreateUserRequest)
  ↓
  domainFromDTO() или NewUserUninitialized()
  ↓
Domain слой (domain.User с ID=UninitializedID)
  ↓
Service слой → Сохраняет в БД
  ↓
БД возвращает реальный ID и Version
  ↓
Domain слой (domain.User с реальным ID и Version)
*/