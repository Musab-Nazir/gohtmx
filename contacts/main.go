package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/samber/lo"
)

func main() {
	runServer()
}

func runServer() {
	router := http.NewServeMux()
	// index
	router.HandleFunc("/", indexHandler)
	// Contacts
	router.HandleFunc("/contacts", contactHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func indexHandler(w http.ResponseWriter, request *http.Request) {
	// parse template
	tmpl := template.Must(template.ParseFiles("templates.html"))

	if err := tmpl.ExecuteTemplate(w, "index", nil); err != nil {
		log.Println(err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func contactHandler(w http.ResponseWriter, request *http.Request) {
	// get query params
	query := request.URL.Query()
	name := query.Get("name")
	// parse template
	tmpl := template.Must(template.ParseFiles("templates.html"))
	// get contact data
	data, err := loadContacts("contacts.json")
	if err != nil {
		log.Println("can't load contacts", err)
	}

	if name != "" {
		log.Printf("Got query param: %s", name)
		data = lo.Filter(data, func(contact Contact, i int) bool {
			return contact.Name == name
		})
	}

	if err := tmpl.ExecuteTemplate(w, "content", data); err != nil {
		log.Println(err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func loadContacts(filename string) ([]Contact, error) {
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	// Unmarshal the JSON data into a slice of Contact structs
	var contacts []Contact
	if err := json.Unmarshal(byteValue, &contacts); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	return contacts, nil
}
