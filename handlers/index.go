package handlers

import (
	"context"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)

	var user struct {
		Username string
		Photos   []string
	}

	if ok {

		err := UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
		if err != nil {
			http.Error(w, "Kullanıcı bilgileri alınamadı!", http.StatusInternalServerError)
			return
		}
	}

	data := struct {
		IsLoggedIn bool
		Username   string
		Photos     []string
	}{
		IsLoggedIn: ok,
		Username:   username,
		Photos:     user.Photos,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}
