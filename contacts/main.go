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
	// get
	router.HandleFunc("/contacts", contactHandler)
	// post
	router.HandleFunc("/contacts/new", newContactHandler)

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
		data = lo.Filter(data, func(contact Contact, i int) bool {
			return contact.Name == name
		})
	}

	if err := tmpl.ExecuteTemplate(w, "content", data); err != nil {
		log.Println(err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func newContactHandler(w http.ResponseWriter, request *http.Request) {
	// get query params
	name := request.FormValue("name")
	email := request.FormValue("email")
	// parse template
	tmpl := template.Must(template.ParseFiles("templates.html"))
	// get contact data
	data, err := loadContacts("contacts.json")
	if err != nil {
		log.Println("can't load contacts", err)
	}

	if name != "" && email != "" {
		c := Contact{Name: name, Email: email}
		data = append(data, c)
	}

	// write json file
	err = saveContacts("contacts.json", data)
	if err != nil {
		log.Println("can't save new contacts file", err)
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

func saveContacts(filename string, contacts []Contact) error {
	// Marshal the contacts slice to JSON
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}

	// Create or open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(data)
	return err
}
