package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"umanagement/pkg/form"
	"umanagement/pkg/model"
	"umanagement/pkg/model/sqlite"
)

// MyCasbinHandler ...
type MyCasbinHandler struct {
	M *sqlite.MyCasbinModel
}

func (h *MyCasbinHandler) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/casbin/index/")
		policies := h.M.GetPoliciesOrderBy(name)
		render(w, r, "./ui/html/casbin/index.html", nil, &model.CasbinIndexModel{
			Role:           &model.Role{Name: name},
			CasbinPolicies: policies,
		})
	}
}

func (h *MyCasbinHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			render(w, r, "./ui/html/casbin/create.html", nil, form.Init(nil))
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
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
				log.Println(err)
			}

			// Save the policy back to DB
			if err = enforcer.SavePolicy(); err != nil {
				log.Println("Save Policy failed, err:", err)
				return
			}

			http.Redirect(w, r, "/casbin/index", http.StatusSeeOther)
		}
	}
}

func (h *MyCasbinHandler) edit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			log.Println("this is get method")
			sub := r.URL.Query().Get("sub")
			obj := r.URL.Query().Get("obj")
			act := r.URL.Query().Get("act")
			log.Println(sub, obj, act)
			policy := &model.CasbinPolicy{Sub: sub, Obj: obj, Act: act}
			render(w, r, "./ui/html/casbin/edit.html", nil, policy)
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
				log.Println(err)
			}
			_, err = enforcer.AddPolicy(sub, obj, act)
			if err != nil {
				log.Println(err)
			}

			// Save the policy back to DB
			if err = enforcer.SavePolicy(); err != nil {
				log.Println("Save Policy failed, err:", err)
				return
			}
			http.Redirect(w, r, "/casbin/index", http.StatusSeeOther)
		}
	}
}

func (h *MyCasbinHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub := r.URL.Query().Get("sub")
		obj := r.URL.Query().Get("obj")
		act := r.URL.Query().Get("act")
		log.Printf("sub: %s, obj: %s, act: %s", sub, obj, act)

		enforcer := h.M.InitCasbin()
		if enforcer == nil {
			log.Fatal("init casbin error")
		}
		_, err := enforcer.RemovePolicy(sub, obj, act)
		if err != nil {
			log.Println(err)
		}
		log.Println("Deleted")
		http.Redirect(w, r, "/casbin/index", http.StatusSeeOther)
	}
}

func (h *MyCasbinHandler) details() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/casbin/details/")
		policies := h.M.GetPoliciesOrderBy(name)
		render(w, r, "./ui/html/casbin/details.html", nil, policies)
	}
}

func (h *MyCasbinHandler) addRolesForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page = "./ui/html/casbin/add_roles_for_user.html"
		urlID := strings.TrimPrefix(r.URL.Path, "/casbin/addr2u/")
		// Convert user id from int to string
		id, err := strconv.Atoi(urlID)
		if err != nil {
			fmt.Fprint(w, "Not Found")
			return
		}

		// Get user by user id
		db := h.M.DB
		userModel := &sqlite.UserModel{DB: db}
		user, err := userModel.GetUser(id)

		if err != nil {
			log.Println(err)
			return
		} else if err == model.ErrNoRecord {
			log.Println(err)
			return
		}

		enforcer := h.M.InitCasbin()
		if enforcer == nil {
			log.Fatal("init casbin error")
			return
		}

		if r.Method == "GET" {
			// Get all roles
			roleModel := &sqlite.RoleModel{DB: db}
			roles, err := roleModel.GetRoles()
			if err != nil {
				log.Println(err)
				return
			}

			rolesForUser, err := enforcer.GetRolesForUser(user.SN)

			if err != nil {
				log.Println("get roles err:", err)
			}

			funcMap := template.FuncMap{
				"rolesChecked": rolesChecked,
				"safe":         safe,
			}
			render(w, r, page, funcMap, &model.CasbinAddRolesForUserModel{
				User:                 user,
				Roles:                roles,
				RolesForSpecificUser: rolesForUser,
			})
		}
		if r.Method == "POST" {
			r.ParseForm()
			roles := r.Form["roles"]

			_, err = enforcer.DeleteRolesForUser(user.SN)
			if err != nil {
				log.Println("delete user :", err)
			}

			// interface conversion: *sqladapter.Adapter is not persist.BatchAdapter: missing method AddPolicies

			// implement these methods in sqladapter's adapter.go:
			/*
				func (p *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
					return nil
				}
				func (p *Adapter) RemovePolicies(sec string, ptype string, rule [][]string) error {
					return nil
				}
			*/

			_, err = enforcer.AddRolesForUser(user.SN, roles) // interface conversion: *sqladapter.Adapter is not persist.BatchAdapter: missing method AddPolicies
			if err != nil {
				log.Println("add err:',err")
			}

			// Save the policy back to DB
			if err = enforcer.SavePolicy(); err != nil {
				log.Println("Save Policy failed, err:", err)
				return
			}

			http.Redirect(w, r, "/casbin/index", http.StatusSeeOther)
		}
	}
}
