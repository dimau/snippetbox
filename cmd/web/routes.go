package main

import (
	"net/http"
)

// Метод app для инициализации и настройки роутера
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Возвращает мультиплексор (роутер), обернутый в несколько слоев middleware обработчиков
	// Тем самым, сначала для каждого запроса последовательно отрабатывает логика каждого middleware
	// А затем уже отрабатывает логика непосредственно роутера и обработчика
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
