package respond

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		Error(w, errors.Wrap(err, "json marshal failed"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}
func Error(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
