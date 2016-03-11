package controllers

import (
        "github.com/17media/heath-ledger/models"
        "github.com/julienschmidt/httprouter"
)

// createUser - Create a new user
func createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        user := &models.User{
                Email: "sam@17.media"
        }
        
        json := map[string][string]
        
        json["user"] getuser
        models.newUser()
        fmt.Fprint(w, "User Created!\n")
}

// getUser - Get User by ID
func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprint(w, "Get User!\n")
}

// listUser - User listing
func listUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprint(w)
        fmt.Fprint(w, "User List!\n")
}

// updateUser - Update User by ID
func updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprint(w, "Update User!\n")
}

// deleteUser - Delete User by ID
func deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprint(w, "Delete User!\n")
}

func getuserList() {
        // HACK:
        
        user := getUser()
        follower := getFolower(user)
        block := getBolck(user)
        
        response := map[string][string]
        
        
}
