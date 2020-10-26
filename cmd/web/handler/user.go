package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"umanagement/pkg/form"
	"umanagement/pkg/model"
	"umanagement/pkg/model/sqlite"
)

// UserHandler ...
type UserHandler struct {
	M *sqlite.UserModel
}

func (h *UserHandler) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.M.GetUsers()
		if err != nil {
			log.Println(err)
			return
		}
		render(w, r, "./ui/html/user/index.html", users)
	}
}

// create a new user
func (h *UserHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page = "./ui/html/user/create.html"

		if r.Method == "GET" {
			render(w, r, page, form.Init(nil))
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				return
			}
			userForm := form.Init(r.PostForm)
			userForm.Required("sn", "name", "password")
			userForm.Match("email", form.EmailReg)
			userForm.MinLength("password", 3)
			userForm.MaxLength("sn", 8)
			userForm.MaxLength("name", 6)

			if !userForm.Valid() {
				render(w, r, page, userForm)
				return
			}

			err = h.M.Create(userForm.Get("sn"), userForm.Get("name"),
				userForm.Get("email"), userForm.Get("password"))
			if err == model.ErrDuplicate {
				userForm.Errors.Add("sn", "user already exist")
				render(w, r, page, userForm)
				return
			} else if err != nil {
				log.Println(err)
				return
			}

			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	}
}

func (h *UserHandler) edit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page = "./ui/html/user/edit.html"
		urlID := strings.TrimPrefix(r.URL.Path, "/user/edit/")
		id, err := strconv.Atoi(urlID)
		if err != nil {
			fmt.Fprint(w, "Not Found")
			return
		}

		if r.Method == "GET" {
			user, err := h.M.GetUser(id)
			if err == model.ErrNoRecord {
				fmt.Fprint(w, "Not Found")
				return
			} else if err != nil {
				log.Println(err)
				return
			}

			userForm := form.Init(r.PostForm)
			data := model.UserEditModel{
				Form: userForm,
				User: *user,
			}
			render(w, r, page, &data)
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				return
			}
			userForm := form.Init(r.PostForm)
			userForm.Required("name", "hashedPassword")
			userForm.MaxLength("name", 6)

			sn := r.PostForm.Get("sn")
			name := r.PostForm.Get("name")
			psw := r.PostForm.Get("hashedPassword")

			user := &model.User{ID: id, SN: sn, Name: name, HashedPassword: []byte(psw)}

			if !userForm.Valid() {
				data := model.UserEditModel{
					Form: userForm,
					User: *user,
				}
				render(w, r, page, &data)
				return
			}

			if err != h.M.Edit(user) {
				log.Println(err)
				return
			}

			http.Redirect(w, r, "/users", http.StatusSeeOther)
		}
	}
}

func (h *UserHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlID := strings.TrimPrefix(r.URL.Path, "/user/delete/")
		id, err := strconv.Atoi(urlID)
		if err != nil {
			fmt.Fprint(w, "Not Found")
			return
		}

		err = h.M.Delete(id)
		if err != nil {
			fmt.Fprint(w, "Not Found")
			return
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
}

func (h *UserHandler) details() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/user/details/")

		fmt.Fprintln(w, id)
	}
}

func (h *UserHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page = "./ui/html/user/login.html"

		if r.Method == "GET" {
			render(w, r, page, form.Init(nil))
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				return
			}

			userForm := form.Init(r.PostForm)
			user, err := h.M.Authenticate(userForm.Get("sn"), userForm.Get("password"))
			if err == model.ErrInvalidCredentials {
				userForm.Errors.Add("generic", "User name or password is not correct!")
				render(w, r, page, userForm)
				return
			} else if err != nil {
				log.Println(err)
				return
			}

			// TODO: add user to session

			log.Println(user.ID)

			http.Redirect(w, r, "/users", http.StatusSeeOther)
		}
	}
}

func (h *UserHandler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO: remove session

		http.Redirect(w, r, "/", 303)
	}
}
