package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	// Получаем полный stack trace ошибки в текущей горутине
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// Output с calldepth=2 вместо обычного Print позволяет вывести stack trace ошибки без учета верхушки стека
	// Верхушка стека бесполезна = названию данного файла и данной строки (так как отсюда производится запись в лог)
	app.errorLog.Output(2, trace)

	// Отправляем ответ на запрос со статусом = 500 (Internal Server Error)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description to the user
// when there's a problem with the request that the user sent
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}