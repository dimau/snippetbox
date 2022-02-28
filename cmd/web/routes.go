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
	mux.Get("/snippet/create", app.session.Enable(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", app.session.Enable(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", app.session.Enable(http.HandlerFunc(app.showSnippet)))

	// Обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Возвращает мультиплексор (роутер), обернутый в несколько слоев middleware обработчиков
	// Тем самым, сначала для каждого запроса последовательно отрабатывает логика каждого middleware
	// А затем уже отрабатывает логика непосредственно роутера и обработчика
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
