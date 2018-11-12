package keys

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/keys"
	"net/http"
	"strings"
)

///////////////////////////
// REST

// get key REST handler
func GetKeyRequestHandler(indent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		info, err := keys.GetKey(name)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("Key %s not found", name)) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}

		keyOutput, err := keys.Bech32KeyOutput(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		keys.PostProcessResponse(w, cdc, keyOutput, indent)
	}
}