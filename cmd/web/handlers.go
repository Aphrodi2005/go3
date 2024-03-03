package main

import (
	"AituNews/pkg/forms"
	"AituNews/pkg/models"
	"errors"
	"fmt"
	_ "github.com/golangcollege/sessions"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.movies.Latest(10)
	if err != nil {
		app.serverError(w, err)
		return
	}
	isAdmin := app.isAdmin(r)

	app.render(w, r, "home.page.tmpl", &templateData{
		Movie2s: s,
		IsAdmin: isAdmin,
	})
}

func (app *application) showMovies(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.movies.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{

		Movies: s,
	})
}

func (app *application) genre(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	genre := segments[len(segments)-1]

	s, err := app.movies.GetMovieByGenre(genre)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Movie2s: s})
}

func (app *application) createMoviesForm(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) || !app.isAdmin(r) {
		app.clientError(w, http.StatusForbidden)
		return
	}
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}

func (app *application) createMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	title := r.PostForm.Get("title")
	genre := r.PostForm.Get("genre")
	ratingStr := r.PostForm.Get("rating")
	sessionTimeStr := r.PostForm.Get("sessionTime")

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	sessionTime, err := time.Parse("2006-01-02T15:04", sessionTimeStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.movies.Create(title, genre, float64(rating), sessionTime)
	if errors.Is(err, models.ErrDuplicate) {
		app.clientError(w, http.StatusBadRequest)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}
func (app *application) updateMoviesForm(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) || !app.isAdmin(r) {
		app.clientError(w, http.StatusForbidden)
		return
	}
	app.render(w, r, "update.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}

func (app *application) updateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := strconv.Atoi(r.PostForm.Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	title := r.PostForm.Get("title")
	genre := r.PostForm.Get("genre")
	ratingStr := r.PostForm.Get("rating")
	sessionTimeStr := r.PostForm.Get("sessionTime")

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	sessionTime, err := time.Parse("2006-01-02T15:04", sessionTimeStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.movies.Update(title, genre, id, float64(rating), sessionTime)
	if errors.Is(err, models.ErrDuplicate) {
		app.clientError(w, http.StatusBadRequest)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.movies.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password", "role")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 3)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	role := form.Get("role")
	if role == "admin" {
		app.session.Put(r, "role", models.RoleAdmin)
	}
	if role == "supplier" {
		app.session.Put(r, "role", models.Rolesupplier)
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"), form.Get("role"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}
func (app *application) isAdmin(r *http.Request) bool {
	role := app.session.GetString(r, "role")
	return role == models.RoleAdmin
}
func (app *application) buyTicket(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
	if !app.isAuthenticated(r) {
		// If not authenticated, redirect to the login page
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get user ID from session
	userID := app.session.GetInt(r, "authenticatedUserID")

	// Get movie ID and other ticket details from the form
	movieID, err := strconv.Atoi(r.Form.Get("movie_id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	movieTitle := r.Form.Get("movie_title")
	sessionTimeStr := r.Form.Get("session_time")
	sessionTime, err := time.Parse("2006-01-02T15:04", sessionTimeStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new ticket
	err = app.tickets.Create(userID, movieID, movieTitle, sessionTime)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to a success page or any other appropriate page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
