package main

import (
	"fmt"
	//"github.com/elgs/gorest"
	//"github.com/elgs/gosqljson"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"gorest"
	"gosqljson"
)

func init() {
	gorest.RegisterGlobalDataInterceptor(&GlobalTokenInterceptor{Id: "GlobalTokenInterceptor"})
}

func checkToken(db *sql.DB, id string, key string, context map[string]interface{}) (bool, error) {
	if context["table_id"] == "gorest.token" {
		return false, errors.New("We think you are invading the system.")
	}
	if id != "" && key != "" && tokenRegistry[id] == key {
		return true, nil
	}
	data, err := gosqljson.QueryDbToMap(db, false, "SELECT * FROM gorest.token WHERE ID=? AND TOKEN_KEY=? AND STATUS=?", id, key, "0")
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if data != nil && len(data) == 1 {
		tokenRegistry[data[0]["ID"]] = data[0]["TOKEN_KEY"]
		return true, nil
	}
	return false, errors.New("Authentication failed.")
}

var tokenRegistry = make(map[string]string)

type GlobalTokenInterceptor struct {
	*gorest.DefaultDataInterceptor
	Id string
}

func (this *GlobalTokenInterceptor) BeforeCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterCreate(ds interface{}, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeLoad(ds interface{}, context map[string]interface{}, id string) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterLoad(ds interface{}, context map[string]interface{}, data map[string]string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeUpdate(ds interface{}, context map[string]interface{}, oldData map[string]interface{}, data map[string]interface{}) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterUpdate(ds interface{}, context map[string]interface{}, oldData map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeDuplicate(ds interface{}, context map[string]interface{}, oldData map[string]interface{}, data map[string]interface{}) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterDuplicate(ds interface{}, context map[string]interface{}, oldData map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeDelete(ds interface{}, context map[string]interface{}, id string) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterDelete(ds interface{}, context map[string]interface{}, id string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeListMap(ds interface{}, context map[string]interface{}, where string, order string, start int64, limit int64, includeTotal bool) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterListMap(ds interface{}, context map[string]interface{}, data []map[string]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeListArray(ds interface{}, context map[string]interface{}, where string, order string, start int64, limit int64, includeTotal bool) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterListArray(ds interface{}, context map[string]interface{}, data [][]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeQueryMap(ds interface{}, context map[string]interface{}, sqlSelect string, sqlSelectCount string, start int64, limit int64, includeTotal bool) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterQueryMap(ds interface{}, context map[string]interface{}, data []map[string]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeQueryArray(ds interface{}, context map[string]interface{}, sqlSelect string, sqlSelectCount string, start int64, limit int64, includeTotal bool) (bool, error) {
	db := ds.(*sql.DB)
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context)
}
func (this *GlobalTokenInterceptor) AfterQueryArray(ds interface{}, context map[string]interface{}, data [][]string, total int64) error {
	return nil
}
