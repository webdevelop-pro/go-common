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

	"github.com/pkg/errors"
)

func (f FixturesManager) updateSeqNum(fixture Fixture) error {
	query := fmt.Sprintf(
		"select setval('%s_id_seq',(select max(id)+1 from %s));",
		fixture.table, fixture.table,
	)

	_, err := f.db.Exec(context.Background(), query)

	return err
}

func (f FixturesManager) LoadFixtures(fixtures []Fixture) error {
	for _, fixture := range fixtures {
		err := f.LoadFixture(fixture.table, fixture.filePath)
		if err != nil {
			return err
		}

		err = f.updateSeqNum(fixture)
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
					if fv, ok := v.(float64); ok {
						v = fv
					}
				// handle json/jsonb
				case reflect.TypeOf(v).Kind() == reflect.Map:
					b := bytes.NewBuffer(nil)
					err = json.NewEncoder(b).Encode(v)
					if err != nil {
						file.Close()
						e := fmt.Errorf("failed to encode json field %s: %w", k, err)
						log.Fatal(e.Error())
					}
					v = b.String()
				}
			}
			vals = append(vals, v)
		}
		q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, strings.Join(fields, ","), valuePlaceholders)
		ct, err := f.db.Exec(context.Background(), q, vals...)
		if err != nil {
			e := fmt.Errorf("failed to insert row #%d: %w\n\nquery: %s\n\nvals: %+v", rowID, err, q, vals)
			log.Fatal(e.Error())
		}
		totalInserted += ct.RowsAffected()
	}

	file.Close()

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

	if _, err := os.Stat(currentDir + "/tests"); os.IsNotExist(err) {
		if path.Base(currentDir) != "tests" {
			currentDir = path.Dir(currentDir)
		}

		currentDir = path.Dir(currentDir)
	}

	return path.Join(currentDir, relativePath)
}
