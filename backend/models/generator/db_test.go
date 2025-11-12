package generator

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
	"testing"
)

const URL = "admin:cjwy837..@(8.130.14.218:3306)/phm?charset=utf8mb4&parseTime=True&loc=Local"

func TestDb(t *testing.T) {
	g := gen.NewGenerator(gen.Config{
		OutPath:        "./test/query",
		Mode:           gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode gen.WithoutContext |
		FieldNullable:  true,
		FieldCoverable: false,
	})

	gormdb, _ := gorm.Open(mysql.Open(URL))
	g.UseDB(gormdb) // reuse your gorm db
	// data type for protobuf
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	}
	//
	g.WithDataTypeMap(dataMap)
	// 将数字类型转为string
	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		toStringField := `balance, `
		if strings.Contains(toStringField, columnName) {
			return columnName + ",string"
		}
		return columnName
	})
	autoUpdateTimeField := gen.FieldJSONTag("update_date_time", "column:update_date_time;type:int unsigned;autoUpdateTime")
	autoCreateTimeField := gen.FieldJSONTag("create_date_time", "column:create_date_time;type:int unsigned;autoCreateTime")
	softDeleteField := gen.FieldJSONTag("is_soft_delete", "soft_delete.DeleteAt")
	fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}
	//Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.ApplyBasic(model.User{})

	IspUser := g.GenerateModel("sys_user")
	allModel := g.GenerateAllTable(fieldOpts...)
	g.ApplyBasic(IspUser)
	g.ApplyBasic(allModel...)
	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`

	// Generate the code
	g.Execute()
	t.Log("db 生成成功")
}
