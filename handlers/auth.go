package handlers

import (
	"context"
	"html/template"
	"net/http"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection
var store = sessions.NewCookieStore([]byte("super-secret-key")) // Güvenli bir key kullanın
var jwtKey = []byte("super-secret-key")                         // Güvenli bir anahtar kullanın

func GenerateJWT(username string) (string, error) {
	// Token yükünü oluştur
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24 saat geçerli
	}

	// Token'ı oluştur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Şifreyi hash'leme
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Şifre hashleme sırasında bir hata oluştu!", http.StatusInternalServerError)
			return
		}

		// MongoDB'ye kullanıcı ekleme
		_, err = UserCollection.InsertOne(context.TODO(), bson.M{
			"username": username,
			"password": string(hashedPassword),
		})
		if err != nil {
			http.Error(w, "Kayıt sırasında bir hata oluştu!", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Kullanıcı bilgilerini kontrol et
		var result bson.M
		err := UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
		if err != nil {
			http.Error(w, "Geçersiz kullanıcı adı veya şifre!", http.StatusUnauthorized)
			return
		}

		// Şifre doğrulama
		storedPassword := result["password"].(string)
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Geçersiz kullanıcı adı veya şifre!", http.StatusUnauthorized)
			return
		}

		// JWT oluştur veya session başlat
		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Oturum (session) sonlandırma
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1 // Oturumu geçersiz yap
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
