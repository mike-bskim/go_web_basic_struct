package handlers

import (
	"GO/Toy_Prj/basic_struct/internal/config"
	"GO/Toy_Prj/basic_struct/internal/forms"
	"GO/Toy_Prj/basic_struct/internal/models"
	"GO/Toy_Prj/basic_struct/internal/render"
	"log"

	// _ "cycle"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

/*
	NewHandlers sets the repository for the handlers
	Repo *Repository 선언을 main에서 하면 쉽게 사용가능하지만,
	다른 파일들에서는 사용하기 불편함. 이유는 다른파일에서 main 을 패키지로 처리하고 접근해야 함.
	Repo *Repository 선언을 핸들러 내부에서 하면 다른 파일에서 패키지 처리하고 접근하기 편함.
	반대로 errors.go, forms.go 에서는 자신의 구조체(errors, Form)를
	(Repo *Repository) 처럼 내부에서 선언을 하지 않음
	이유는 각 화면(view)당 각각 에러 및 필드 검증을 하기 위해서 임
*/
func NewHandlers(r *Repository) {
	Repo = r
}

// Reservation renders the make a reservation page and displays form
// 초기에 get 에 대한 응답시, 오류 메시지도 없고 데이터도 없는 상태임.
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	// 응답 페이지에는 빈 데이터가 전달됨
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
// post 에 대한 응답시, 필드 데이터를 검증하여 오류시 오류 메시지 및 원래 데이터를 전달한다.
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// post 로 전달된 데이터 저장, 아래에서 응답시 화면으로 전달함
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	log.Println("reservation >>", reservation)

	// 각 필드데이터에 접근은 아래처럼 해도 된다.
	form := forms.New(r.PostForm)
	log.Println("handlers.go >>> forms.New:", form.Get("first_name"))

	// form.Has("first_name")
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() { // 해당 폼에 오류 메시지가 하나도 없어야 실행됨
		log.Println("!form.Valid()")
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
}
