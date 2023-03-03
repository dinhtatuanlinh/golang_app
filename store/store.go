package store

import (
	"database/sql"
	"fmt"
	"reflect"
	"server/configs"
	"strings"
)

type IDbManipulation[T any] interface {
	CreateSchema(name string)
	getSchemas() map[string]string
	FindAll() ([]map[string]interface{}, error)
}

type DbManipulation[T any] struct {
	DB *sql.DB
}

func NewDbManipulation[T any](db *sql.DB) IDbManipulation[T] {
	return &DbManipulation[T]{DB: db}
}

func (DbManipulation *DbManipulation[T]) CreateSchema(name string) {
	query := fmt.Sprintf("CREATE SCHEMA%s", name)
	result, err := DbManipulation.DB.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func (DbManipulation *DbManipulation[T]) getSchemas() map[string]string {
	query := fmt.Sprintf("SELECT nspname FROM pg_catalog.pg_namespace")
	rows, err := DbManipulation.DB.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var schemas []string
	results := make(map[string]string)
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			panic(err)
		}
		schemas = append(schemas, schema)
	}
	for _, v := range schemas {
		results[v] = v
	}
	return results
}
func (DbManipulation *DbManipulation[T]) GetById() (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM user WHERE deleted_at IS NULl;")
	rows, err := DbManipulation.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return rows, nil
}
func (DbManipulation *DbManipulation[T]) FindAll() ([]map[string]interface{}, error) {
	var data T
	t := reflect.TypeOf(data)
	TableName := t.Name()
	query := fmt.Sprintf("SELECT * FROM %s.\"%s\" WHERE deleted_at IS NULl", configs.SchemaName, strings.ToLower(TableName))
	rows, err := DbManipulation.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	row := make([]interface{}, len(cols))
	rowPtr := make([]interface{}, len(cols))
	val := make(map[string]interface{})
	datas := []map[string]interface{}{}
	for i := range row {
		rowPtr[i] = &row[i]
	}

	for rows.Next() {
		err = rows.Scan(rowPtr...)
		if err != nil {
			fmt.Println("cannot scan row:", err)
		}
		for i, v := range row {
			var key = cols[i]
			//t := fmt.Sprintf("%v", reflect.TypeOf(v))
			//switch t {
			//case "int64":
			//	v, _ = v.(int)
			//	fmt.Println(v)
			//case "string":
			//	v = fmt.Sprintf("%v", v)
			//case "time.time":
			//	v = v.(time.Time)
			//}
			val[key] = v
		}
		datas = append(datas, val)
	}

	return datas, nil
}

//func getTypes()
