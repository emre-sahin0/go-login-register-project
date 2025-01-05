package handlers

import (
	"html/template"
	"net/http"

	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)

	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user struct {
		Username string
		Photos   []string
	}

	err := UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "Kullanıcı bilgileri alınamadı!", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Photos   []string
	}{
		Username: user.Username,
		Photos:   user.Photos,
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, data)
}
