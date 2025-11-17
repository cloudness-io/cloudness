package variable

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/variable"
)

func HandleList(varCtrl *variable.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderVariablePage(w, r, varCtrl)
	}
}
