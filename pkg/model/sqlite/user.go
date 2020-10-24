package sqlite

import (
	"database/sql"
	"strings"
	"umanagement/pkg/model"

	"golang.org/x/crypto/bcrypt"
)

// UserModel ...
type UserModel struct {
	DB *sql.DB
}

// GetUsers get all users from database.
func (m *UserModel) GetUsers() ([]*model.User, error) {
	rows, err := m.DB.Query("SELECT id,sn,name FROM user;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.SN, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Create add a new user
func (m *UserModel) Create(sn, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO user(sn,name,email,hashed_password)
VALUES(?,?,?,?)`
	_, err = m.DB.Exec(stmt, sn, name, email, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return model.ErrDuplicate
		}
	}
	return err
}

// Edit modify a user
func (m *UserModel) Edit(u *model.User) error {
	hashedPsw, err := bcrypt.GenerateFromPassword([]byte(u.HashedPassword), 12)
	if err != nil {
		return err
	}
	stmt, err := m.DB.Prepare("UPDATE user SET sn=?,name=?,hashed_password=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.SN, u.Name, string(hashedPsw), u.ID)
	return err
}

// Delete delete a user from database
func (m *UserModel) Delete(id int) error {
	stmt, err := m.DB.Prepare("DELETE FROM user WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

// GetUser method fetch details for a specific user
func (m *UserModel) GetUser(id int) (*model.User, error) {
	s := &model.User{}
	stmt := "SELECT id,sn,name FROM user WHERE id=?"
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&s.ID, &s.SN, &s.Name)
	if err == sql.ErrNoRows {
		return nil, model.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// Authenticate verify where a user exist with the user sn and password
// This will return the relevant user struct
func (m *UserModel) Authenticate(sn, password string) (*model.User, error) {
	user := &model.User{}
	row := m.DB.QueryRow(`SELECT id,hashed_password FROM user WHERE sn=?`, sn)
	err := row.Scan(&user.ID, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, model.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, model.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
