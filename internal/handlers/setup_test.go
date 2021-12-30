package handlers

import (
	"GO/Toy_Prj/basic_struct/internal/config"
	"GO/Toy_Prj/basic_struct/internal/models"
	"GO/Toy_Prj/basic_struct/internal/render"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig        // from config.go
var session *scs.SessionManager // from package
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	// os.Exit(m.Run())
	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}

func getRoutes() http.Handler {

	// what am I doing to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production, 보안강화 적용안함.
	app.InProduction = false

	// 세션관련 설정
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // middleware.go > Secure 에서도 참조함
	app.Session = session

	tc, err := CreateTestTemplateCache() // tmpl 파일을 조립하여 메모리로 로딩
	fmt.Println(tc)
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true // render.CreateTemplateCache 사용 못하게 막아야 함.

	repo := NewRepo(&app) // main에서 선언한 AppConfig 변수를 handlers.go 와 공유
	// main에서 선언한 repo 객체를 handlers 에 전달하여 Repo 객체와 repo 객체의 메모리 매핑.
	// 다른 파일에서 Repo를 통해서 handlers 내부 함수에 접근 가능함.  사용예시 routes.go
	NewHandlers(repo)

	render.NewTemplates(&app) // main에서 선언한 AppConfig 변수를 render.go 와 공유

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf) // cause error at POST page.
	mux.Use(SessionLoad)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf andds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	log.Println("call NoSurf")
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteDefaultMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every requests
func SessionLoad(next http.Handler) http.Handler {
	log.Println("call SessionLoad")
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map, page + base
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	// 폴더이름+파일명 저장
	pages, err := filepath.Glob("./../../templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		// 폴더정보를 제거하고 파일 이름만 저장
		name := filepath.Base(page)

		// 페이지 정보 로딩
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob("./../../templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			// 페이지 정보에 base 정보를 추가 결함.
			ts, err = ts.ParseGlob("./../../templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}
