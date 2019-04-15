package web

import (
	"github.com/annakozyreva1/person_registrator/log"
	"github.com/annakozyreva1/person_registrator/person"
	"net/http"
)

var (
	logger = log.Logger
)

func registrate(reg *person.Registrator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		firstName := r.FormValue("firstname")
		if firstName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("firstname parameter is absent"))
			return
		}
		lastName := r.FormValue("lastname")
		if lastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("lastname parameter is absent"))
			return
		}
		isAdded := reg.Add(firstName, lastName)
		if isAdded {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Run(address string, reg *person.Registrator) {
	http.HandleFunc("/person", registrate(reg))
	if err := http.ListenAndServe(address, nil); err != nil {
		logger.Fatalf("failed web server: %v", err.Error())
	}
}
