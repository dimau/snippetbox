package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/Dimau/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	session *sessions.Session
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Обрабатываем конфигурационные параметры приложения
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32 bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// Инициализируем логгеры
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Инициализируем пул соединений с базой данных
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Инициализируем кэш шаблонов веб-страниц приложения
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true // Set the Secure flag on our session cookies

	// Инициализируем инстанс структуры application, который будет содержать все зависимости для handler-ов HTTP запросов
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		session: session,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want the server to use
	//
	// The tls.Config.PreferServerCipherSuites field controls whether the HTTPS connection should use
	// Go’s favored cipher suites or the user’s favored cipher suites.
	// By setting this to true — Go’s favored cipher suites are given preference
	// and we help increase the likelihood that a strong cipher suite which also supports forward secrecy is used
	//
	// The tls.Config.CurvePreferences field lets us specify which elliptic curves should be given preference
	// during the TLS handshake. Go supports a few elliptic curves, but as of Go 1.11 only
	// tls.CurveP256 and tls.X25519 have assembly implementations. The others are very CPU intensive,
	// so omitting them helps ensure that our server will remain performant under heavy loads.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Инициализация сервера и роутера на базе пакета net/http
	srv := &http.Server{
		Addr: *addr,                     // адрес и/или порт
		ErrorLog: errorLog,              // логгер для ошибок сервера
		Handler: app.routes(),           // что использовать в качестве handler запросов
		TLSConfig: tlsConfig,            // конфиги для TLS (HTTPS) соединения
		IdleTimeout: time.Minute,        // Таймаут сервера по всем запросам
		ReadTimeout: 5 * time.Second,    // Таймаут сервера по всем запросам
		WriteTimeout: 10 * time.Second,  // Таймаут сервера по всем запросам
	}

	// Запуск сервера на базе пакета net/http
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}