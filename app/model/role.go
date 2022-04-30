package model

import (
	"errors"
)

// Role
type Role string

const (
	RoleInvalid  Role = ""
	RoleAdmin    Role = "admin"
	RoleConsumer Role = "consumer"
)

var (
	roles = []Role{RoleAdmin, RoleConsumer}
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
