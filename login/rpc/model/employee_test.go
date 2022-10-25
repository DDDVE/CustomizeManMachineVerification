package model

import (
	"testing"

	"github.com/jinzhu/gorm"
)

func TestInsert(t *testing.T) {
	d, _ := gorm.Open("mysql", "root:123456@tcp(43.143.208.232:3306)/cmmvplat?charset=utf8&parseTime=true&loc=Local")
	defer d.Close()
	var g DefaultEmployeeModel
	g.db = d
	g.Insert(&Employee{MobileNum: "13139200815"})
}
