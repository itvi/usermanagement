// ViewModel

package model

import (
	"umanagement/pkg/form"
)

// UserEditModel ...
type UserEditModel struct {
	Form *form.UserForm
	User User
}

// RoleEditModel ...
type RoleEditModel struct {
	Form *form.UserForm
	Role Role
}

// CasbinIndexModel ...
type CasbinIndexModel struct {
	Role           *Role
	CasbinPolicies []*CasbinPolicy
}

// CasbinAddRolesForUserModel ...
type CasbinAddRolesForUserModel struct {
	User                 *User
	Roles                []*Role
	RolesForSpecificUser []string
}
