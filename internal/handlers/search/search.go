package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sdn/internal/app"
	"strings"
)

func NewSearch(repo *app.EntryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		n := r.URL.Query().Get("name")
		f := app.T_WEAK

		if strings.EqualFold(r.URL.Query().Get("type"), "strong") {
			f = app.T_STRONG
		}

		entries, err := repo.Search(n, f)
		if err != nil {
			// TODO: log error
			fmt.Printf("[error] %s\n", err)

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(entries)
	}
}
