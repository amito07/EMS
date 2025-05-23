package student

import "net/http"

func New() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the EMS!"))
	}
}

func Create() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create a new student"))
	}
}