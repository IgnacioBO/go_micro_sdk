package mock

import (
	"errors"

	"github.com/IgnacioBO/gomicro_domain/domain"
)

type UserSdkMock struct {
	GetMock func(id string) (*domain.User, error)
}

func (m *UserSdkMock) Get(id string) (*domain.User, error) {
	// Esta validacion es para asegurarnos de que se setee el mock antes de usarlo (ya que debemos establecer que hara el mock)
	if m.GetMock == nil {
		return nil, errors.New("GetMock is not set")
	}
	return m.GetMock(id)
}
