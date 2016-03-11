package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/17media/heath-ledger/models"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

// createUser - Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w)
	}
	fmt.Fprint(w, body)
	fmt.Fprint(w, "User Created!\n")
}

// getUser - Get User by ID
func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Get User!\n")
}

// listUser - User listing
func listUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users := models.GetUsers()
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, response)
}

// updateUser - Update User by ID
func updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Update User!\n")
}

// deleteUser - Delete User by ID
func deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Delete User!\n")
}
