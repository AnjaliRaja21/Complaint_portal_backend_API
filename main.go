package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID         string      `json:"id"`
	SecretCode string      `json:"secretCode"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints,omitempty"`
}

type Complaint struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Rating   int    `json:"rating"`
	Resolved bool   `json:"resolved"`
}

var usersDB map[string]User = make(map[string]User)

func ReturnJsonResponse(res http.ResponseWriter, resMessage []byte) {
	res.Header().Set("content-type", "application/json")
	res.Write(resMessage)
}

func loginUser(secretCode string) (User, error) {

	for _, user := range usersDB {
		if user.SecretCode == secretCode {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("invalid secret code")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" :"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	var requestBody struct {
		SecretCode string `json:"secretCode"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" : "Error parsing the req body data",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	user, err := loginUser(requestBody.SecretCode)

	if err != nil {
		HandlerMessage := []byte(`{
			"success":false,
			"message":"Wrong secret code/user not found",
 }`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	userJSON, err := json.MarshalIndent(user, "", "\t")

	if err != nil {
		HandlerMessage := []byte(`{
   "success":false,
   "message":"Error parsing the user data",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	HandlerMessage := []byte(`{
		"success" : true,
		"message" : "user sign-in successfully",
	}`)
	ReturnJsonResponse(w, HandlerMessage)
	ReturnJsonResponse(w, userJSON)
	return
}

func generateUniqueID() string {

	rand.Seed(time.Now().UnixNano())

	uniqueID := strconv.Itoa((rand.Intn(10000)))

	return uniqueID
}

func generateUniqueSecretCode() string {

	rand.Seed(time.Now().UnixNano())

	secretCode := strconv.Itoa((rand.Intn(1000000)))

	return secretCode

}

func registerUser(name, email string) User {
	userID := generateUniqueID()
	secretCode := generateUniqueSecretCode()
	newUser := User{
		ID:         userID,
		SecretCode: secretCode,
		Name:       name,
		Email:      email,
		Complaints: []Complaint{},
	}
	usersDB[userID] = newUser
	return newUser
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" :"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	var newUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" : "Error parsing the request body data",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}
	_, ok := usersDB[newUser.Email]
	if ok {
		HandlerMessage := []byte(`{
	 		"success" : false,
	 		"message" : "User already exist",
	 	}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	user := registerUser(newUser.Name, newUser.Email)

	userJSON, err := json.MarshalIndent(user, "", "\t")

	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"Error parsing the  data",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	HandlerMessage := []byte(`{
		"success" : true,
		"message" : "New user sign-up",
	}`)
	ReturnJsonResponse(w, HandlerMessage)
	ReturnJsonResponse(w, userJSON)
	return
}

func submitComplaint(userID, title, summary string, rating int) (Complaint, error) {
	complaintID := generateUniqueID() // Implement a function to generate a unique complaint ID
	newComplaint := Complaint{
		ID:       complaintID,
		Title:    title,
		Summary:  summary,
		Rating:   rating,
		Resolved: false,
	}
	user, ok := usersDB[userID]
	if ok {
		// Handle user not found error
		return newComplaint, fmt.Errorf("user not found")
	}
	user.Complaints = append(user.Complaints, newComplaint)
	usersDB[userID] = user
	return newComplaint, nil
}

func submitComplaintHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" :"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	var newComplaint struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Rating  int    `json:"rating"`
	}
	err := json.NewDecoder(r.Body).Decode(&newComplaint)
	if err != nil {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" : "Error parsing the complain data",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	userID := "348"

	complaint, err := submitComplaint(userID, newComplaint.Title, newComplaint.Summary, newComplaint.Rating)
	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"User not found",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	complaintJSON, err := json.MarshalIndent(complaint, "", "\t")
	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"Error parsing the  data",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	HandlerMessage := []byte(`{
		"success" : true,
		"message" : "complain was successfully created",
	}`)
	ReturnJsonResponse(w, HandlerMessage)
	ReturnJsonResponse(w, complaintJSON)
}

// func getAllComplaintsForUser(userID string) ([]Complaint, error) {
// 	user, ok := usersDB[userID]
// 	if !ok {
// 		return nil, fmt.Errorf("user not found")
// 	}
// 	return user.Complaints, nil
// }

// func getAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != "GET" {
// 		HandlerMessage := []byte(`{
// 			"success":false,
//             "message":"check your HTTP method : Invalid HTTP method executed",
// 		}`)
// 		ReturnJsonResponse(w, HandlerMessage)
// 		return
// 	}

// 	userID := "6806"

// 	complaints, err := getAllComplaintsForUser(userID)

// 	if err != nil {
// 		HandlerMessage := []byte(`{
// 	       "success":false,
// 	       "message":"User not found",
// }`)
// 		ReturnJsonResponse(w, HandlerMessage)
// 		return
// 	}

// 	complaintsJSON, err := json.MarshalIndent(complaints, "", "\t")

// 	if err != nil {
// 		HandlerMessage := []byte(`{
// 	       "success":false,
// 	       "message":"Error parsing the user data",
// }`)
// 		ReturnJsonResponse(w, HandlerMessage)
// 		return
// 	}

// 	ReturnJsonResponse(w, complaintsJSON)
// }

func getAllComplaintsForAdmin() []Complaint {
	var allComplaints []Complaint
	for _, user := range usersDB {
		allComplaints = append(allComplaints, user.Complaints...)
	}
	return allComplaints
}

func getAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		HandlerMessage := []byte(`{
			"success":false,
            "message":"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	complaints := getAllComplaintsForAdmin()

	complaintsJSON, err := json.MarshalIndent(complaints, "", "\t")
	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"Error parsing the user data",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	ReturnJsonResponse(w, complaintsJSON)
}

func viewComplaint(userID, complaintID string) (Complaint, error) {
	user, ok := usersDB[userID]
	if !ok {
		return Complaint{}, fmt.Errorf("user not found")
	}
	for _, complaint := range user.Complaints {
		if complaint.ID == complaintID {
			return complaint, nil
		}
	}
	return Complaint{}, fmt.Errorf("complaint not found")
}

func viewComplaintHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		HandlerMessage := []byte(`{
			"success":false,
            "message":"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	userID := "user123"

	complaintID := r.URL.Query().Get("complaintID")

	complaint, err := viewComplaint(userID, complaintID)
	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"Not Found",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	complaintJSON, err := json.MarshalIndent(complaint, "", "\t")
	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"Error parsing the user data",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	ReturnJsonResponse(w, complaintJSON)
}

func resolveComplaint(complaintID string) error {
	for _, user := range usersDB {
		for idx, complaint := range user.Complaints {
			if complaint.ID == complaintID {
				user.Complaints[idx].Resolved = true
				usersDB[user.ID] = user
				return nil
			}
		}
	}
	return fmt.Errorf("complaint not found")
}

func resolveComplaintHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "PUT" {
		HandlerMessage := []byte(`{
			"success":false,
            "message":"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	complaintID := r.URL.Query().Get("complaintID")

	err := resolveComplaint(complaintID)

	if err != nil {
		HandlerMessage := []byte(`{
	       "success":false,
	       "message":"Complaint Not Found",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	HandlerMessage := []byte(`{
		"success":false,
		"message":"Resolved Succesfully",
}`)
	ReturnJsonResponse(w, HandlerMessage)
}

func main() {

	log.Println("Complaint API")

	// http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/register", registerHandler)
	// http.HandleFunc("/submitComplaint", submitComplaintHandler)
	// http.HandleFunc("/getAllComplaintsForUser", getAllComplaintsForUserHandler)
	// http.HandleFunc("/getAllComplaintsForAdmin", getAllComplaintsForAdminHandler)
	// http.HandleFunc("/viewComplaint", viewComplaintHandler)
	// http.HandleFunc("/resolveComplaint", resolveComplaintHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
