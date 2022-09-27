package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	// If already connected, just return that connection
	if DB != nil {
		return DB, nil
	}

	dbAddress := viper.Get("dbAddress").(string)
	dbUser := viper.Get("dbUsername").(string)
	dbPass := viper.Get("dbPassword").(string)
	dbPort := viper.Get("dbPort").(string)
	var dbSSL string
	// Needs enabled/disabled
	if viper.Get("dbSSL").(bool) {
		dbSSL = "enable"
	} else {
		dbSSL = "disable"
	}

	// Change this to match docker settingss
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=GMT",
		dbAddress, dbUser, dbPass, dbUser, dbPort, dbSSL)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to %v db at '%v:%v' (ssl=%v) with username '%v'",
			dbUser, dbAddress, dbPort, dbSSL, dbUser)
		log.Fatalf("DB connection error: %v", err)

		return nil, err
	}

	log.Printf("Connected to %v db at '%v:%v' (ssl=%v) with username '%v'",
		dbUser, dbAddress, dbPort, dbSSL, dbUser)

	return DB, nil
}

func MigrateSchema() {
	if DB == nil {
		log.Printf("Could not migrate database schema... DB returned nil.")
	} else {
		log.Printf("Migrating database schema...")
		defer log.Printf("Done migrating database schema.")

		log.Printf("Migrating hosts table...")
		DB.AutoMigrate(&Host{})

		log.Printf("Migrating network interfaces table...")
		DB.AutoMigrate(&NetworkInterface{})

		log.Printf("Migrating remote hosts table...")
		DB.AutoMigrate(&RemoteHost{})

		log.Printf("Migrating network flow table...")
		DB.AutoMigrate(&NetFlow{})

		log.Printf("Migrating users table...")
		DB.AutoMigrate(&User{})
	}
}

func InitializeDB() {
	ConnectDB()
	MigrateSchema()
}
