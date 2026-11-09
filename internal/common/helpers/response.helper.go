package helpers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//recover from panic
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println("panic", rec)
				WriteJson(w, http.StatusInternalServerError, map[string]any{
					"statusCode": http.StatusInternalServerError,
					"message":    "Internal Server Error",
				})
			}
		}()

		err := h(w, r)
		if err != nil {
			if apiErr, ok := err.(*ApiError); ok {
				WriteJson(w, apiErr.StatusCode, apiErr)
			} else {
				//WriteJson(w, http.StatusBadRequest, err.Error())
				errResp := map[string]any{
					"statusCode": http.StatusBadRequest,
					"message":    err.Error(),
				}
				WriteJson(w, http.StatusBadRequest, errResp)
			}
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
		}
		slog.Info("HTTP API success", "path", r.URL.Path)
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteJsonString(w http.ResponseWriter, status int, v string) error {
	fmt.Println("serve form cache")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(v)); err != nil {
		return err
	}

	return nil
}
