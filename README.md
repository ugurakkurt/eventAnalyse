# EventAnalyser
An event analyser RestAPI written in Golang. It is an analyze tool that can be used
for storing and analyzing some specified type data using MySQL database technology.
If you have some data for requests of specific users using specific apps and if you 
want to analyze your response time per app in this data, then you can use this app
by storing your data toa JSON file as event definitions that composes of userIDs,
ApiIds and response time info.

Installation:

Assuming that go is set up already, install following three dependencies
before moving on. Here go-chart library is used to analyze stored event data.

go get "github.com/go-sql-driver/mysql"

go get "github.com/gin-gonic/gin"

go get github.com/wcharczuk/go-chart



For client-server connection, Gin is used in this API. Gin is a HTTP web
framework written in Go (Golang). It features a Martini-like API with much 
better performance -- up to 40 times faster. 

Before running the API, firstly create a database in MySQL explicitly. Then give
the name of the database, MySQL username and password in  const variables section
of createTables.go file. Also make sure that you are using port 3306 in MySQL. If
that port is used by another application in your computer, you can change the port
name in const variable list as well. Also 3232 port is used in the connection 
between client and server using Gin framework. If another application is using that
port in your computer, you need to change that port in const variables sections in
both client and server package.

Now your server is ready to go.
You can run the server. Gin-gonic creates the environment and port connection. When
the app is run, at first database relations are created in MySQL in main method if
they had not been created before by calling createTables() method just in the first
line of main method. Then server listens for GET and POST requests. POST request is
made for adding events and GET is for analyzing events.

In client package, there is a JSON data of 100 test events that can be sent to
server to be stored in database automatically. You can do it by simply running 
client.exe in client package. You can send your own event data by manipulating 
the pre-created JSON file that is in the client package.


Now run the the client side. Event logs from
the JSON file will be automatically send to server and inserted to database using
POST requests. JSON file can be updated to send more event definitions automatically.
An example request can be defined as follows:

http://localhost:3232/event/?api_key=10&user_id=1&unix_timestamp=1485084687

As a second option, event data can be sent via POST request manually using Postman
application. You can download Postman from the following link:
https://www.getpostman.com/


Running the client side and having some event storage in database, now analysis part 
can be done sending a GET request to server. Just open your browser and type the
following:

http://localhost:3232/event (change port number if yours is different)

Now you will see a chart that shows percentages of events corresponding to some specified
buckets in terms of their processing times. You can analyze your event data using this
chart and customize it for yourself if you wish.








