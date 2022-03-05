package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// ***** Initialization *****
	// Создаем инстанс application структуры, которая хранит все общие зависимости уровня приложения.
	// Не обязательно инициализировать все поля этой структуры (зависимости приложения),
	// достаточно только те, без которых не сработает цепочка обработчиков, задействованных в обработке
	// данного конкретного тестового запроса (в данном примере достаточно mock-а только для логгеров)
	app := &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
	}

	// Инициализируем тестовый HTTPS сервер (цепляется к случайному порту на машине на время проведения теста),
	// который будет в качестве обработчика для всех запросов использовать наш роутер - из app.routes()
	// (соответственно, мы сможем проверить работу приложения от этапа роутинга и до выдачи ответа на запрос
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// ***** Execution *****
	// С помощью поля ts.URL получаем фактический адрес и порт, по которому слушает тестовый сервер.
	// Собираем HTTP запрос с методом GET и path="/ping" и отправляем его в тестовый сервер.
	// Записываем полученную в ответ структуру типа http.Response в переменную rs
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	// ***** Validation *****
	// Убеждаемся, что HTTP ответ содержит ожидаемый Status Code, в данном случае = 200 "OK"
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// Достаем тело HTTP ответа и убеждаемся, что тело HTTP ответа = "OK"
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
