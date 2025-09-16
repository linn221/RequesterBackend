package middlewares

import (
	"fmt"
	"net/http"

	"github.com/linn221/RequesterBackend/utils"
)

// Host: host portion(localhost)
// SecretPath: uri without the slash (start-session)
// RedirectUrl: absolute url to redirect
type SecretConfig struct {
	SecretFunc  func() string
	SecretPath  string
	RedirectUrl string
	Host        string
	// expiration time.Time // for later
}

func (cfg *SecretConfig) Middleware() func(h http.Handler) http.Handler {
	theSecret := cfg.SecretFunc()
	fmt.Printf("Magic auth link: %s/%s?secret=%s\n", cfg.Host, cfg.SecretPath, theSecret)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUrl := r.URL.Path
			if currentUrl == "/"+cfg.SecretPath {
				pretendSecret := r.URL.Query().Get("secret")
				if pretendSecret == theSecret {
					// authentication success
					utils.SetCookies(w, "secret", theSecret)
					http.Redirect(w, r, cfg.RedirectUrl, http.StatusTemporaryRedirect)
					return
				}
				http.Error(w, "please visit the magic link for auth", http.StatusUnauthorized)
				return
			}
			cookies, err := r.Cookie("secret")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "please visit the magic link for auth", http.StatusUnauthorized)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pretendSecret := cookies.Value
			if pretendSecret == theSecret {
				h.ServeHTTP(w, r)
				return
			}

			utils.RemoveCookies(w, "secret")
			http.Error(w, "please visit the magic link for auth", http.StatusUnauthorized)
		})
	}
}
