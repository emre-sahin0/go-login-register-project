package handlers

import (
	"context"
	"html/template"
	"net/http"
	"regexp"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection
var store = sessions.NewCookieStore([]byte("super-secret-key"))
var jwtKey = []byte("super-secret-key")

func GenerateJWT(username string) (string, error) {
	// Token yükünü oluştur
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

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
		email := r.FormValue("email")
		password := r.FormValue("password")

		// E-posta doğrulama
		if !isValidEmail(email) {
			http.Error(w, "Invalid email format!", http.StatusBadRequest)
			return
		}

		// Şifre hashleme
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// Kullanıcıyı veritabanına kaydet
		_, err := UserCollection.InsertOne(context.TODO(), bson.M{
			"email":    email,
			"password": string(hashedPassword),
			"photos":   []string{},
		})

		if err != nil {
			http.Error(w, "Could not register user!", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func isValidEmail(email string) bool {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		// bilgilerini kontrol et
		var result bson.M
		err := UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)
		if err != nil {
			http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
			return
		}

		// Şifre doğrulama
		storedPassword := result["password"].(string)
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
			return
		}

		// Session başlat
		session, _ := store.Get(r, "session")
		session.Values["email"] = email
		session.Save(r, w)

		// Dashboard'a yönlendirme
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Oturum (session) sonlandırma
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
