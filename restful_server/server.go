package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart"
)

//global variables to be used by any function on server side
var db *sql.DB
var stmt *sql.Stmt
var err error
var rows *sql.Rows

const bucketSize = 6

var endPoints = [bucketSize]int{1, 5, 10, 20, 50, 100}
var n = [bucketSize]float64{0, 0, 0, 0, 0, 0} // will be used when defining buckets to analyze response times data

/*var isBucketChanged = 1  ---> isBucketChanged is used as a flag, if it is 1, responses are retrived again from database
for analysis. if you do not want to retrive response times from database if client did not explictly send post reqests and updated
response times list, then use it by uncommenting the codeline above, the if statement(it is in router.GET("/analysis") function
 -around line 205, uncomment isBucketChanged = 0 code-around line 288 and the code block around line 154 in createTables.go )*/

var responseSize = 0 // will be used for response times analysis
var API_KEY = ""
var USER_ID = ""
var UNIX_TIMESTAMP = ""

func main() {
	port := "3232"

	createTables()
	db, err = sql.Open(DB_TYPE, USERNAME+":"+PASSWORD+"@tcp(127.0.0.1:"+DB_PORT+")/"+DB_NAME)
	if err != nil {
		panic(err)
	}
	// make sure db is closed after everything
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	//event parameters definition

	router := gin.Default() //create a default gin fouter

	// POST new event method
	// http://localhost:3232/event/?api_key=10&user_id=1&unix_timestamp=1485084687
	router.POST("/event", func(c *gin.Context) {

		c.Writer.Header().Add("answer", "")
		//retrive the parameters from the post request
		API_KEY = c.Query("api_key")
		USER_ID = c.Query("user_id")
		UNIX_TIMESTAMP = c.Query("unix_timestamp")
		//c.Writer.Header().Add("Answer", "")
		//check error conditions
		if checkErrorCond(c) == false {
			return
		} else { //no error, move on
			// it is okay if user id does not exist, if the user id in the processed event definition does not
			// exist in database, create a new user in the users table
			if doesUserExist(USER_ID) == false {
				//this user id did not exist, so add..
				addUser(USER_ID)
			}
			//now ready to process event
			addEvent(API_KEY, USER_ID, UNIX_TIMESTAMP)
			//randomly generate 1ms-100ms (including 1 and 100ms)
			var randSleepMS = time.Duration(rand.Int31n(100)) + 1 //between 1 and 100(1 and 100 inclusive)
			//pause the execution
			//time.Sleep(time.Millisecond * randSleepMS)
			//save this time info to database
			addResponseTime(randSleepMS.String())
		}

	})

	//used for creating analysis report and printing on web page
	router.GET("/analysis", func(c *gin.Context) {
		var responseTime int
		// if there is an update in response times list, retrive response times again from database
		/*!!NOTE: This algorithm works correctly assuming that responsetimes table is not
		manipulated manually in MySQL!! Otherwise, an update done manually in MySQL is not shown
		unless another update from client package is done within POST request		*/
		//if isBucketChanged == 1 { //may be uncommented if needed
		//first updating data..
		rows, err = db.Query("SELECT responsetime FROM responsetimes")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		//set buckets to 0 before updating them again
		for i := 0; i <= bucketSize-1; i++ {
			n[i] = 0
		}
		responseSize = 0
		//iterate over response times list and increment necessary bucket representations accordingly
		for rows.Next() {
			responseSize = responseSize + 1
			err := rows.Scan(&responseTime)
			checkErr(err)
			if (responseTime > 0) && (responseTime <= endPoints[0]) {
				//var endPoints = [bucketSize]int{1, 5, 10, 20, 50, 100}
				n[0] = n[0] + 1
			}
			if (responseTime > endPoints[0]) && (responseTime <= endPoints[1]) {
				n[1] = n[1] + 1
			}
			if (responseTime > endPoints[1]) && (responseTime <= endPoints[2]) {
				n[2] = n[2] + 1
			}
			if (responseTime > endPoints[2]) && (responseTime <= endPoints[3]) {
				n[3] = n[3] + 1
			}
			if (responseTime > endPoints[3]) && (responseTime <= endPoints[4]) {
				n[4] = n[4] + 1
			}
			if (responseTime > endPoints[4]) && (responseTime <= endPoints[5]) {
				n[5] = n[5] + 1
			}
		}
		//print bucket resresponseSizeults
		for i := 0; i <= bucketSize-1; i++ {
			n[i] = n[i] / float64(responseSize) * 100
		}
		//}
		graph := chart.BarChart{
			Title: "Response Time Analysis / (interval - percentage)\n Total num of events: " + strconv.Itoa(responseSize),
			TitleStyle: chart.Style{
				Show: true,
				Padding: chart.Box{
					Top:    50,
					Left:   300,
					Right:  10,
					Bottom: 100,
				},
				FontSize:    15,
				StrokeWidth: 30,
			},
			Width:    1000, //550
			Height:   750,  //700
			BarWidth: 140,
			XAxis: chart.Style{
				Show:     true,
				FontSize: 13,
			},

			YAxis: chart.YAxis{
				Name:      "The YAxis",
				NameStyle: chart.StyleShow(),
				Style:     chart.StyleShow(),
				TickStyle: chart.Style{
					Show: true,
				},

				Range: &chart.ContinuousRange{
					Min: 0.0,
					Max: 100.0,
				},
			},

			Bars: []chart.Value{
				{Value: n[0], Label: "[0-1)\nValue:\n%" + strconv.FormatFloat(n[0], 'f', 2, 64)},
				{Value: n[1], Label: "[1-5)\nValue:\n%" + strconv.FormatFloat(n[1], 'f', 2, 64)},
				{Value: n[2], Label: "[5-10)\nValue:\n%" + strconv.FormatFloat(n[2], 'f', 2, 64)},
				{Value: n[3], Label: "[10-20)\nValue:\n%" + strconv.FormatFloat(n[3], 'f', 2, 64)},
				{Value: n[4], Label: "[20-50)\nValue:\n%" + strconv.FormatFloat(n[4], 'f', 2, 64)},
				{Value: n[5], Label: "[50-100)\nValue:\n%" + strconv.FormatFloat(n[5], 'f', 2, 64)},
			},
		}
		c.Writer.Header().Set("Content-Type", "image/png")
		err := graph.Render(chart.PNG, c.Writer)
		if err != nil {
			fmt.Printf("Error when rendering chart: %v\n", err)
		}
		message := "msg"
		nick := "nck"
		c.JSON(200, gin.H{
			"message": message,
			"nick":    nick,
		})
		//isBucketChanged = 0 //may be uncommented if needed
	})

	// THE PART BELOW IS DEPRECATED!
	//It can also be used for raw analysis but results are not graphical
	// to analyze results, define responseTime
	/*var responseTime int
	//analyze events within a GET request,
	router.GET("/eventt", func(c *gin.Context) {
		// if there is an update in response times list, retrive response times again from database
		rows, err = db.Query("SELECT responsetime FROM responsetimes")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		//set buckets to 0 before updating them again
		for i := 0; i <= bucketSize-1; i++ {
			n[i] = 0
		}
		responseSize = 0
		//iterate over response times list and increment necessary bucket representations accordingly
		for rows.Next() {
			responseSize = responseSize + 1
			err := rows.Scan(&responseTime)
			checkErr(err)
			if (responseTime > 0) && (responseTime <= endPoints[0]) {
				//var endPoints = [bucketSize]int{1, 5, 10, 20, 50, 100}
				n[0] = n[0] + 1
			}
			if (responseTime > endPoints[0]) && (responseTime <= endPoints[1]) {
				n[1] = n[1] + 1
			}
			if (responseTime > endPoints[1]) && (responseTime <= endPoints[2]) {
				n[2] = n[2] + 1
			}
			if (responseTime > endPoints[2]) && (responseTime <= endPoints[3]) {
				n[3] = n[3] + 1
			}
			if (responseTime > endPoints[3]) && (responseTime <= endPoints[4]) {
				n[4] = n[4] + 1
			}
			if (responseTime > endPoints[4]) && (responseTime <= endPoints[5]) {
				n[5] = n[5] + 1
			}
		}
		fmt.Println("response size: " + strconv.Itoa(responseSize))
		fmt.Println(n[5])
		//print bucket resresponseSizeults
		for i := 0; i <= bucketSize-1; i++ {
			n[i] = n[i] / float64(responseSize) * 100
		}

		for i := 0; i <= bucketSize-1; i++ {
			end := int(n[i])
			if i == 0 {
				c.Writer.WriteString(" 0 < x <=" + strconv.Itoa(endPoints[i]))
			} else {
				c.Writer.WriteString(" " + strconv.Itoa(endPoints[i-1]) + " < x <=" + strconv.Itoa(endPoints[i]))
			}

			c.Writer.WriteString(" (%" + strconv.FormatFloat(n[i], 'f', 2, 64) + "):")
			c.Writer.WriteString("\t\t")

			for j := 1; j <= end; j++ {
				c.Writer.WriteString("|")
			}
			c.Writer.WriteString("\n\n\n\n")
		}

		//These are the exact results
		for i := 0; i <= bucketSize-1; i++ {
			c.Writer.WriteString("%" + strconv.FormatFloat(n[i], 'f', -1, 64) + "\n")
		}
		c.Writer.WriteString("\n")
		//response times list is updated so make the flag(isBucketChanged)false
		//so that when there is a analyis request from client without making any
		//insertion to response times list, it would not retrive response lists
		//from database
		err = rows.Err()
		checkErr(err)
	})*/
	router.Run(":" + port) // listen all requests on port 3232 to insert event definitions
}
func checkErrorCond(c *gin.Context) bool {
	if len(API_KEY) < 1 {
		fmt.Println("api key parameter is missing in URL")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Set("Answer", "input_apiKeyMissing")
		c.String(404, "API KEY parameter is not found in URL")
		return false
	}
	if len(USER_ID) < 1 {
		fmt.Println("user id parameter is missing in URL")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Set("Answer", "input_userIdMissing")
		c.String(404, "USER ID parameter is not found in URL")
		return false
	}
	if len(UNIX_TIMESTAMP) < 1 {
		fmt.Println("unix timestamp parameter is missing in URL")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Set("Answer", "input_timeMissing")
		c.String(404, "UNIX TIMESTAMP parameter is not found in URL")
		return false
	}
	//check if user id parameter is given as integer value
	_, err = strconv.Atoi(USER_ID)
	if err != nil {
		fmt.Println("user id value must be integer")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Set("Answer", "input_userIdError")
		c.String(404, "user id value must be integer")
		return false
	}
	//check if API_KEY is convertible to integer
	_, err = strconv.Atoi(API_KEY)
	if err != nil {
		fmt.Println("The api key must be integer")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Set("Answer", "input_apiKeyTypeError")
		c.String(404, "The api key must be integer")
		return false
	}
	//check if unix timestamp parameter is given as integer value
	_, err = strconv.ParseInt(UNIX_TIMESTAMP, 10, 64)
	if err != nil {
		fmt.Println("Unix timestamp value must be integer")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Set("Answer", "input_TimeError")
		c.String(404, "Unix timestamp value must be integer")
		return false
	}
	if doesApiKeyExist(API_KEY) == false {
		fmt.Println("The api key: ", API_KEY, " does not exist in database.")
		c.Writer.Header().Del("Answer")
		c.Writer.Header().Del("apikey")
		c.Writer.Header().Set("Answer", "apikeyNotFound")
		c.Writer.Header().Set("apikey", API_KEY)
		c.String(404, "Api key: (%s) does not exist", API_KEY)
		return false
	}
	return true
}
