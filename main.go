package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	
	b64 "encoding/base64"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

)

type BannedGroup struct {
	GroupID int `json:"groupId"`
}

func main() {

	sdk, _ := b64.StdEncoding.DecodeString(os.Getenv("FIREBASE_SDK"))
	log.Println(sdk, _)
	opt := option.WithCredentialsJSON(sdk)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}

	client, err := app.DatabaseWithURL(context.Background(), getEnv("FIREBASE_LINK"))
	if err != nil {
		log.Fatalf("Error initializing database client: %v\n", err)
	}

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/index.html")

		if err != nil {
			log.Printf("Error parsing HTML template: %v", err)
			return
		}

		tmpl.Execute(w, nil)
	}

	blacklistHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./templates/blacklist.html")
		if err != nil {
			log.Printf("Error parsing HTML template: %v", err)
			return
		}

		ref := client.NewRef("/blacklist/groups")

		var data map[string]BannedGroup
		if err := ref.Get(context.Background(), &data); err != nil {
			log.Fatalf("Error getting value: %v\n", err)
		}

		tmpl.Execute(w, data)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/blacklist", blacklistHandler)

	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getEnv(key string) string {
	return os.Getenv(key)
}
