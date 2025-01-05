package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Oturumdan kullanıcı bilgilerini al
		session, _ := store.Get(r, "session")
		email, ok := session.Values["email"].(string)

		if !ok || email == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Dosya işlemleri
		file, header, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "Dosya yüklenemedi!", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Dosyayı kaydet
		savePath := fmt.Sprintf("static/uploads/%s", header.Filename)
		out, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Dosya kaydedilemedi!", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Dosya yazılamadı!", http.StatusInternalServerError)
			return
		}

		// Veritabanına fotoğrafı ekle
		_, err = UserCollection.UpdateOne(context.TODO(),
			bson.M{"email": email},
			bson.M{"$push": bson.M{"photos": header.Filename}},
		)
		if err != nil {
			http.Error(w, "Fotoğraf veritabanına kaydedilemedi!", http.StatusInternalServerError)
			return
		}

		// Dashboard'a yönlendirme
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
