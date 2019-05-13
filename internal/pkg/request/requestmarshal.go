package request

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func ParseRequest(r *http.Request, t interface{}) (error error) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		return errors.Wrap(err, "error when decode http request")
	}
	return
}
