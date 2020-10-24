package handler

import (
	"log"
	"net/http"
	"strconv"
)

// MyCasbinHandler ...
type MyCasbinHandler struct {
	M *sqlite.MyCasbinModel
}

func (h *MyCasbinHandler) index(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleName := r.URL.Query().Get("name")
		policies := h.M.GetPoliciesOrderBy(roleName)
		render(w, r, "./ui/html/casbin/index.html", &templateData{
			Role:           &model.Role{Name: roleName},
			CasbinPolicies: policies,
		})
	}
}

func (h *MyCasbinHandler) create(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			config.App.render(w, r, "casbin.policy.create.page.tmpl", &templateData{Form: forms.New(nil)})
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				config.App.clientError(w, http.StatusBadRequest)
				return
			}

			sub := r.PostFormValue("sub")
			obj := r.PostFormValue("obj")
			act := r.PostFormValue("act")

			ref := r.Header["Referer"]
			log.Println("ref:", ref)
			enforcer := h.M.InitCasbin()
			if enforcer == nil {
				log.Fatal("init casbin error")
			}
			_, err = enforcer.AddPolicy(sub, obj, act)
			if err != nil {
				customErrorPage(w, r, ErrorInfo{Err: err})
			}

			config.App.Session.Put(r, "flash", "添加成功")
			http.Redirect(w, r, "/casbin/policies", http.StatusSeeOther)
		}
	}
}

func (h *MyCasbinHandler) edit(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			log.Println("this is get method")
			sub := r.URL.Query().Get("sub")
			obj := r.URL.Query().Get("obj")
			act := r.URL.Query().Get("act")
			log.Println(sub, obj, act)
			policy := &model.CasbinPolicy{Sub: sub, Obj: obj, Act: act}
			config.App.render(w, r, "casbin.policy.edit.page.tmpl", &templateData{CasbinPolicy: policy})
			//http.Redirect(w, r, "/casbin/policies?sub="+sub+"&obj="+obj+"&act="+act, 303)
		}
		if r.Method == "POST" {
			// original data
			oSub := r.URL.Query().Get("sub")
			oObj := r.URL.Query().Get("obj")
			oAct := r.URL.Query().Get("act")
			log.Printf("old sub: %s, obj: %s, act: %s", oSub, oObj, oAct)

			// new data
			sub := r.PostFormValue("sub")
			obj := r.PostFormValue("obj")
			act := r.PostFormValue("act")
			log.Printf("new sub: %s, obj: %s, act: %s", sub, obj, act)

			// delete first then create
			enforcer := h.M.InitCasbin()
			if enforcer == nil {
				log.Fatal("init casbin error")
			}
			_, err := enforcer.RemovePolicy(oSub, oObj, oAct)
			if err != nil {
				customErrorPage(w, r, ErrorInfo{Err: err})
			}
			_, err = enforcer.AddPolicy(sub, obj, act)
			if err != nil {
				customErrorPage(w, r, ErrorInfo{Err: err})
			}

			config.App.Session.Put(r, "flash", "更新成功")
		}
	}
}

func (h *MyCasbinHandler) delete(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// from ajax form
		sub := r.PostFormValue("sub")
		obj := r.PostFormValue("obj")
		act := r.PostFormValue("act")
		log.Printf("new sub: %s, obj: %s, act: %s", sub, obj, act)

		enforcer := h.M.InitCasbin()
		if enforcer == nil {
			log.Fatal("init casbin error")
		}
		_, err := enforcer.RemovePolicy(sub, obj, act)
		if err != nil {
			customErrorPage(w, r, ErrorInfo{Err: err})
		}

		config.App.Session.Put(r, "flash", "删除成功")
	}
}

func (h *MyCasbinHandler) details(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get(":name")
		policies := h.M.GetPoliciesOrderBy(name)
		config.App.render(w, r, "casbin.policy.details.page.tmpl", &templateData{
			CasbinPolicies: policies,
		})
	}
}
func (h *MyCasbinHandler) assignRoles(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Convert user id from int to string
			id, err := strconv.Atoi(r.URL.Query().Get(":id"))
			if err != nil {
				config.App.notFound(w)
				return
			}

			// Get user by user id
			user, err := config.App.User.M.Get(id)
			if err != nil {
				config.App.serverError(w, err)
				return
			} else if err == model.ErrNoRecord {
				config.App.notFound(w)
				return
			}

			// Get all roles
			roles, err := config.App.Role.M.GetRoles()
			if err != nil {
				config.App.serverError(w, err)
				return
			}

			enforcer := h.M.InitCasbin()
			if enforcer == nil {
				log.Fatal("init casbin error")
			}
			//enforcer.LoadPolicy()

			rolesForUser, err := enforcer.GetRolesForUser(user.SN)

			if err != nil {
				log.Println("get roles err:", err)
			}

			config.App.render(w, r, "user.assign.page.tmpl", &templateData{
				User:                 user,
				Roles:                roles,
				RolesForSpecificUser: rolesForUser,
			})
		}
		if r.Method == "POST" {
			r.ParseForm()
			roles := r.Form["roles"]

			id, err := strconv.Atoi(r.URL.Query().Get(":id"))
			if err != nil {
				config.App.notFound(w)
				return
			}

			//user, err := config.App.Users.Get(id)
			user, err := config.App.User.M.Get(id)
			if err != nil {
				config.App.serverError(w, err)
				return
			} else if err == model.ErrNoRecord {
				config.App.notFound(w)
				return
			}

			enforcer := h.M.InitCasbin()
			if enforcer == nil {
				log.Fatal("init casbin error")
			}

			_, err = enforcer.DeleteRolesForUser(user.SN)
			if err != nil {
				log.Println("delete user :", err)
				customErrorPage(w, r, ErrorInfo{Err: err})
			}
			_, err = enforcer.AddRolesForUser(user.SN, roles)
			if err != nil {
				log.Println("add err:',err")
				customErrorPage(w, r, ErrorInfo{Err: err})
			}

			config.App.Session.Put(r, "flash", "分配角色成功")
			http.Redirect(w, r, "/users", http.StatusSeeOther)
		}
	}
}
