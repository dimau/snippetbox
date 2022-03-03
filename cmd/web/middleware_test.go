package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// ***** Initialization *****
	// Инициализируем httptest.ResponseRecorder, который будем использовать вместо http.ResponseWriter
	rr := httptest.NewRecorder()

	// Инициализируем http.Request, который будем использовать вместо реального HTTP запроса
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем mock HTTP handler-а, который будем передавать нашему middleware
	// Mockup очень простой: записывает Status code = 200 и "OK" в тело ответа
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// ***** Execution *****
	// Обращаемся к тестируемому middleware обработчику, передавая ему заготовленный HTTP handler.
	// Так как middleware возвращает в качестве результата своей работы новый HTTP handler
	// (обертку над переданным), мы можем обратиться к методу ServeHTTP у возвращенного значения
	// Нам важно проверить, как HTTP handler после middleware обертки будет работать с тестовым HTTP запросом
	secureHeaders(next).ServeHTTP(rr, r)

	// ***** Validation *****
	// Вызываем метод Result() у нашего http.ResponseRecorder, чтобы получить http.Response
	rs := rr.Result()

	// Проверяем, что middleware корректно установил заголовок "X-Frame-Options" в HTTP Response
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	// Проверяем, что middleware корректно установил заголовок "X-XSS-Protection" в HTTP Response
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	// Проверяем, что middleware корректно вызывал next HTTP handler
	// Для этого убеждаемся, что status code и тело HTTP ответа соответствуют ожиданиям
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
