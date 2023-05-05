package scraping_zenn

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func ErrorResponse(w http.ResponseWriter, code int, message string) {
	Response(w, code, map[string]string{"error": message})
}
func Response(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Fatal().Msg(err.Error())
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err = w.Write(res); err != nil {
		log.Fatal().Msg(err.Error())
		os.Exit(1)
	}
}
