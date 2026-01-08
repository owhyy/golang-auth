package main

import (
	"net/http"
	"owhyy/simple-auth/internal/models"
	"owhyy/simple-auth/internal/types"
	"owhyy/simple-auth/ui/templates"
	"strconv"

	"github.com/a-h/templ"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, title string, main templ.Component) {
	w.WriteHeader(status)

	navComponent := templates.Nav(r.URL.Path, app.isAuthenticated(r))
	templates.Base(title, navComponent, main).Render(r.Context(), w)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.errorLog.Println(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	session, err := app.cookieStore.Get(r, "auth-session")
	if err != nil {
		return false
	}
	return session.Values["userID"] != nil && !session.IsNew
}

func (app *application) getAuthenticatedUser(r *http.Request) *models.User {
	session, err := app.cookieStore.Get(r, "auth-session")
	if err != nil {
		return nil
	}
	id, ok := session.Values["userID"].(uint)
	if !ok || session.IsNew {
		return nil
	}

	user, err := app.users.GetByID(id)
	if err != nil {
		return nil
	}
	return user
}

func (app *application) renderHTMXSuccess(w http.ResponseWriter, msg string) {
	w.Write([]byte(`<p class="pico-color-green-600">` + msg + "</p>"))
}

func (app *application) renderHTMXError(w http.ResponseWriter, msg string) {
	w.Write([]byte(`<p class="pico-color-red-600">` + msg + "</p>"))
}

func (app *application) newPagination(r *http.Request) (*types.PaginationData, error) {
	var err error

	data := &types.PaginationData{}

	curPage := 1
	if s := r.URL.Query().Get("page"); s != "" {
		curPage, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	data.CurrentPage = curPage

	// TODO: implement this on frontend
	perPage := 30
	if s := r.URL.Query().Get("per_page"); s != "" {
		perPage, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	if perPage == 30 || perPage == 60 || perPage == 90 {
		data.PerPage = perPage
	} else {
		data.PerPage = 30
	}

	data.Prev = curPage - 1
	data.Next = curPage + 1

	return data, nil
}
