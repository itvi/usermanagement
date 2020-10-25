package handler

import "net/http"

// Application ...
type Application struct {
	Handlers *Handlers
}

// Handlers ...
type Handlers struct {
	User *UserHandler
	Role *RoleHandler
	// ...
	// ...
	// ...
}

// Routes ...
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/users", app.Handlers.User.index())
	mux.Handle("/user/create", app.Handlers.User.create())
	mux.Handle("/user/edit/", app.Handlers.User.edit())
	mux.Handle("/user/delete/", app.Handlers.User.delete())
	mux.Handle("/user/details/", app.Handlers.User.details())
	mux.Handle("/user/login", app.Handlers.User.login())
	mux.Handle("/user/logout", app.Handlers.User.logout())

	mux.Handle("/roles", app.Handlers.Role.index())
	return mux
}
