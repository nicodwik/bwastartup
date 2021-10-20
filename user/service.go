package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	CheckEmailAvailability(input CheckInputEmail) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(password) // data berupa byte, perlu convert string
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil

	// mapping struct input ke struct User
	// simpan struct User melalui repository

}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return user, errors.New("User not found")
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) // return error automatically

		if err != nil {
			return user, err
		} else {
			return user, nil
		}
	}

}

func (s *service) CheckEmailAvailability(input CheckInputEmail) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.Id == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {

	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	user.Avatar = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return user, err
	}

	return updatedUser, nil
}
