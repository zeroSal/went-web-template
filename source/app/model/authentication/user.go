package authentication

import (
	"slices"

	"github.com/zeroSal/went-web/user"
)

const (
	RoleAdmin = "ROLE_ADMIN"
	RoleUser  = "ROLE_USER"
	RoleEditor = "ROLE_EDITOR"
)

var _ user.Interface = (*User)(nil)

type User struct {
	id      int
	username string
	password string
	roles    []string
}

func NewUser(
	id int,
	username string,
	password string,
	roles ...string,
) *User {
	return &User{
		id:      id,
		username: username,
		password: password,
		roles:    roles,
	}
}

func (u *User) GetID() any {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetRoles() []string {
	return u.roles
}

func (u *User) HasRole(role string) bool {
	return slices.Contains(u.roles, role)
}