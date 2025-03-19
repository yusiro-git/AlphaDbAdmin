package handlers

import (
	"net/http"
)

// func AuthHandler(extraValue string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Use the extraValue inside the handler
// 		fmt.Fprintf(w, "AuthHandler called with extra value: %s\n", extraValue)
// 	}
// }

func AuthHandler(magicKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Serve the authentication page
			http.ServeFile(w, r, "templates/auth.html")
			return
		}

		// Handle POST request (key submission)
		key := r.FormValue("key")
		if key == magicKey {
			// Set auth cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: "723hbo2uipfir3]g1h734807jr9-237fgy3bigni1fjp0h84gy3t",
			})
			// Redirect to homepage
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			// Redirect back to auth page if key is incorrect
			http.Redirect(w, r, "/auth", http.StatusFound)
		}
	}
}
