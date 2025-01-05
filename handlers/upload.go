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
		username, ok := session.Values["username"].(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Yüklenen dosyayı al
		file, header, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "File upload error!", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Dosyayı kaydet
		savePath := fmt.Sprintf("static/uploads/%s", header.Filename)
		out, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "File save error!", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "File write error!", http.StatusInternalServerError)
			return
		}

		// Veritabanına fotoğrafı ekle
		_, err = UserCollection.UpdateOne(context.TODO(),
			bson.M{"username": username},
			bson.M{"$push": bson.M{"photos": header.Filename}},
		)
		if err != nil {
			http.Error(w, "Database update error!", http.StatusInternalServerError)
			return
		}

		// Başarıyla yönlendir
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
