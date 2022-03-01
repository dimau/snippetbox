package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)

// Метод application для инициализации и настройки роутера
func (app *application) routes() http.Handler {
	mux := pat.New()
	// Все обработчики с динамическим контентом оборачиваем в middleware
	// для чтения/записи сессионных куки "app.session.Enable"
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", app.session.Enable(app.requireAuthentication(http.HandlerFunc(app.createSnippetForm))))
	mux.Post("/snippet/create", app.session.Enable(app.requireAuthentication(http.HandlerFunc(app.createSnippet))))
	mux.Get("/snippet/:id", app.session.Enable(http.HandlerFunc(app.showSnippet)))
	mux.Get("/user/signup", app.session.Enable(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", app.session.Enable(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", app.session.Enable(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", app.session.Enable(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", app.session.Enable(app.requireAuthentication(http.HandlerFunc(app.logoutUser))))

	// Обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Возвращает мультиплексор (роутер), обернутый в несколько слоев middleware обработчиков
	// Тем самым, сначала для каждого запроса последовательно отрабатывает логика каждого middleware
	// А затем уже отрабатывает логика непосредственно роутера и обработчика
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
