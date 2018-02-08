package uid

import (
    "fmt"
    "io/ioutil"
    "log"
    "math"
    "os"
    "strconv"
)

const persistenceFile = "last-used-id.txt"

var (
     nextID uint64 = 1
     infoLog       = log.New(os.Stdout, "INFO: ", log.LstdFlags)
     errLog        = log.New(os.Stdout, "ERROR: ", log.LstdFlags)
)

// InitPersistenceLayer initializes the persistence layer of the application by:
//   1) creating the 'last-used-id.txt' file (if it does not exist),
//   2) setting the file permissions to validate read/writability,
//   3) reading the last used ID from the file, and
//   4) setting the nextID value appropriately
// If there are any errors with the persistence layer, including the program
// cannot create/read the file, set file permissions, or the ID data in the file
// cannot be converted to uint64, then the program exits.
func InitPersistenceLayer() {
    // Check if any persistence layer data exists. Create persistence layer
    // file if it does not exist.
    _, err := os.Stat(persistenceFile);
    if os.IsNotExist(err) {
        file, err := os.Create(persistenceFile)
        if err != nil {
            errLog.Printf("Cannot create persistence layer data.\n", err)
            panic(err)
        }
        file.Close()
    }

    // Make sure file permissions are set appropriately.
    err = os.Chmod(persistenceFile, 0644)
    if err != nil {
        errLog.Printf("Persistence layer file permissions cannot be verified.\n", err)
        panic(err)
    }

    // Read the last used ID from storage.
    content, err := ioutil.ReadFile(persistenceFile)
    if err != nil {
        errLog.Printf("Cannot read persistence layer data.\n", err)
        panic(err)
    }

    // Prepare the next unique ID for service.
    if len(content) != 0 {
        lastUsedID, err := strconv.ParseUint(string(content), 10, 64)
        if err != nil {
           errLog.Printf("Persistence layer data is corrupted.\n", err)
           panic(err)
        }
        nextID = lastUsedID + 1
    }
}

// GetNextID() returns a unique ID. Calls to this method should be synchronized. Before
// returning the unique ID, the ID is written to a file to keep memory of IDs that have
// already been used.
func GetNextID() (uint64, error) {
    id  := nextID;
    err := writeNextIDToFile()
    if err != nil {
        errLog.Printf("Error writing to file. ", err)
        return id, err
    }
    
    err = prepareNextID()
    if err != nil {
        errLog.Printf("Error preparing next ID. ", err)
    }
    return id, err
}

// writeNextIDToFile writes the next usable unique ID to persistent storage (a file). Every
// time a new ID is requested, the ID is saved to the file 'last-used-id.txt' before returning
// the ID to the client. This is the persistent storage layer, allowing the ID program to 
// continue generating unique IDs after a restart.
func writeNextIDToFile() error {
    // Convert the nextID to bytes.
    d1 := []byte(strconv.FormatUint(nextID, 10))

    // Write nextID bytes to the file.
    return ioutil.WriteFile(persistenceFile, d1, 0644)
}

// prepareNextID increments the nextID variable so it will be unique for the next request.
func prepareNextID() error {
    nextID++
    if nextID == math.MaxUint64 {
        return fmt.Errorf("Unique ID pool is depleted!")
    }
    return nil
}
