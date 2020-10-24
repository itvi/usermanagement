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

// RoleHandler ...
type RoleHandler struct {
	M *sqlite.RoleModel
}

func (h *RoleHandler) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := h.M.GetRoles()
		if err != nil {
			log.Println(err)
			return
		}

		render(w, r, "./ui/html/role/index.html", roles)
	}
}

func (h *RoleHandler) details() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlID := strings.TrimPrefix(r.URL.Path, "/role/details/")
		id, err := strconv.Atoi(urlID)
		if err != nil {
			log.Println(err)
			return
		}

		role, err := h.M.GetRoleByID(id)
		if err == model.ErrNoRecord {
			log.Println(err)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		render(w, r, "./ui/html/role/details.html", role)
	}
}

func (h *RoleHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page = "./ui/html/role/create.html"
		if r.Method == "GET" {
			render(w, r, page, form.Init(nil))
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				return
			}

			form := form.Init(r.PostForm)
			form.Required("name")
			form.MaxLength("name", 20)
			form.MaxLength("desc", 50)

			if !form.Valid() {
				render(w, r, page, form)
				return
			}

			name := form.Get("name")
			desc := form.Get("desc")

			role := &model.Role{Name: name, Description: desc}
			err = h.M.CreateRole(role)
			if err != nil {
				log.Println(err)
				return
			}

			http.Redirect(w, r, "/roles", http.StatusSeeOther)
		}
	}
}

func (h *RoleHandler) edit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page = "./ui/html/role/edit.html"
		urlID := strings.TrimPrefix(r.URL.Path, "/role/edit/")
		id, err := strconv.Atoi(urlID)
		if err != nil {
			fmt.Fprint(w, "Not Found")
			return
		}

		if r.Method == "GET" {
			role, err := h.M.GetRoleByID(id)
			if err == model.ErrNoRecord {
				log.Println(err)
				return
			} else if err != nil {
				log.Println(err)
				return
			}
			render(w, r, page, role)
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				return
			}

			form := form.Init(r.PostForm)
			form.Required("name")
			form.MaxLength("name", 20)

			name := r.PostForm.Get("name")
			desc := r.PostForm.Get("desc")
			role := &model.Role{ID: id, Name: name, Description: desc}

			if !form.Valid() {
				render(w, r, page, &templateData{
					Form: form,
					Role: role,
				})
				return
			}

			err = h.M.EditRole(role)
			if err != nil {
				config.App.serverError(w, err)
				return
			}

			config.App.Session.Put(r, "flash", "更新成功")
			http.Redirect(w, r, "/roles", http.StatusSeeOther)
		}
	}
}

func (h *RoleHandler) delete(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil {
			config.App.notFound(w)
			return
		}
		err = h.M.DeleteRole(id)
		if err != nil {
			config.App.serverError(w, err)
			return
		}

		config.App.Session.Put(r, "flash", "删除成功")
	}
}
