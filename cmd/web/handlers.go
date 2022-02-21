package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что пользователь обращался именно к корневой странице сайта
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Нам нужен путь к файлу с HTML шаблоном конкретной страницы - home.page.tmpl
	// Также нам нужен путь к файлу с общим лейаутом для всех страниц сайта - base.layout.tmpl
	// Некоторые части общего лейаута могут быть вынесены
	//    для удобства переиспользования в отдельный файл - footer.partial.tmpl
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Непосредственно парсим все нужные для формирования конкретной страницы файлы с шаблонами
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Выполняем файлы с шаблонами и отдаем конечную HTML страницу
	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
