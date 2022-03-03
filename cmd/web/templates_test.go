package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Создаем таблицу sub-tests (в рамках одного большого тест-кейса) - на основе массива структур
	// Каждая структура представляет собой описание одного sub-test
	// В поле name будем хранить имя sub-test, в поле tm - входные данные, в поле want - ожидаемый результат
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		}, {
			name: "Empty",
			tm:   time.Time{},
			want: "",
		}, {
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	// Проходим в цикле по каждому sub-test
	for _, tt := range tests {
		// Для запуска каждого sub-test внутри одного тест-кейса используется t.Run()
		// Первый аргумент - название sub-test (чтобы идентифицировать его в логах/выводе)
		// Второй аргумент - функция, которая фактически содержит описание sub-test
		t.Run(tt.name, func(t *testing.T) {
			// Initialization and execution of sub-test
			hd := humanDate(tt.tm)
			// Validation of sub-test
			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}
