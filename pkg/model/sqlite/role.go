package sqlite

import (
	"database/sql"
	"umanagement/pkg/model"
)

// RoleModel handle roles in database
type RoleModel struct {
	DB *sql.DB
}

// GetRoles get all roles in database.
func (m *RoleModel) GetRoles() ([]*model.Role, error) {
	rows, err := m.DB.Query("SELECT * FROM role")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*model.Role{}
	for rows.Next() {
		r := &model.Role{}
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID get role by id
func (m *RoleModel) GetRoleByID(id int) (*model.Role, error) {
	role := &model.Role{}
	err := m.DB.QueryRow("SELECT * FROM role WHERE id=?", id).Scan(
		&role.ID, &role.Name, &role.Description,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return role, nil
}

// CreateRole add a new role to database
func (m *RoleModel) CreateRole(role *model.Role) error {
	stmt, err := m.DB.Prepare("INSERT INTO role(name,description) VALUES(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	return err
}

// EditRole edit a current role
func (m *RoleModel) EditRole(r *model.Role) error {
	q := "UPDATE role SET name=?,description=? WHERE id=?"
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.Name, r.Description, r.ID)
	return err
}

// DeleteRole delete a role
func (m *RoleModel) DeleteRole(id int) error {
	stmt, err := m.DB.Prepare("DELETE FROM role WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
