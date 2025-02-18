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
			"static/"+"layout"+".html",
			"static/"+name+".html",
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
	request *http.Request,
) {

	if err := templates["welcome"].Execute(writer, struct{}{}); err != nil {
		log.Println(err)
	}
}

func listHandler(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if err := templates["list"].Execute(writer, responses); err != nil {
		log.Println(err)
	}
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
		if err := templates["form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		}); err != nil {
			log.Println(err)
		}
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
		responses = append(responses, &responseData)
		if responseData.WillAttend {
			if err := templates["thanks"].Execute(writer, responseData.Name); err != nil {
				log.Println(err)
			}
		} else {
			if err := templates["sorry"].Execute(writer, responseData.Name); err != nil {
				log.Println(err)
			}
		}
	}
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
