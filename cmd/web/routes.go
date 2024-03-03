package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate, app.admin)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/genre/horror", dynamicMiddleware.ThenFunc(app.genre))
	mux.Get("/genre/comedy", dynamicMiddleware.ThenFunc(app.genre))
	mux.Get("/genre/drama", dynamicMiddleware.ThenFunc(app.genre))
	mux.Get("/genre/scifi", dynamicMiddleware.ThenFunc(app.genre))

	mux.Get("/movies/create", dynamicMiddleware.Append(app.requireAuthentication, app.requireAdmin).ThenFunc(app.createMoviesForm))
	mux.Post("/movies/create", dynamicMiddleware.Append(app.requireAuthentication, app.requireAdmin).ThenFunc(app.createMovies))
	mux.Post("/updateMovie", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateMovie))
	mux.Del("/movies/delete", dynamicMiddleware.Append(app.requireAuthentication, app.requireAdmin).ThenFunc(app.deleteMovie))
	mux.Get("/movies/:id", dynamicMiddleware.ThenFunc(app.showMovies))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
