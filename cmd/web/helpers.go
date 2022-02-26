package main

import (
	"bytes"
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
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Если рендеринг HTML страницы прошел успешно, пишем содержимое буфера в http.ResponseWriter клиенту
	buf.WriteTo(w)
}

