package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// func FieldsConstructor(q interface{}) []string {
// 	t := reflect.TypeOf(q)
// 	i := 0
// 	var res []string
// 	shouldAppendField := func(field reflect.StructField) (string, bool) {
// 		dbTag := t.Field(i).Tag.Get("json")
// 		if dbTag == "" {
// 			return "", false
// 		}
// 		return dbTag, true
// 	}
// 	if reflect.ValueOf(q).Kind() == reflect.Struct {
// 		for ; i < t.NumField()-1; i++ {
// 			if dbTag, should := shouldAppendField(t.Field(i)); should {
// 				res = append(res, dbTag)
// 			}
// 		}
// 	}
// 	if dbTag, should := shouldAppendField(t.Field(i)); should {
// 		res = append(res, dbTag)
// 	}
// 	return res
// }

func GetDate() string {
	return time.Now().Format(time.RFC3339)
}

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02 15:04:11")
}

func SetColumnPrefix(prefix string, column string) string {
	return fmt.Sprint(prefix, ".", strings.ToLower(column))
}

// return columns in the format of "prefix[0].column as prefix.column"
func Prefixer(prefix string, columns []string) []string {
	ans := make([]string, len(columns))
	for i, v := range columns {
		ans[i] = fmt.Sprint(string(prefix[0]), ".", v, " as '", SetColumnPrefix(prefix, v), "'")
	}
	return ans
}

func RollbackOnErr(err error, conn *sqlx.Tx) {
	if err != nil {
		conn.Rollback()
	}
}

// func Deprefix(prefix string) *reflectx.Mapper {
// 	return reflectx.NewMapperTagFunc("db", func(s string) string {
// 		return SetColumnPrefix(prefix, s)
// 	}, func(tag string) string {
// 		return SetColumnPrefix(prefix, tag)
// 	})
// }

// type RowStructMap map[string]interface{}

// func MultiStructScan(source map[string]interface{}, target RowStructMap) error {
// 	for prefix, v := range target {
// 		refV := reflect.ValueOf(v).Elem()
// 		if refV.Kind() != reflect.Struct {
// 			return errors.New("not struct")
// 		}
// 		t := reflect.TypeOf(v).Elem()
// 		for i := 0; i < t.NumField(); i++ {
// 			f := t.Field(i)
// 			vf := refV.FieldByName(f.Name)
// 			if !vf.CanSet() {
// 				return errors.New("can not set value")
// 			}
// 			// get column name
// 			columnName := ""
// 			dbTag := f.Tag.Get("db")
// 			if dbTag == "-" {
// 				// is ignored
// 				continue
// 			}
// 			if dbTag == "" {
// 				// there is no db tag, use field name
// 				columnName = strings.ToLower(f.Name)
// 			} else {
// 				columnName = dbTag
// 			}
// 			// add prefix to column name
// 			columnName = SetColumnPrefix(prefix, columnName)
// 			sourceV, ok := source[columnName]
// 			if !ok {
// 				// TODO error
// 				continue
// 			}
// 			refSourceV := reflect.ValueOf(sourceV)
// 			if vf.Kind() != refSourceV.Kind() {
// 				return errors.New("type invalid")
// 			}
// 			vf.Set(refSourceV)
// 		}
// 	}
// 	return nil
// }
