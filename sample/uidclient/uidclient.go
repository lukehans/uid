package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
	"net/http"
	"os"
)

const totalReq = 500

type UniqueID struct {
    Value uint64 `json:"value"`
}

var (
    infoLog = log.New(os.Stdout, "INFO: ", log.LstdFlags)
    errLog  = log.New(os.Stdout, "ERROR: ", log.LstdFlags)
    idMap   = make(map[uint64]bool)
    idChan  = make(chan uint64)
    done    = make(chan bool)
)

// main starts a client test to use the Unique ID service. A goroutine is started
// to listen to the idChan, register IDs as they are received from the server, and
// check for duplicates. main() also starts goroutines to send a number of HTTP 
// requests to the Unique ID service. The number of HTTP requests to be send is
// specified by 'totalReq'.
func main() {
	infoLog.Printf("Client started. Beginning to send %d HTTP requests.", totalReq)

	go checkIDs()
	for i := 0; i < totalReq; i++ {
		go sendRequest()
	}
	<-done
}

// sendRequest sends one HTTP request to the Unique ID service to request a unique
// ID. The HTTP response body is parsed into JSON, the unique ID value is extracted,
// and the unique ID is sent to the idChan to registed that particular ID as
// received. 
func sendRequest() {
	resp, err := http.Post("http://localhost:8080/sample/uid", "text/html", nil)

	if err != nil {
		errLog.Printf("Error in HTTP response.", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var id UniqueID
	err = json.Unmarshal([]byte(string(body)), &id)
	infoLog.Printf("Successful HTTP response. Parsed ID is %d", id.Value)
	idChan <- id.Value
}

// checkIDs listens to the idChan, where IDs will be sent once received from the
// ID service in an HTTP response. When an ID is received, the idMap is checked to
// see if the ID has already been put in the map (idMap[id] == true in this case).
// If a duplicate is detected, then panic. Otherwise, add the id as a map key with
// value 'true'. Once the idMap length is the same size as how many requests were
// sent, then signal program exit by closing the 'done' channel.
func checkIDs() {
	for id := range idChan {
		if idMap[id] == true {
			panic("Uh oh - Duplicate ID detected!")
		}
		idMap[id] = true

		if len(idMap) == totalReq {
			infoLog.Println("All responses are received successfully. No duplicates!")
			close(done)
		}
    }
}