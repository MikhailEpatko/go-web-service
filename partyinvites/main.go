package main

import (
	"html/template"
	"log"
	"net/http"
)

type Rsvp struct {
	Name,
	Email,
	Phone string
	WillAttend bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{
		"welcome",
		"form",
		"thanks",
		"sorry",
		"list",
	}
	for i, name := range templateNames {
		if t, err := template.ParseFiles(
			"templates/"+"layout"+".html",
			"templates/"+name+".html",
		); err != nil {
			panic(err)
		} else {
			templates[name] = t
			log.Println("Loaded template", i, name)
		}
	}
}

func welcomeHandler(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	_ = templates["welcome"].Execute(writer, struct{}{})
}

func listHandler(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	_ = templates["list"].Execute(writer, responses)
}

type formData struct {
	*Rsvp
	Errors []string
}

func formHandler(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method == http.MethodGet {
		_ = templates["form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		if err := request.ParseForm(); err != nil {
			log.Println(err)
			return
		}
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}
		errors := errors(responseData)
		if len(errors) > 0 {
			_ = templates["form"].Execute(writer, formData{
				Rsvp:   &responseData,
				Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				_ = templates["thanks"].Execute(writer, responseData.Name)
			} else {
				_ = templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func errors(responseData Rsvp) []string {
	errors := make([]string, 0, 3)
	if responseData.Name == "" {
		errors = append(errors, "Please enter your name")
	}
	if responseData.Email == "" {
		errors = append(errors, "Please enter your email")
	}
	if responseData.Phone == "" {
		errors = append(errors, "Please enter your phone number")
	}
	return errors
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
