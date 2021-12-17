package render

import (
	"GO/Toy_Prj/basic_struct/internal/config"
	"GO/Toy_Prj/basic_struct/internal/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
// msin 에서 정의한 데이터 공유.
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
	// 아래 부분이 html 에서 불러오는 이름("{{.CSRFToken}}")과 같아야 함.
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
// TemplateData 를 이용하여 서버와 클라이언트 사이에 정보 공유함.
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		// 운영시에는 true 로 변경하여 tmpl 화면을 메모리에서 호출하여 속도를 빠르게 함.
		tc = app.TemplateCache
	} else {
		// DEV mode 에서는 UseCache==false 이므로. read cache everytime.
		// 매번 화면 호출시마다 tmpl 파일을 계속 읽음. 수정시 서버 재시작 필요없음
		tc, _ = CreateTemplateCache()
	}

	// map에 원하는 페이지가 있는지 확인
	t, ok := tc[tmpl]
	if !ok {
		// return errors.New("Could not get template from template cache")
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer) // buf 생성
	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td) // 해당 페이지를 buf 에 저장

	_, err := buf.WriteTo(w) // client 에게 전송
	if err != nil {
		// log.Println(err)
		fmt.Println("error writing template to browser", err)
		// return err
	}

}

// CreateTemplateCache creates a template cache as a map, page + base
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	// 폴더이름+파일명 저장
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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

		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			// 페이지 정보에 base 정보를 추가 결함.
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}
