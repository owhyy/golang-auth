package main

import (
	"fmt"
	"net/http"
	"owhyy/simple-auth/internal/models"
	"strconv"
	"strings"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	prefix, _, _ := strings.Cut(page, ".")
	data.BaseURL = prefix

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
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

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) newPagination(r *http.Request) (*paginationData, error) {
	var err error

	data := &paginationData{}

	curPage := 1
	if s := r.URL.Query().Get("page"); s != "" {
		curPage, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	data.CurrentPage = curPage

	perPage := 20
	if s := r.URL.Query().Get("per_page"); s != "" {
		perPage, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	if perPage == 20 || perPage == 40 || perPage == 60 {
		data.PerPage = perPage
	} else {
		data.PerPage = 20
	}

	data.Prev = curPage - 1
	data.Next = curPage + 1

	return data, nil
}
