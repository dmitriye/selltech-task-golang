package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sdn/internal/app"
)

type Response struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
	Code   int    `json:"code"`
}

func NewUpdate(s *app.AppState, u *app.Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		if err := u.Run(); err != nil {
			// TODO: log error
			fmt.Printf("[error] %s\n", err)

			json.NewEncoder(w).Encode(Response{
				Result: false,
				Info:   http.StatusText(http.StatusInternalServerError),
				Code:   http.StatusInternalServerError,
			})
			return
		}

		json.NewEncoder(w).Encode(Response{
			Result: true,
			Info:   http.StatusText(http.StatusOK),
			Code:   http.StatusOK,
		})
	}
}
