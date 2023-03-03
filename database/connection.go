package database

import (
	"database/sql"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"server/configs"
	"server/models"
	"server/pkg/utils"
)

import (
	_ "github.com/lib/pq"
)

type Config struct {
	Config struct {
		Database struct {
			Environment struct {
				POSTGRES_HOST     string `yaml:"POSTGRES_HOST"`
				POSTGRES_DB       string `yaml:"POSTGRES_DB"`
				POSTGRES_USER     string `yaml:"POSTGRES_USER"`
				POSTGRES_PASSWORD string `yaml:"POSTGRES_PASSWORD"`
				Ports             int64  `yaml:"ports"`
			} `yaml:"environment"`
		} `yaml:"database"`
	} `yaml:"config"`
}

func Connection() *sql.DB {
	c := &Config{}
	result, err := utils.ReadFileYaml("./configs/config_server.yaml")

	if err != nil {
		fmt.Println(err)
	} else {
		mapstructure.Decode(*result, c)
	}

	dbInfo := c.Config.Database.Environment

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.POSTGRES_HOST, dbInfo.Ports, dbInfo.POSTGRES_USER, dbInfo.POSTGRES_PASSWORD, dbInfo.POSTGRES_DB)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
	// sqlite
	//log.Println("Creating sqlite-database.db...")
	//file, err := os.Create("sqlite-database.db") // Create SQLite file
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//file.Close()
	//log.Println("sqlite-database.db created")
	//
	//sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	//defer sqliteDatabase.Close()                                     // Defer Closing the database
	//return sqliteDatabase

}
func InitDatabase() {
	db := Connection()
	createTable(
		db, models.User{},
	)

}

func checkTableExist(name string, Db *sql.DB) (exist bool) {

	query := fmt.Sprintf(
		"SELECT EXISTS ( SELECT * FROM %s.%s)", configs.SchemaName, name)
	_, err := Db.Exec(query)
	if err != nil {
		return false
	}
	return true
}

func createTable(Db *sql.DB, tables ...interface{}) {
	for _, v := range tables {
		t := reflect.TypeOf(v)
		TableName := t.Name()
		tableExist := checkTableExist(TableName, Db)
		if !tableExist {
			var fields = "id SERIAL, "
			if reflect.ValueOf(v).Kind() == reflect.Struct {
				value := reflect.ValueOf(v)
				numberOfFields := value.NumField() // get number of struct fields
				for i := 1; i < numberOfFields; i++ {
					fieldType := fmt.Sprintf("%v", value.Field(i).Kind())
					jsonValueTag := t.Field(i).Tag.Get("json")
					fields += " " + jsonValueTag + " "
					if fieldType == "string" {
						fields += "text, "
					} else if fieldType == "int" {
						fields += "INTEGER, "
					} else if fieldType == "bool" {
						fields += "BOOL, "
					}
				}
			}

			query := fmt.Sprintf("CREATE TABLE %s.%s ( %s "+
				"created_at date not null, "+
				"updated_at date,"+
				"deleted_at date)", configs.SchemaName,
				TableName, fields)
			_, err := Db.Query(query)
			if err != nil {
				fmt.Printf("creating %s table error", TableName)

			}
			fmt.Printf("create %s table successfully", TableName)
		} else {
			fmt.Printf("%s table existing", TableName)
		}
	}
}
