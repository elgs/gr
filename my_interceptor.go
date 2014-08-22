package main

import (
	"fmt"
	//"github.com/elgs/gorest"
	_ "github.com/go-sql-driver/mysql"
	"gorest"
)

func init1() {
	tableId := "test.TEST"
	gorest.RegisterDataInterceptor(tableId, &MyDataInterceptor{Id: "Local"})
	gorest.RegisterGlobalDataInterceptor(&MyDataInterceptor{Id: "Global"})
}

type MyDataInterceptor struct {
	*gorest.EchoDataInterceptor
	Id string
}

func (this *MyDataInterceptor) BeforeLoad(ds interface{}, context map[string]interface{}, id string) (bool, error) {
	fmt.Println(this.Id, ": Here I'm in BeforeLoad!")
	return true, nil
}
