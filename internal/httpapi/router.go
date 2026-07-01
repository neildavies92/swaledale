package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/neildavies/swaledale/internal/store"
)

type Options struct {
	SecureCookies bool
}

func NewRouter(store store.Store, opts Options) http.Handler {
	api := &API{store: store, secureCookies: opts.SecureCookies}
	r := chi.NewRouter()
	r.Use(cors)
	r.Get("/healthz", api.health)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", api.register)
		r.Post("/auth/login", api.login)

		r.Group(func(r chi.Router) {
			r.Use(api.requireAuth)
			r.Get("/auth/me", api.me)
			r.Post("/auth/logout", api.logout)
			r.Get("/summary", api.summary)
			r.Get("/members", api.members)
			r.Get("/members/{memberID}/budget", api.memberBudget)
			r.Put("/members/{memberID}/budget-items/{itemID}", api.updateBudgetItem)
			r.Get("/joint-account", api.jointAccount)
			r.Put("/joint-account/items/{itemID}", api.updateJointAccountItem)
			r.Get("/goals", api.goals)
			r.Put("/goals/{goalID}", api.updateGoal)
		})
	})

	return r
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
