package provider

import (
	"fmt"
	"webtemplate/app/model/authentication"

	"github.com/zeroSal/went-web/session"
	"github.com/zeroSal/went-web/user"
)

var _ session.ProviderInterface = (*Session)(nil)

type Session struct{}

func NewSession() session.ProviderInterface {
	return &Session{}
}

func (p *Session) Load(credentials any) (user.Interface, error) {
	username, ok := credentials.(string)
	if !ok {
		return nil, fmt.Errorf("invalid credentials type: %T", credentials)
	}

	if username == "" {
		return nil, fmt.Errorf("empty credentials provided")
	}

	return authentication.NewUser(1, "admin", "password", "ROLE_ADMIN"), nil
}

