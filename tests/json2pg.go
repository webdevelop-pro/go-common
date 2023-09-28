package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (f FixturesManager) LoadFixtures(fixtures []Fixture) error {
	for _, fixture := range fixtures {
		err := f.LoadFixture(fixture.table, fixture.filePath)
		if err != nil {
			return err
		}
	}
	// fmt.Printf("applied %d fixtures", len(fixtures))
	return nil
}

func (f FixturesManager) LoadFixture(tableName, fileName string) error {
	file, err := os.Open(GetAbsPath("/tests/" + fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	var inputData []map[string]interface{}
	err = json.NewDecoder(file).Decode(&inputData)
	if err != nil {
		return err
	}
	if len(inputData) == 0 {
		return err
	}

	cols, err := f.columns(f.cfg.Database, tableName)
	if err != nil {
		return err
	}

	errors := make([]error, 0)
	var totalInserted int64
	for rowID, row := range inputData {
		var valuePlaceholders string
		fields := make([]string, 0, len(row))
		vals := make([]interface{}, 0, len(row))
		var i int
		for k, v := range row {
			if _, ok := cols[k]; !ok {
				continue
			}
			i++
			if i > 1 {
				valuePlaceholders += ","
			}
			valuePlaceholders += "$" + strconv.Itoa(i)
			fields = append(fields, `"`+k+`"`)

			if v != nil {
				switch {
				// handle number -> timestamp
				case reflect.TypeOf(v).Kind() == reflect.Float64 && strings.Contains(cols[k], "timestamp"):
					v = time.Unix(int64(v.(float64)), 0)
				// handle json/jsonb
				case reflect.TypeOf(v).Kind() == reflect.Map:
					b := bytes.NewBuffer(nil)
					err = json.NewEncoder(b).Encode(v)
					if err != nil {
						e := fmt.Errorf("Failed to encode json field %s: %v\n", k, err)
						log.Fatal(e.Error())
						errors = append(errors, e)
					}
					v = b.String()
				}
			}
			vals = append(vals, v)
		}
		q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, strings.Join(fields, ","), valuePlaceholders)
		ct, err := f.db.Exec(context.Background(), q, vals...)
		if err != nil {
			e := fmt.Errorf("Failed to insert row #%d: %v\n\nquery: %s\n\nvals: %+v\n", rowID, err, q, vals)
			log.Fatal(e.Error())
			errors = append(errors, e)
		}
		totalInserted += ct.RowsAffected()
	}
	// fmt.Printf("Inserted %d rows into %s\n", totalInserted, tableName)
	if len(errors) > 0 {
		fmt.Printf("Errors occured during execution (%d):\n", len(errors))
		for i, err := range errors {
			fmt.Printf("#%d\n%s\n", i, err)
		}
		os.Exit(1)
	}
	return nil
}

func (f FixturesManager) columns(dbName, tableName string) (map[string]string, error) {
	rows, err := f.db.Query(
		context.Background(),
		`SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = $1 AND table_catalog=$2`,
		tableName, dbName,
	)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}
	defer rows.Close()
	cols := make(map[string]string)
	for rows.Next() {
		var n, t string
		err = rows.Scan(&n, &t)
		if err != nil {
			return nil, errors.Wrap(err, "scan failed")
		}
		cols[n] = t
	}
	return cols, nil
}

func GetAbsPath(relativePath string) string {
	currentDir, _ := os.Getwd()

	for path.Base(currentDir) != "tests" {
		currentDir = path.Dir(currentDir)
	}

	currentDir = path.Dir(currentDir)

	return path.Join(currentDir, relativePath)
}
