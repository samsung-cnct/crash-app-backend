package backendproxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestReverseProxy(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()

	target, err := url.Parse(backendServer.URL)
	if err != nil {
		log.Fatal(err)
	}
	frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(target))
	defer frontendProxy.Close()

	resp, err := http.Get(frontendProxy.URL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}

func TestValidMessageOldIndex(t *testing.T) {
	jsonmsg := CrashAppMessage{KrakenLog: "This is the log data PLAY TASK [blahblah]", FailedTask: "roles/kraken.readiness : Get all nodes"}
	b, err := json.Marshal(jsonmsg)
	if err != nil {
		fmt.Println("error:", err)
	}

	var msg CrashAppMessage
	erru := json.Unmarshal(b, &msg)
	if erru != nil {
		t.Error("error unmarshall", erru)
	}

	validated, err := msg.Validate(esIndexOld)
	if err != nil {
		fmt.Println(err)
	}
	if !validated {
		t.Error("should have been validated, got ", validated)
	}
}

func TestValidMessageNewIndex(t *testing.T) {
	jsonmsg := CrashAppMessage{Date: "2017-09-26T15:56:49.012Z", KrakenLog: "This is the log data PLAY TASK [blahblah]", FailedTask: "roles/kraken.readiness : Get all nodes"}
	b, err := json.Marshal(jsonmsg)
	if err != nil {
		fmt.Println("error:", err)
	}

	var msg CrashAppMessage
	erru := json.Unmarshal(b, &msg)
	if erru != nil {
		t.Error("error unmarshall", erru)
	}

	validated, err := msg.Validate(esIndex)
	if err != nil {
		fmt.Println(err)
	}
	if !validated {
		t.Error("should have been validated, got ", validated)
	}
}
func TestValidDate(t *testing.T) {
	jsonmsg := CrashAppMessage{Date: "2017-09-26T15:56:49.012Z", KrakenLog: "This is the log data PLAY TASK [blahblah]", FailedTask: "roles/kraken.readiness : Get all nodes"}
	b, err := json.Marshal(jsonmsg)
	if err != nil {
		fmt.Println("error:", err)
	}

	var msg CrashAppMessage
	erru := json.Unmarshal(b, &msg)
	if erru != nil {
		t.Error("error unmarshall", erru)
	}

	validated, err := msg.Validate(esIndex)
	if err != nil {
		fmt.Println(err)
	}
	if !validated {
		t.Error("should have been true, got ", validated)
	}
}
func TestInvalidDate(t *testing.T) {
	jsonmsg := CrashAppMessage{Date: "2017-09-26T15:56:49Z", KrakenLog: "This is the log data PLAY TASK [blahblah]", FailedTask: "roles/kraken.readiness : Get all nodes"}
	b, err := json.Marshal(jsonmsg)
	if err != nil {
		fmt.Println("error:", err)
	}

	var msg CrashAppMessage
	erru := json.Unmarshal(b, &msg)
	if erru != nil {
		t.Error("error unmarshall", erru)
	}

	validated, err := msg.Validate(esIndex)
	if err != nil {
		fmt.Println(err)
	}
	if validated {
		t.Error("should have been false, got ", validated)
	}
}

func TestInvalidKrakenLog(t *testing.T) {
	jsonmsg := CrashAppMessage{Date: "2017-09-26T15:56:49.022Z", KrakenLog: "This is the log data [blahblah]", FailedTask: "roles/kraken.readiness : Get all nodes"}
	b, err := json.Marshal(jsonmsg)
	if err != nil {
		fmt.Println("error:", err)
	}

	var msg CrashAppMessage
	erru := json.Unmarshal(b, &msg)
	if erru != nil {
		t.Error("error unmarshall", erru)
	}

	validated, err := msg.Validate(esIndex)
	if err != nil {
		fmt.Println(err)
	}
	if validated {
		t.Error("should have been false, got ", validated)
	}
}

func TestInvalidKrakenLogSize(t *testing.T) {
	jsonmsg := CrashAppMessage{Date: "2017-09-26T15:56:49.022Z", KrakenLog: "This is the log data [blahblah]", FailedTask: "roles/kraken.readiness : Get all nodes"}
	b, err := json.Marshal(jsonmsg)
	if err != nil {
		fmt.Println("error:", err)
	}

	var msg CrashAppMessage
	erru := json.Unmarshal(b, &msg)
	if erru != nil {
		t.Error("error unmarshall", erru)
	}

	maxCharsKrakenLog = 2
	validated, err := msg.Validate(esIndex)
	if err != nil {
		fmt.Println(err)
	}
	if validated {
		t.Error("should have been false, got ", validated)
	}
}
