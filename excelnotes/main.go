package main

import (
	"flag"
	"fmt"
	"goNotes/dbnotes/dbhelper"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var dbIP = flag.String("dbIP", "127.0.0.1", "db ip")
var dbPort = flag.Int("dbPort", 3306, "db port")
var dbName = flag.String("dbName", "dbnote", "db name")
var dbUser = flag.String("dbUser", "root", "db user")
var dbPass = flag.String("dbPass", "123456", "db password")

func main() {

	flag.Parse()

	db := dbhelper.GetDB(*dbIP, *dbPort, *dbName, *dbUser, *dbPass)
	xlsx, err := excelize.OpenFile("../../doc/job_qualityV2.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	keys := []string{}
	indexs := []int{}
	combinationKeyIndex := map[int]bool{}
	combinationValueIndex := -1
	combinationKeys := []string{}
	values := [][]interface{}{}

	rows := xlsx.GetRows("job_quality")
	for rowno, row := range rows {
		fmt.Println("rowno:", rowno)
		if 1 == rowno {
			for index, colCell := range row {
				if len(colCell) > 0 {

					if "key" == colCell {
						combinationKeyIndex[index] = true
					}

					if "value" == colCell {
						combinationValueIndex = index
					}

				}
			}
		}

		if 2 == rowno {
			for k, _ := range combinationKeyIndex {

				combinationKeys = append(combinationKeys, row[k])

			}
		}

		if 3 == rowno {

			for index, colCell := range row {
				if len(colCell) > 0 && !combinationKeyIndex[index] {
					keys = append(keys, colCell)
					indexs = append(indexs, index)
				}
			}
		}

		if rowno > 3 {
			values = append(values, []interface{}{})

			for _, index := range indexs {

				values[rowno-4] = append(values[rowno-4], row[index])

				if index == combinationValueIndex {

					combinationValue := ""

					for i := 0; i < len(combinationKeyIndex)-1; i++ {
						combinationValue += combinationKeys[i] + ":" + row[i] + ","
					}

					combinationValue += combinationKeys[len(combinationKeyIndex)-1] + ":" + row[len(combinationKeyIndex)-1]

					values[rowno-4][len(values[rowno-4])-1] = combinationValue
				}
			}

		}

	}

	sql := "insert into hero (" + strings.Join(keys, ",") + ") values (" + strings.Repeat("?,", len(keys)-1) + "?)"

	for _, value := range values {
		_, err := db.Exec(sql, value...)

		if nil != err {
			fmt.Println(err)
		}
	}

}
