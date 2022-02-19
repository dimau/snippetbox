package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Используем стандартный пакет "flag"
	// Объявляем параметр "addr"
	// В первом аргументе - название ключа командной строки, чье значение парсим
	// Во втором аргументе - значение по-умолчанию (если ключ не указан)
	// В третьем аргументе - короткое текстовое объяснение, за что флаг отвечает
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Непосредственно выполняем парсинг всех флагов в соответствующие переменные
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", *addr)

	// Значение, возвращаемое функцией flag.String(), является указателем на
	//    текстовое значение флага, а не самим значением.
	// Поэтому, чтобы получить само текстовое значение, нам нужно добавить перед указателем символ *
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
