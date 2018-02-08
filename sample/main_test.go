package main

import (
	"io/ioutil"
    "math"
    "os"
	"strconv"
	"testing"
	"uid/sample/uid"
)

const persistenceFile = "last-used-id.txt"

// Test listenForIDRequests:
// Test that listenForIDRequests() is listening on the idRequestChan and
// it responds with a positive uint64 value through the idReceiverChan.
func TestListenForIDRequests(t *testing.T) {
	go listenForIDRequests()

    idReceiverChan := make(chan uint64)
    idRequestChan <- idReceiverChan
    id := <- idReceiverChan

    uint64ID := uint64(id)

	if uint64ID < 1 {
		t.Errorf("Received ID is an unexpected number: %d", id)
	}
}

// Test listenForIDRequests:
// Test that an error number is returned through the idReceiverChan when there is
// any error from uid.GetNextID() - in this case, the error is ID is at max value).
func TestListenForIDRequestsReceiveErrFromGetNextID(t *testing.T) {
	d1 := []byte(strconv.FormatUint(math.MaxUint64, 10))
    ioutil.WriteFile(persistenceFile, d1, 0644)

	uid.InitPersistenceLayer()
	go listenForIDRequests()

    idReceiverChan := make(chan uint64)
    idRequestChan <- idReceiverChan
    id := <- idReceiverChan

	if id != uint64(0) {
		t.Errorf("Expected error number from channel (zero), but got %d", id)
	}

	os.Remove(persistenceFile)
}
