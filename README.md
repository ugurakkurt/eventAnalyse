# EventAnalyser
Event analyser is a RestAPI written in Golang. It is an analyzer tool that can be used
for storing and analyzing some specif type of data using MySQL database technology.
If you have some data of requests of specific users using specific apps and if you 
want to analyze the response times for each of these defined apps for the request of users,
then you can use this app by storing your data in a JSON file as event definitions that is in
userID,ApiId and responseTime(unix timestamp) format.

Installation:

Assuming that go is set up already, install following three dependencies
before moving on. 

go get "github.com/go-sql-driver/mysql"

go get "github.com/gin-gonic/gin"

go get github.com/wcharczuk/go-chart



Here go-chart library is used to analyze stored event data and Gin is used in this API
for client-server connection. It is a HTTP web framework written in Go (Golang).
If your data is big, processing time is
important and here Gin serves you well with its performance. It is simple 
to learn how Gin works but hard to master. It has a huge community of 
followers and it is a mature framework having a development history of 2 years
which probably makes it almost bug free Via Gin framework you can render
your plots into your web-page easily unlike some other HTTP frameworks like Mux.
This app requires parsing parameters from URL and it is easy to parse parameters
with Gin. Gin has the same philosophy like httpRouter as it is built on top of 
it which makes its performance better thanks to fast data passing between middlewares.
So if performance is everything for you, Gin can serve you well.

Having told why I chose Gin, now, before running the API, firstly create a database
in MySQL explicitly. Then give the name of the database, MySQL username and password
in  const variables section of createTables.go file. Also make sure that you are using
port 3306 in MySQL. If that port is used by another application in your computer, you
can change the port name in const variable list as well. Also 3232 port is used in the
connection between client and server using Gin framework. If another application is
using that port in your computer, you need to change that port which is stored in port
variable of main function in server.go file. in both client and server package.

Now your server is ready to go.
Run the server. Gin-gonic will create the environment and start listening for requests. When
the app is run, at first database tables are created in MySQL in main method if
these tables had not been created before by calling createTables() method just in the first
line of main method. Then server listens for GET and POST requests from the same port.
POST request is made for adding events and GET is for analyzing events.

Now having server running, you can run the client side. In client package, 
there is a JSON data of 100 test event logs that can be sent to
server within POST reqests to be stored in database automatically. You can do it by simply running 
client.exe in client package. Also, you can send your own event data by manipulating 
the pre-created JSON file that resides in the client package. JSON has a readable,
user-friendly format and less grammer. JSON is processed quite fast as it has 
a simpler structure so less markup overhead compated to XML. 

As a second option, event data can be sent via POST request manually. As an option, Postman
application can be used for sending POST requests. You can download Postman 
from the following link:
https://www.getpostman.com/

An example POST request can be defined as follows:

http://localhost:3232/event/?api_key=10&user_id=1&unix_timestamp=1485084687
This POST request above, will be processed by server. Server will firstly check for error 
conditions. All three parameters must be found in the URL. If at least one of them is missing, 
server will return a proper error message. Also if any of the parameters is not integer, then server 
will return an error message as well. Assuming all data is sent correctly like in the example
request above, server will continue its work after checking these error conditions. If this user
id did not make any request before, its id will be saved in users table. Then it will check api id.
If api id is not one of the predefined api's then server returns error message again. There are
predefined api id's and for each api there is a table created in database to store events for 
related api. In other words events that belong to the same API key are grouped in the same table.
So the info of which user made the request for the related api and how much time passed for response
will be saved in the related api's table which would be "api_key_10" table for this example request above.
At last, response time of server will be saved in responsetimes table for future analysis. although actually
response times are seved in related api's table, I also added these response times specifically in responsetimes
table again to utilize data locality which is needed in a data analysis where all the response times info 
is needed without considering api id's. So it would be faster retrive all the response times if only
one table was processed. If you have enough storage area in your database, this approach is better in case
of performance for retriving all response times data for analysis. It adds a small overhead during event
insertion but it is a better approach if you need performance during analysis part. So this is actually 
a trade-off decision that needed to be made.

Now, Running the client side and having some event storage in database, analysis part 
can be done sending a GET request to server. Just open your browser and type the
following:

http://localhost:3232/event (change port number if yours is different)

Now you will see a chart that shows percentages of events corresponding to some specified
buckets in terms of their processing times. You can analyze your event data using this
chart and customize it for yourself if you wish.


