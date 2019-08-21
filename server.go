package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session
var db *mgo.Database
var AppConfig Configuration

// Starting connection wiht db
func connectToDatabase() {

	// connect to the database
	dialInfo, err := mgo.ParseURL(AppConfig.DbConnection)
	//fmt.Printf("%+v\n", dialInfo)

	if err != nil {

		panic(err)
	}

	fmt.Println("Trying to open DB")

	//Below part is similar to above.
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	//session, err := mgo.Dial(AppConfig.DbConnection)
	if err != nil {
		panic(err)
	}

	session, err1 := mgo.DialWithInfo(dialInfo)
	//defer session.Close()

	if err1 != nil {
		fmt.Println("something bad happened")
		panic(err1)
	}
	fmt.Println("Opened DB")
	session.SetMode(mgo.Monotonic, true)
	// set the global
	mgoSession = session
	// db = session.DB(DATABASE)
}

// MyServer .....
type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func respondWithJsonn(w http.ResponseWriter, code int, reslt []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(reslt)
}

func registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	var data StructKeys
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Printf("There was an error decoding the json. err = %s", err)
		responseToSender := make(map[string]string)
		responseToSender["Status"] = "Error"
		responseToSender["message"] = "There was an error decoding the json"
		js, _ := json.Marshal(responseToSender)
		respondWithJsonn(w, http.StatusOK, js)
	} else {
		obj, yubierr := findYubikey(data.Yubikey, *mgoSession)
		fmt.Println(obj)
		if yubierr != nil {
			if yubierr.Error() == "not found" {
				obj1, mobilerr := findMobilkey(data.MobileKey, *mgoSession)
				fmt.Println(obj1)
				if mobilerr != nil {
					if mobilerr.Error() == "not found" {
						insErr := writeToUsersCollection(data, *mgoSession)
						if insErr != nil {
							fmt.Printf("Error in inserting data : %s", insErr)
							responseToSender := make(map[string]string)
							responseToSender["Status"] = "Error"
							responseToSender["message"] = "There was an error in inserting data"
							js, _ := json.Marshal(responseToSender)
							respondWithJsonn(w, http.StatusOK, js)
						} else {
							responseToSender := make(map[string]string)
							responseToSender["Status"] = "Success"
							responseToSender["message"] = "Keys registered successfully"
							js, _ := json.Marshal(responseToSender)
							respondWithJsonn(w, http.StatusOK, js)
						}
					} else {
						fmt.Printf("Error in fetching data : %s", mobilerr)
						responseToSender := make(map[string]string)
						responseToSender["Status"] = "Error"
						responseToSender["message"] = "Error in fetching mobile key"
						js, _ := json.Marshal(responseToSender)
						respondWithJsonn(w, http.StatusOK, js)
					}
				} else {
					fmt.Printf("mobileKey already registered : %s", mobilerr)
					responseToSender := make(map[string]string)
					responseToSender["Status"] = "Error"
					responseToSender["message"] = "mobileKey already registered"
					js, _ := json.Marshal(responseToSender)
					respondWithJsonn(w, http.StatusOK, js)
				}
			} else {
				fmt.Printf("Error in fetching data : %s", yubierr)
				responseToSender := make(map[string]string)
				responseToSender["Status"] = "Error"
				responseToSender["message"] = "Error in fetching yubi key"
				js, _ := json.Marshal(responseToSender)
				respondWithJsonn(w, http.StatusOK, js)
			}
		} else {
			fmt.Printf("Yubikey already registered : %s", yubierr)
			responseToSender := make(map[string]string)
			responseToSender["Status"] = "Error"
			responseToSender["message"] = "Yubikey already registered"
			js, _ := json.Marshal(responseToSender)
			respondWithJsonn(w, http.StatusOK, js)
		}
	}
}

func authentication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	var data StructKeys
	err := json.NewDecoder(r.Body).Decode(&data)
	fmt.Println(data)
	if err != nil {
		fmt.Printf("There was an error decoding the json. err = %s", err)
		responseToSender := make(map[string]string)
		responseToSender["Status"] = "Error"
		responseToSender["message"] = "There was an error decoding the json"
		js, _ := json.Marshal(responseToSender)
		respondWithJsonn(w, http.StatusOK, js)
	} else {
		User, getErr := getUserByKeys(data, *mgoSession)
		if getErr != nil {
			fmt.Println("Error in fetching data : ", getErr)
			responseToSender := make(map[string]string)
			responseToSender["Status"] = "Error"
			responseToSender["message"] = "Invalid Key"
			js, _ := json.Marshal(responseToSender)
			respondWithJsonn(w, http.StatusOK, js)
		} else {
			fmt.Println("User With keys is : ", User)
			responseToSender := make(map[string]string)
			responseToSender["Status"] = "Success"
			responseToSender["message"] = "Valid Keys"
			js, _ := json.Marshal(responseToSender)
			respondWithJsonn(w, http.StatusOK, js)
		}
	}
}

func main() {
	fmt.Println("Reading config file")

	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Error occured while reading config")
		return
	}
	json.Unmarshal(raw, &AppConfig)

	fmt.Printf("Running with configuration = \r\n%+v\n", AppConfig)

	r := mux.NewRouter()
	connectToDatabase()
	r.HandleFunc("/dacs/u2fRegistration", registration)
	r.HandleFunc("/dacs/u2fAuthentication", authentication)
	// Configure websocket route
	//http.HandleFunc("/ws", handleConnections)
	http.Handle("/", &MyServer{r})
	//http.ListenAndServe(":2000", nil)

	http.ListenAndServe(AppConfig.Port, nil)
}
