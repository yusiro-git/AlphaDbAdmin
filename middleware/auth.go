package middleware

import (
	"net/http"
)

// AuthMiddleware ensures that only authorized users can access protected routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil || cookie.Value != "723hbo2uipfir3]g1h734807jr9-237fgy3bigni1fjp0h84gy3t" {
			http.Redirect(w, r, "/auth", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
