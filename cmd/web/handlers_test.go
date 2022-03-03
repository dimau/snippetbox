package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// ***** Initialization *****
	// Инициализируем httptest.ResponseRecorder, который будем в тесте использовать вместо http.ResponseWriter
	rr := httptest.NewRecorder()

	// Инициализируем http.Request для использования в тесте
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ***** Execution *****
	// Вызываем тестируемый handler (обработчик HTTP запросов) - в нашем случае функция ping
	ping(rr, r)

	// ***** Validation *****
	// Вызываем метод Result() у нашего http.ResponseRecorder, чтобы получить http.Response,
	// сгенерированный проверяемым в тесте обработчиком.
	// То есть по сути мы получаем и далее проверяем HTTP ответ, который выдал наш handler
	rs := rr.Result()

	// Проверяем Status Code ответа, полученного от handler. Мы ожидаем = 200 "OK"
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// Проверяем тело HTTP ответа от handler. Мы ожидаем его = "OK"
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
