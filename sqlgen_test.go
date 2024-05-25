package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestSqlGen(t *testing.T) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "dal/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
	})

	dsn := "root:root@tcp(127.0.0.1:3306)/api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	g.UseDB(db)
	g.ApplyBasic(g.GenerateModel("interface_info"), g.GenerateModel("user"))
	g.Execute()
}
