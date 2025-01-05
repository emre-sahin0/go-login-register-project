package handlers

import (
	"html/template"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Token'ı cookie'den al
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenStr := cookie.Value

	// Token'ı doğrula
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Kullanıcı adını al
	username := (*claims)["username"].(string)

	// Şablona veri gönder
	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, data)
}
