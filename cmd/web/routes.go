package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf)
	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/category/student", dynamicMiddleware.ThenFunc(app.category))
	mux.Get("/category/staff", dynamicMiddleware.ThenFunc(app.category))
	mux.Get("/category/applicant", dynamicMiddleware.ThenFunc(app.category))
	mux.Get("/category/research", dynamicMiddleware.ThenFunc(app.category))
	mux.Get("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNewsForm))
	mux.Post("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNews))

	mux.Get("/news/:id", dynamicMiddleware.ThenFunc(app.showNews))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))

	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
