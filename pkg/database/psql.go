package database

import (
	
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	utils "github.com/rcarrata/rck-auth/pkg/utils"

)

// User store User details
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Connect to the Postgresql Database
// TODO: Use Viper handle the psql parameters
func GetDatabase() *gorm.DB {

	// DB params are first retrieved from env variables and if not, default params are applied
	db_name := utils.getEnv("DB_NAME", "postgres")
	db_pass := utils.getEnv("DB_PASS", "1312")
	db_host := utils.getEnv("DB_HOST", "127.0.0.1")
	databaseurl := "postgres://postgres:" + db_pass + "@" + db_host + "/" + db_name + "?sslmode=disable"

	// fmt.Println(databaseurl)
	connection, err := gorm.Open(db_name, databaseurl)

	sqldb := connection.DB()

	// Check the Database URL
	if err != nil {
		logrus.Fatalln("Wrong database url")
	}

	// Check the connection towards the Postgresql
	if err := sqldb.Ping(); err != nil {
		logrus.Fatalln("Error in make ping the DB " + err.Error())
		return nil
	}

	logrus.Info("DB connected")
	return connection

}


// Add the InitialMigration for the DB
func InitialMigration() {
	connection := GetDatabase()
	defer Closedatabase(connection)
	connection.AutoMigrate(User{})
	// CreateRecord(connection, User)
}

// Close the database connection opened
func Closedatabase(connection *gorm.DB) {
	// Only for debug
	// log.Println("Closing DB connection")
	sqldb := connection.DB()
	sqldb.Close()
}

// Function to test the generation of records in the DB
func CreateRecord(db *gorm.DB) {
	user := User{Name: "Rober", Email: "rober@test.com", Password: "test", Role: "Admin"}
	result := db.Create(&user)

	if result.Error != nil {
		logrus.Fatalln("Not able to generate the record")
	}
}

// Query records in example function
func QueryRecord(db *gorm.DB, user User) {
	result := db.First(&user)

	if result.Error != nil {
		logrus.Println("Not record present")
	}
}
