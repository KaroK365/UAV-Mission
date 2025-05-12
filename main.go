package main

import (
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB 

//UAV struct to represent UAV data

type UAV struct{
	ID  		uint      `json:"id" gorm:"primaryKey"`
	NAME		string	  `json:"name"`
	SPEED		int		  `json:"speed"` //in km/h	
	FUEL    	int		  `json:"fuel"` // in litres
	DESCRIPTION	string	  `json:"description"`
	Missions	[]Mission `json:"missions" gorm:"foreignKey:UAVID"`
}

type Mission struct{
	ID 			uint		`json:"id" gorm:"primaryKey"`
	Name		string		`json:"name"`
	UAVID		uint		`json:"uav_id"`
	UAV			UAV			`gorm:"foreignKey:UAVID" json:"-"`
	Waypoints	[]Waypoint	`json:"waypoints" gorm:"foreignKey:MissionID"`
}

type Waypoint struct{
	ID 			uint		`json:"id" gorm:"primaryKey"`
	MissionID	uint		`json:"mission_id"`
	Latitude 	float64		`json:"latitude"`
	Longitude 	float64		`json:"longitude"`
	Altitude	float64		`json:"altitude"`
}

// Initializing database connection
func init(){
	var err error

	err = godotenv.Load()
	if err!= nil{
		log.Fatal("Error loading '.env' file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	//Connect database
	db, err = gorm.Open("postgres", "host=localhost user=" + dbUser + " dbname = "+ dbName +" password=" + dbPassword + " sslmode=disable")

	if err!=nil{
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&UAV{}, &Mission{}, &Waypoint{}) // Automatically migrate UAV model
}

func main() {
    router := gin.Default()
	router.Use(cors.Default())

    // Test route to check if server is running
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

	//Route to create a new UAV
	router.POST("/uav", func(c *gin.Context) {
		var uav UAV
		if err:= c.BindJSON(&uav); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid data"})
			return 
		}

		if uav.NAME == "" || uav.SPEED <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{ "error" : "Name is required and speed must be positive"})

			return
		}

		db.Create(&uav)
		c.JSON(http.StatusCreated, gin.H{"message":"UAV created", "uav":uav})
	})

	//Route to create a Mission
	router.POST("/mission", func(c *gin.Context) {
		var mission Mission
		if err:= c.BindJSON(&mission); err!= nil{
			c.JSON(400, gin.H{"error" : "Invalid data"})
			return
		}

		db.Create(&mission)
		c.JSON(http.StatusCreated, gin.H{"message":"Mission created", "mission": mission})
	})

	//Route to create a Waypoint
	router.POST("/waypoint", func (c *gin.Context)  {
		var waypoint Waypoint
		if err := c.BindJSON(&waypoint); err != nil{
			c.JSON(400, gin.H{"error" : "Invalid data"})
			return
		} 

		db.Create(&waypoint)
		c.JSON(http.StatusCreated, gin.H{"message":"Waypoint Created", "waypoint": waypoint})
	})


	// Get Uavs
	router.GET("/uav", func(c *gin.Context){
		var uavs []UAV
		db.Preload("Missions").Preload("Missions.Waypoints").Find(&uavs)
		c.JSON(http.StatusOK, uavs)
	})

	//Get Uav by id
	router.GET("/uav/:id", func(c *gin.Context) {
		id := c.Param("id")
		var uav UAV

		if err := db.Preload("Missions").Preload("Missions.Waypoints").Find(&uav, id).Error; err!=nil{
			c.JSON(http.StatusNotFound, gin.H{"error" : "UAV not found"})
			return
		}

		c.JSON(http.StatusOK, uav)
	})

	// Get Missions
	router.GET("/mission", func(c *gin.Context){
		var missions []Mission
		db.Preload("UAV").Preload("Waypoints").Find(&missions)
		c.JSON(http.StatusOK, missions)
	})

	//Get Missions by id
	router.GET("/mission/:id", func(c *gin.Context) {
		id := c.Param("id")
		var mission Mission
		if err:=db.Preload("UAV").Preload("Waypoints").First(&mission, id).Error; err!=nil{
			c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
			return
		}
		c.JSON(http.StatusOK, mission)
	})

	// Summary for mission
	router.GET("/mission/:id/summary",func(c *gin.Context){
		id:= c.Param("id")

		var mission Mission

		// Fetching mission details 
		if err:= db.Preload("Waypoints").Preload("UAV").First(&mission, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error":"Mission not found"})
			return
		}

		// calculating total distance between waypoints

		totalDistance := 0.0
		waypoints := mission.Waypoints

		for i:= 0; i < len(waypoints) - 1; i++{
			wp1 := waypoints[i]
			wp2 := waypoints[i+1]
			totalDistance += haversine(wp1.Latitude, wp1.Longitude, wp2.Latitude, wp2.Longitude)
		}

		//Estimate fuel usage
		fuelEfficiency := 0.2
		estimatedFuel := totalDistance * fuelEfficiency

		//Travel time calculation
		speed:= float64(mission.UAV.SPEED)
		estimatedTime := totalDistance / speed

		c.JSON(http.StatusOK, gin.H{
			"mission_id" : mission.ID,
			"uav" : mission.UAV.NAME,
			"total_distance" : totalDistance,
			"fuel_required" : estimatedFuel,
			"travel_time" : estimatedTime * 60,
		})

	})

    // Start the server on port 8080
    router.Run(":8080")
}


func haversine(lat1, long1, lat2, long2 float64) float64{
	const R = 6371

	dLat:= (lat2 - lat1) * math.Pi / 180
	dLon:= (long2 - long1) * math.Pi / 180

	a:= math.Sin(dLat/2) * math.Sin(dLat/2) +
		math.Cos(lat1 * math.Pi/180) * math.Cos(lat2 * math.Pi/180) *
		math.Sin(dLon/2) * math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))

	return R * c
}

