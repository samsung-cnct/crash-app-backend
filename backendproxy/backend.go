package backendproxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/gorilla/mux"
)

// remove esIndexOld soon
const esIndexOld string = "/k2crashreporter/k2crashes"
const esIndex string = "/krakencrashreporter/krakencrashes"

// InputValidation is the interface for validation
type InputValidation interface {
	Validate(path string) error
}

var target = "http://google.com"
var remote *url.URL
var maxCharsKrakenLog = 2000000
var maxCharsFailedTask = 400
var rateLimit = 10

// CrashAppMessage is the expected json from Crash App
type CrashAppMessage struct {
	Date       string `json:"date,omitempty"`
	KrakenLog  string `json:"k2_log"`
	FailedTask string `json:"failed_task"`
}

// Validate method for CrashAppMessage
func (t CrashAppMessage) Validate(path string) (bool, error) {
	switch path {
	case esIndex:
		if t.Date == "" {
			return false, errors.New("Invalid CrashApp data found, Expected date")
		}
		// validate Date
		reDate := regexp.MustCompile(`^2[0-9]{3}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{3}Z`)
		if !reDate.MatchString(t.Date) {
			return false, errors.New("Invalid CrashApp data found, Invalid date format")
		}
	case esIndexOld:
		if t.Date != "" {
			return false, errors.New("Invalid CrashApp data found, Unexpected date")
		}
	default:
		return false, errors.New("Invalid CrashApp data found, Invalid elasticsearch index")
	}

	if tLen := len(t.FailedTask); tLen > maxCharsFailedTask {
		return false, fmt.Errorf("Invalid CrashApp data found, FailedTask is too large %d > max %d", tLen, maxCharsFailedTask)
	}
	reFailed := regexp.MustCompile("(.+?) : (.+?)")
	if !reFailed.MatchString(t.FailedTask) {
		return false, errors.New("Invalid CrashApp data found, Task not found in KrakenLog")
	}

	if logLen := len(t.FailedTask); logLen > maxCharsKrakenLog {
		return false, fmt.Errorf("Invalid CrashApp data found, KrakenLog is too large %d > max %d", logLen, maxCharsKrakenLog)
	}
	re := regexp.MustCompile(`TASK \[(.+?)\]`)
	matchedTask := re.MatchString(t.KrakenLog)
	re2 := regexp.MustCompile("PLAY")
	matchedPlay := re2.MatchString(t.KrakenLog)
	if !matchedPlay || !matchedTask {
		return false, errors.New("Invalid CrashApp data found, Invalid KrakenLog")
	}

	return true, nil
}

func handlerHealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm Alive!"))
}

// HandlerCrashAppValidation to handle CrashAppMessage
func HandlerCrashAppValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in HandlerCrashAppValidation")
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	buf, _ := ioutil.ReadAll(r.Body)
	rdr := ioutil.NopCloser(bytes.NewBuffer(buf))

	var msg CrashAppMessage

	errd := json.NewDecoder(rdr).Decode(&msg)
	if errd != nil {
		http.Error(w, "Invalid json data - can't decode", http.StatusBadRequest)
		return
	}
	validated, verr := msg.Validate(r.URL.Path)
	if !validated {
		fmt.Println(verr)
		http.Error(w, verr.Error(), http.StatusBadRequest)
		return
	}

	// proxy original request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	proxy.ServeHTTP(w, r)
}

// InitConfig initialize variables based on arguments
func InitConfig(t string, klogmax int, ktaskmax int, ratelimit int) {
	if t != "" {
		target = t
	}

	tempRemote, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}
	remote = tempRemote
	maxCharsKrakenLog = klogmax
	maxCharsFailedTask = ktaskmax
	rateLimit = ratelimit
}

func handleRequests() {

	router := mux.NewRouter()
	router.HandleFunc(esIndex, HandlerCrashAppValidation).Methods("POST")
	router.HandleFunc(esIndexOld, HandlerCrashAppValidation).Methods("POST")
	router.HandleFunc("/healthcheck", handlerHealthCheck).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}

// Server starts the web server
func Server(t string, klogmax int, ktaskmax int, ratelimit int) {
	InitConfig(t, klogmax, ktaskmax, ratelimit)
	handleRequests()
}
