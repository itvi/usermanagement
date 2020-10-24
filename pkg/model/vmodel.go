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
