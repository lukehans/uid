package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    "uid/sample/uid"
)

const errID = 0

type UniqueID struct {
    Value uint64 `json:"value"`
}

var (
    idRequestChan = make(chan chan uint64)
    infoLog       = log.New(os.Stdout, "INFO: ", log.LstdFlags)
    errLog        = log.New(os.Stdout, "ERROR: ", log.LstdFlags)
)

// main starts a service to reserve unique IDs. Unique IDs are uint64 values and are
// requested by sending a POST request to the service, indicating a request to 
// reserve a unique ID from the pool of possible IDs. The ID is returned to the
// requester via a JSON object in an HTTP response message payload. 
//
// The service can handle multiple concurrent users, and preserves ID uniqueness
// upon restart. The service can be stopped with an operating system kill command,
// and the service attempts to gracefully shut down, immediately preventing new
// connections but allowing existing service requests 10 more seconds to complete.
func main() {
    uid.InitPersistenceLayer()

    go listenForIDRequests()

    infoLog.Println("Server is starting...")
    srv := startSrv()

    done := make(chan bool)
    exit := make(chan os.Signal, 1)
    signal.Notify(exit, os.Interrupt)
    go func() {
        <-exit
        shutDownGracefully(srv)
        close(done)
    }()

    infoLog.Println("Server is ready.")

    // Wait until channel 'done' is closed before ending.
    <-done
    infoLog.Println("Server stopped. Program exiting.")
}

// startSrv starts an HTTP server listening on port 8080 and registers a handler
// to listen for requests hitting the URL path '/uid'. 
func startSrv() *http.Server {
    srv := &http.Server{Addr: ":8080"}
    http.HandleFunc("/sample/uid", serveUniqueID)

    go func() {
        err := srv.ListenAndServe()
        if err != nil {
            errLog.Printf("HTTP Server interrupted: %s", err)
        }
    }()

    return srv
}

// shutDownGracefully attempts to gracefully shut down the given server. Existing
// connections are allowed to work for 10 seconds until the connections are forced
// to close.
func shutDownGracefully(srv *http.Server) {
    infoLog.Println("Server is shutting down...")
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        errLog.Fatalf("Graceful shut down failed - error: %v\n", err)
    }
}

// serveUniqueID returns a unique ID to the requester through the ResponseWriter. It
// first validates that the HTTP request is a POST request, then opens an ID 
// receiver channel (idReceiverChan). The ID receiver channel is sent through the ID
// request channel (a channel is sent in another channel). A goroutine is waiting to
// accept from the ID request channel and will send back a unique ID through the ID
// receiver channel.
//
// The unique ID is returned in a JSON Object in the HTTP response payload. The JSON
// object is in the following format:
//     {"value":"<uint64 unique ID>"}
//
func serveUniqueID(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        errLog.Println("Received non-POST request.")
        http.Error(w, "POST requests only, please.", http.StatusMethodNotAllowed)
        return
    }
   
    // Create receiving channel for a unique ID
    var idReceiverChan = make(chan uint64)
    
    // Send ID receiving channel to the ID request channel
    idRequestChan <- idReceiverChan

    // Read an ID from the receiving channel.
    id := <- idReceiverChan

    // If the received ID is 0, there was an underlying error.
    if id == errID {
        w.WriteHeader(http.StatusInternalServerError)
        errLog.Printf("Error returned to requester. Something went wrong with ID generation.")
    } else {
        // Write the unique ID to the ResponseWriter
        uniqueID := UniqueID{Value: id}
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(uniqueID)
        infoLog.Printf("ID sent to requester: %d", uniqueID)
    }
}

// listenForIDRequests listens for ID requests from the idRequestChan. ID requests
// are themselves represented as channels (idReturnChan), providing the mechanism to
// return a unique ID to the caller.
//
// In other words, the idRequestChan is a channel of channels. The inner channel is
// a response channel where IDs (uint64 integers) are returned to the requester.
// 
// When a channel is received from the idRequestChan (indicating a new ID request),
// a unique ID is retrieved using the uid package. Next, the ID is sent back to the
// requester through the idReturnChan. Finally, the method loops back to continue 
// listening for more ID requests.
func listenForIDRequests() {
    for idReturnChan := range idRequestChan {
        // Get a unique ID
        id, err := uid.GetNextID()
        if err != nil {
            errLog.Printf("Error getting the next ID: %s", err)
            idReturnChan <- errID
        }
        // Send a unique ID to the ID return channel
        idReturnChan <- id
    }
}
