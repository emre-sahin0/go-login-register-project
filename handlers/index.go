package handlers

import (
	"context"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcı oturum bilgilerini kontrol et
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)

	var user struct {
		Username string
		Photos   []string
	}

	if ok {
		// Kullanıcı fotoğraflarını veritabanından getir
		err := UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
		if err != nil {
			http.Error(w, "Kullanıcı bilgileri alınamadı!", http.StatusInternalServerError)
			return
		}
	}

	// Şablona veri gönder
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
