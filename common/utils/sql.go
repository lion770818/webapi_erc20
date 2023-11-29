package utils

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// SQLPagination 增加 sql limit, offset 條件
func SQLPagination(db *gorm.DB, rowCount, offset int) *gorm.DB {
	return db.Limit(rowCount).Offset(offset)
}

// SQLAppend 依照 condition 判斷是否拼湊 SQL Where 條件
func SQLAppend(db *gorm.DB, condition bool, query, columnValue interface{}) *gorm.DB {
	if !condition {
		return db
	}

	return db.Where(query, columnValue)
}

// SQLRawAppend 依照 condition 判斷是否拼湊 SQL Where 條件，給組 Raw SQL 場景
func SQLRawAppend(condition bool, sql *strings.Builder, query string, values []interface{}, value interface{}) []interface{} {
	if !condition {
		return values
	}

	sql.WriteString(query)

	strs, ok := value.([]string)
	if !ok {
		return append(values, value)
	}

	for _, v := range strs {
		values = append(values, v)
	}

	return values
}

// SQLArrayToString 字串陣列轉為 SQL AND (column OR column) 語句
func SQLArrayToString(strs []string, column string) string {
	if len(strs) == 0 {
		return ""
	}

	var sql strings.Builder

	sql.WriteString(" AND (")

	for i := 0; i < len(strs); i++ {
		if i != 0 {
			sql.WriteString(" OR ")
		}

		sql.WriteString(column + "= ?")
	}

	sql.WriteString(")")

	return sql.String()
}

func Assign(origin, target interface{}, excludes ...string) {
	valOrigin := reflect.ValueOf(origin).Elem()
	valTarget := reflect.ValueOf(target).Elem()

	for i := 0; i < valOrigin.NumField(); i++ {
		if !valTarget.FieldByName(valOrigin.Type().Field(i).Name).IsValid() {
			continue
		}

		isExclude := false
		for _, col := range excludes {
			if valOrigin.Type().Field(i).Name == col {
				isExclude = true
				break
			}
		}

		if isExclude {
			continue
		}

		tmpOrigin := valOrigin.Field(i)
		tmpTarget := valTarget.FieldByName(valOrigin.Type().Field(i).Name)
		if reflect.TypeOf(tmpOrigin.Interface()) != reflect.TypeOf(tmpTarget.Interface()) {
			continue
		}

		tmpTarget.Set(tmpOrigin)
	}
}
