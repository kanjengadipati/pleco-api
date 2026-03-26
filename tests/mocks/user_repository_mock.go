package mocks

import (
	"errors"

	"go-auth-app/models"
)

type MockUserRepo struct {
	User *models.User

	FindByEmailErr error
	FindByIDErr    error
}

func (m *MockUserRepo) Create(user *models.User) error {
	m.User = user
	return nil
}

func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	if m.FindByEmailErr != nil {
		return nil, m.FindByEmailErr
	}
	if m.User == nil {
		return nil, errors.New("user not found")
	}
	return m.User, nil
}

func (m *MockUserRepo) FindByID(id uint) (*models.User, error) {
	if m.FindByIDErr != nil {
		return nil, m.FindByIDErr
	}
	if m.User == nil {
		return nil, errors.New("user not found")
	}
	return m.User, nil
}

func (m *MockUserRepo) Update(user *models.User) error {
	m.User = user
	return nil
}
