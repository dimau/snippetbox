package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	// Если в middleware пришел nil вместо данных для отображения в шаблоне, мы все-равно должны
	// создать структуру templateData и добавить в нее данные по умолчанию
	if td == nil {
		td = &templateData{}
	}

	// Текущий год мы используем в подвале каждом страницы сайта
	td.CurrentYear = time.Now().Year()

	// Add the flash message to the template data, if one exists.
	// Если обработчик HTTP запроса добавлял в сессию куку с flash сообщением
	// (которое нужно показать только один раз - на следующей странице пользователю)
	// То достаем это flash сообщение в templateData для отрисовки при рендеринге шаблона страницы
	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	td.Flash = app.session.PopString(r, "flash")

	// Add the authentication status to the template data
	// It's useful for rendering of almost each page of the site
	td.IsAuthenticated = app.isAuthenticated(r)

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// По имени файла-шаблона страницы (например, 'home.page.tmpl') достаем из кэша шаблонов
	//   весь набор необходимых для ее рендеринга файлов с шаблонами (template set)
	// Если не нашли в кэше соответствующий набор шаблонов, будем отвечать ошибкой сервера
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Инициализируем буфер, в который отрендерим HTML страницу перед отправкой клиенту
	// Так как у нас есть динамический контент на странице, при рендеринге в runtime могут быть ошибки.
	// Чтобы отловить ошибку рендеринга и не отправлять клиенту некорректный заголовок "200 OK" с половиной страницы
	// Мы сначала будем рендерить страницу в буфер, а затем (если нет ошибок) отправлять ее клиенту
	buf := new(bytes.Buffer)

	// Пытаемся отрендерить HTML страницу с динамическим контентом, результат пишем в буфер
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Если рендеринг HTML страницы прошел успешно, пишем содержимое буфера в http.ResponseWriter клиенту
	buf.WriteTo(w)
}

// Return true if the current request is from authenticated user, otherwise return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}