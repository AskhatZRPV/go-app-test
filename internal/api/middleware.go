package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func Log(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logrus.Info(r.Method, ": request sent to ", r.URL, " from ", r.RemoteAddr)
		n(w, r, ps)
	}
}
