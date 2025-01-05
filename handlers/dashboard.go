package handlers

import (
	"context"
	"html/template"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// e-posta bilgi
	session, _ := store.Get(r, "session")
	email, ok := session.Values["email"].(string)

	if !ok || email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//  Mongoden al
	var user struct {
		Email  string
		Photos []string
	}
	err := UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "Kullanıcı bilgileri alınamadı!", http.StatusInternalServerError)
		return
	}

	//  "@" işaretine kadar olan
	idx := strings.Index(email, "@")
	username := email
	if idx > 0 {
		username = email[:idx]
	}

	// Dashboard'a gönderilecek veri
	data := struct {
		Username string
		Photos   []string
	}{
		Username: username,
		Photos:   user.Photos,
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, data)
}
