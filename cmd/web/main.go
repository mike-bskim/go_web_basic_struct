package main

import (
	"GO/Toy_Prj/basic_struct/internal/config"
	"GO/Toy_Prj/basic_struct/internal/handlers"
	"GO/Toy_Prj/basic_struct/internal/render"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.AppConfig        // from config.go
var session *scs.SessionManager // from package

// main is the main function
func main() {

	// change this to true when in production, 보안강화 적용안함.
	app.InProduction = false

	// 세션관련 설정
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // middleware.go > Secure 에서도 참조함
	app.Session = session

	tc, err := render.CreateTemplateCache() // tmpl 파일을 조립하여 메모리로 로딩
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false // false: (DEV mode)read cache everytime.

	repo := handlers.NewRepo(&app) // main에서 선언한 AppConfig 변수를 handlers.go 와 공유
	// main에서 선언한 repo 객체를 handlers 에 전달하여 Repo 객체와 repo 객체의 메모리 매핑.
	// 다른 파일에서 Repo를 통해서 handlers 내부 함수에 접근 가능함.  사용예시 routes.go
	handlers.NewHandlers(repo)

	render.NewTemplates(&app) // main에서 선언한 AppConfig 변수를 render.go 와 공유

	tmp := fmt.Sprintf("Staring application on port %s", portNumber)
	fmt.Println(tmp)

	svr := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = svr.ListenAndServe()
	log.Fatal(err)
}
