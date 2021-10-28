package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/dittrichlucas/nykshortener/handler"
	"github.com/dittrichlucas/nykshortener/store"
)

const port = ":3000"

type TodoPageData struct {
	PageTitle string
	Name      string
}

func main() {
	store.Redis()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	tmpl := template.Must(template.ParseFiles("./static/generated.html"))
	http.HandleFunc("/generated", func(w http.ResponseWriter, r *http.Request) {
		test := handler.CreateShortURL(w, r)

		data := TodoPageData{
			PageTitle: "nykshortener",
			Name:      test.ShortURL,
		}
		tmpl.Execute(w, data)
	})

	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
