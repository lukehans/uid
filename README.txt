-------------------------------------------------------------------------------
-- DESCRIPTION
-------------------------------------------------------------------------------

This project is a sample project to generate unique IDs via a REST service in
the Go language. Unique IDs are 64 bit positive integers.

Features:

   + Multithreaded
       Supports multiple concurrent users while still guaranteeing that each
       returned ID is unique.

   + Restart tolerant
       A persistence layer keeps newly generated IDs unique after restarting.

   + Graceful shutdown
       On shutdown, the http server prevents new service requests, but existing
       threads are given 10 seconds to complete.

   + Fault tolerant
       Several checks are performed to ensure the application can function 
       normally, such as verifying the ID pool is not exhausted, verifying 
       write access to the persistence file, verifying persistnce data is not
       corrputed, etc.


-------------------------------------------------------------------------------
-- PREREQUISUTES
-------------------------------------------------------------------------------

1. Install Git

2. Install Go


-------------------------------------------------------------------------------
-- INSTALLATION/EXECUTION INSTRUCTIONS
-------------------------------------------------------------------------------

1. Navigate to 'src' in your $GOPATH
cd $GOPATH/src

2. Clone the git repository
git clone https://github.com/lukehans/uid.git

3. Compile packages and dependencies
go build gouid/uid/uid.go
go build gouid/uidclient/uidclient.go
go build gouid/main.go

4. Test packages
cd $GOPATH/src/gouid/uid
go test
cd $GOPATH/src/gouid
go test

5. [Linux/Mac only] Increase ulimit. If you cannot modify ulimit, decrease the
   variable 'totalReq' in uidclient.go to 100.
sudo ulimit -n 1000

6. Compile and run the main program.
cd $GOPATH/src/uid
go run main.go

7. Open a new terminal to run the test client.

8. [Linux/Max only] Increase ulimit
sudo ulimit -n 1000

9. Compile and run the test client program.
<in the new terminal>
cd $GOPATH/src/uid/uidclient
go run uvclient.go

7. When ready, kill the processes
Ctrl+C


---------------------------------------------------------------
-- REST API DOCUMENTATION
---------------------------------------------------------------

Note: The application will boot up on port 8080.

URL Path              : /uid
HTTP Request Method   : POST
HTTP Request Body     : <Empty>
HTTP Success Response : JSON Object of the following format:
			{"value":<NUMBER>}
						
			Example Response Body:  
			{"value":1234}
Success Status Code   : 201 Created

