package model

import (
	"errors"
)

// Role
type Role string

const (
	RoleInvalid  Role = ""
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

var (
	roles = []Role{RoleAdmin, RoleCustomer}
)

func (r *Role) UnmarshalText(b []byte) error {
	for _, v := range roles {
		if string(v) == string(b) {
			*r = v
			return nil
		}
	}
	return errors.New("invalid role")
}
