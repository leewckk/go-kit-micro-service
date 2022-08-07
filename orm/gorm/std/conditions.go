package std

import (
	"fmt"
	"strings"
	"time"
)

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

func DateFormat(t time.Time) string {
	dateFmt := "2006-01-02"
	return fmt.Sprintf("%v", t.Format(dateFmt))
}

// sql build where
func WhereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 3 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)

		}

		if whereSQL != "" {
			whereSQL += " AND "

		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")

				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")

				}
			default:
				whereSQL += fmt.Sprint(k, " = ?")
				vals = append(vals, v)

			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, " = ?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, " > ?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, " >= ?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, " < ?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, " <= ?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, " != ?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, " != ?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)

			}
			break
		case 3:
			k1 := ks[0]
			switch ks[1] {
			case "not":
				whereSQL += fmt.Sprintf("%v %v %v (?)", k1, ks[1], ks[2])
				vals = append(vals, v)
			}
			break
		}
	}
	return
}
