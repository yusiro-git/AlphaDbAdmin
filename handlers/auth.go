package handlers

import (
	"net/http"
)

const authKey = "0976t3fjr98rfgyujok[394tfnip2jh806g9p0[21]]"

// AuthHandler handles the authentication page and key validation
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Serve the authentication page
		http.ServeFile(w, r, "templates/auth.html")
		return
	}

	// Handle POST request (key submission)
	key := r.FormValue("key")
	if key == authKey {
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
