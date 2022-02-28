package main

import (
	"github.com/Dimau/snippetbox/pkg/forms"
	"github.com/Dimau/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear int
	Form *forms.Form
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

// Create a humanDate function which returns a nicely formatted string representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализируем map для хранения кэша шаблонов веб-приложения
	cache := map[string]*template.Template{}

	// С помощью функции filepath.Glob получаем массив путей ко всем файлам с расширением '.page.tmpl'
	// По сути массив всех шаблонов страниц веб-приложения
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Проходим в цикле по каждой странице
	for _, page := range pages {
		// Достаем имя файла (например 'home.page.tmpl') из полного пути к файлу
		name := filepath.Base(page)

		// Парсим соответствующий файл с шаблоном страницы в template set
		// Функции для использования в шаблонах объявили заранее в глобальной переменной functions
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// С помощью метода ParseGlob добавляем в template set шаблоны всех layout-ов
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// С помощью метода ParseGlob добавляем в template set шаблоны всех partial файлов
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Добавляем полученный template set в map, который будет служить кэшем шаблонов
		// В качестве ключа будет выступать название шаблона страницы (например, 'home.page.tmpl')
		cache[name] = ts
	}

	// Возвращаем заготовленный кэш шаблонов
	return cache, nil
}
