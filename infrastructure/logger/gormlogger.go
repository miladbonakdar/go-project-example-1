package logger

import "fmt"

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {

	case "sql":
		WithData(
			map[string]interface{}{
				"module":  "gorm",
				"type":    "sql",
				"rows":    v[5],
				"src_ref": v[1],
				"values":  v[4],
			},
		).Debug(fmt.Sprintf("%v", v[3]))
	case "log":
		WithData(map[string]interface{}{"module": "gorm", "type": "log"}).Debug(fmt.Sprintf("%v", v[2]))
	case "error":
		WithData(map[string]interface{}{"module": "gorm", "type": "error"}).Error(fmt.Sprintf("%v", v[2]))
	}
}
