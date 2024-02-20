package state

import (
	"encoding/json"
	"net/http"
	"sdn/internal/app"
)

type Response struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

func NewState(s *app.AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		json.NewEncoder(w).Encode(Response{
			Result: s.IsOk(),
			Info:   s.GetName(),
		})
	}
}
