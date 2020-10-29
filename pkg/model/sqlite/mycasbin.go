package sqlite

import (
	"database/sql"
	"log"
	"umanagement/pkg/model"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
)

// MyCasbinModel ...
type MyCasbinModel struct {
	DB *sql.DB
}

// InitCasbin initialize casbin
func (m *MyCasbinModel) InitCasbin() *casbin.Enforcer {
	tableName := "casbin"

	adpt, err := sqladapter.NewAdapter(m.DB, "sqlite3", tableName)

	if err != nil {
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer("./cmd/web/auth/perm.conf", adpt)
	if err != nil {
		panic(err)
	}

	return enforcer
}

// GetPolicies get policy by role name or get all policies
func (m *MyCasbinModel) GetPolicies(roleName string) []*model.CasbinPolicy {
	enforcer := m.InitCasbin()
	if enforcer == nil {
		log.Fatal("init casbin error")
	}

	var ps [][]string
	if roleName == "" {
		ps = enforcer.GetPolicy()
	} else {
		ps = enforcer.GetFilteredPolicy(0, roleName)
	}

	policies := []*model.CasbinPolicy{}
	for _, p := range ps {
		policy := &model.CasbinPolicy{
			Sub: p[0],
			Obj: p[1],
			Act: p[2],
		}
		policies = append(policies, policy)
	}
	return policies
}

// GetPoliciesOrderBy order by role name from database
func (m *MyCasbinModel) GetPoliciesOrderBy(roleName string) []*model.CasbinPolicy {
	var query string

	if roleName == "Administrator" || roleName == "" {
		query = `SELECT v0,v1,v2 FROM casbin WHERE p_type='p'
				 ORDER BY v0,v1`
	} else {
		query = `SELECT v0,v1,v2 FROM casbin 
				 WHERE p_type='p' AND v0='` + roleName + `'
				 ORDER BY v0,v1`
	}

	rows, err := m.DB.Query(query)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	policies := []*model.CasbinPolicy{}
	for rows.Next() {
		p := &model.CasbinPolicy{}
		if err := rows.Scan(&p.Sub, &p.Obj, &p.Act); err != nil {
			log.Println(err)
			return nil
		}
		policies = append(policies, p)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil
	}

	return policies
}
