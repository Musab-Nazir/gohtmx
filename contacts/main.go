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

	// Contacts
	router.HandleFunc("/contacts", contactHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func contactHandler(w http.ResponseWriter, request *http.Request) {
	// get query params
	query := request.URL.Query()
	name := query.Get("q")
	// parse template
	tmpl := template.Must(template.ParseFiles("public/index.html"))
	// get contact data
	data, err := loadContacts("contacts.json")
	if name != "" {
		log.Printf("Got query param: %s", name)
		data = lo.Filter(data, func(contact Contact, i int) bool {
			return contact.FirstName == name
		})
	}
	if err != nil {
		log.Println(err)
	}
	// return view
	tmpl.Execute(w, data)
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
