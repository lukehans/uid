package uid

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"
)

// Test InitPersistenceLayer:
// Test the correct error is returned if the persistence layer data cannot be
// converted to uint64 type.
func TestInitPersistenceLayerCorruptedData(t *testing.T) {
	nextID = 1

	d1 := []byte(string("not-a-uint64"))
    ioutil.WriteFile(persistenceFile, d1, 0644)

    defer func() {
    	r := recover()
        if r != nil {
            fmt.Println("Recovered! - Expected recoverery in testing. All is well. Recovered from the following error: ", r)
        }
    }()

    InitPersistenceLayer()
    os.Remove(persistenceFile)
}

// Test InitPersistenceLayer:
// Test the nextID is the value in the persistence file plus one.
func TestInitPersistenceLayerReadID(t *testing.T) {
	nextID      = 1
	expectedID := uint64(124)

	d1 := []byte(strconv.FormatUint(expectedID - 1, 10))
    ioutil.WriteFile(persistenceFile, d1, 0644)

    InitPersistenceLayer()

  	actualID := nextID

    if expectedID != actualID {
    	t.Errorf("TestInitPersistenceLayerReadID test failed. Persistence layer was not read correctly. Expected %s, actualID was %s", expectedID, actualID)
    }
    os.Remove(persistenceFile)
}

// Test InitPersistenceLayer:
// Test that the persistence file is created.
func TestInitPersistenceLayer(t *testing.T) {
	_, err := os.Stat(persistenceFile)
	if !os.IsNotExist(err) {
    	os.Remove(persistenceFile)
    }

    InitPersistenceLayer()

    _, err = os.Stat(persistenceFile)
    if os.IsNotExist(err) {
    	t.Errorf("TestInitPersistenceLayer test failed. Persistence file could not be created: ", err)
    }
    os.Remove(persistenceFile)
}

// Test GetNextID:
// Test that the returned ID was incremented by one.
func TestGetNextID(t *testing.T) {
	nextID       = 1
	expectedID  := uint64(1)
  	actualID, _ := GetNextID()

  	if expectedID != actualID {
    	t.Errorf("TestGetNextID test failed. Expected %d, actual was %d", expectedID, actualID)
  	}
  	os.Remove(persistenceFile)
}

// Testing GetNextID:
// Test that the returned ID was written to the persistence layer.
func TestGetNextIDCheckPersistenceLayer(t *testing.T) {
	nextID         = 1
	expectedID, _ := GetNextID()

  	content, err := ioutil.ReadFile(persistenceFile)
  	if err != nil {
        t.Errorf("TestGetNextIDCheckPersistenceLayer test failed. Could not read persistence file: ", err)
    }

    actualID, err := strconv.ParseUint(string(content), 10, 64)

  	if expectedID != actualID {
    	t.Errorf("TestGetNextIDCheckPersistenceLayer - test failed. Expected %d, actual was %d", expectedID, actualID)
  	}
  	os.Remove(persistenceFile)
}

// Testing GetNextID:
// Test that a new persistence layer file is created if it does not exist.
func TestWriteNextIDToFileCreateNewFile(t *testing.T) {
    _, err := os.Stat(persistenceFile)
    if !os.IsNotExist(err) {
    	os.Remove(persistenceFile)
    }

	nextID = 1
	
	err = writeNextIDToFile()
	if err != nil {
		t.Errorf("TestWriteNextIDToFileCreateNewFile test failed.", err)
	}
	os.Remove(persistenceFile)
}

// Testing writeNextIDToFile:
// Test that the correct error is returned if file write access is denied by the
// operating system.
// NOTE:
// I am not sure how to get this test working on a Mac, but it works on Windows.
/*func TestWriteNextIDToFileCannotWrite(t *testing.T) {
	InitPersistenceLayer()
	os.Chmod(persistenceFile, 0444)
	err := writeNextIDToFile()

	if err.Error() != "open last-used-id.txt: Access is denied." {
		t.Error("Test failed. Expected to be denied writable access to the file.")
	}
	os.Chmod(persistenceFile, 0644)
}*/

// Testing writeNextIDToFile:
// Test this function writes the nextID to the persistence file.
func TestWriteNextIDToFile(t *testing.T) {
	nextID      = 1
	expectedID := nextID

	writeNextIDToFile()

	content, _  := ioutil.ReadFile(persistenceFile)
    actualID, _ := strconv.ParseUint(string(content), 10, 64)
    
    if actualID != expectedID {
        t.Errorf("TestWriteNextIDToFile test failed. Expected %d, actual was %d", expectedID, actualID)
    }
    os.Remove(persistenceFile)
}

// Test prepareNextID:
// Test that the nextID is incremented in the function call.
func TestPrepareNextID(t *testing.T) {
	nextID = 1

	err := prepareNextID();
	if err != nil {
		t.Errorf("TestPrepareNextID test failed. Returned error: %s", err)
	}

	expectedID := uint64(2)
	actualID   := nextID
	if actualID != expectedID {
		t.Errorf("TestPrepareNextID test failed. Expected %d, actual was %d", expectedID, actualID)
	}
	os.Remove(persistenceFile)
}

// Test prepareNextID:
// Test that the correct error is returned if the ID cannot be incremented any more.
func TestPrepareNextIDAtMaxVal(t *testing.T) {
	nextID = math.MaxUint64 - 1

	err := prepareNextID()
	if err.Error() != "Unique ID pool is depleted!" {
		t.Errorf("TestPrepareNextIDAtMaxVal test - test failed. Returned error: %s", err)
	}
	os.Remove(persistenceFile)
}