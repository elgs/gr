package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elgs/gorest"
	"github.com/elgs/gosqljson"
	"io/ioutil"
	"time"
)

func init() {
	loadACL()
	gorest.RegisterGlobalDataInterceptor(&GlobalTokenInterceptor{Id: "GlobalTokenInterceptor"})
}

var acl = make(map[string]map[string]bool)
var tokenRegistry = make(map[string]map[string]string)

func checkToken(db *sql.DB, id string, key string, context map[string]interface{}, tableId string) (bool, error) {
	if id != "" && key != "" && len(tokenRegistry[id]) > 0 && tokenRegistry[id]["TOKEN_KEY"] == key {
		context["user_token"] = tokenRegistry[id]
		return true, nil
	}
	tokenTable := context["token_table"]
	if tokenTable != nil {
		tokenTableString := tokenTable.(string)
		if len(tokenTableString) > 0 && tableId == tokenTable {
			return false, errors.New("We think you are invading the system.")
		}

		gorest.MysqlSafe(&tokenTableString)
		data, err := gosqljson.QueryDbToMap(db, "upper", fmt.Sprint("SELECT * FROM ", tokenTableString, " WHERE ID=? AND TOKEN_KEY=?"), id, key)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		if data != nil && len(data) == 1 {
			record := data[0]
			tokenRegistry[record["ID"]] = record
			context["user_token"] = record
			return true, nil
		}
	}
	return false, errors.New("Authentication failed.")
}

func loadACL() {
	// load acl from configuration files.
	configFile := "gorest_acl.json"
	aclConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(configFile, " not found, default policies are used.")
	}
	err = json.Unmarshal(aclConfig, &acl)
	if err != nil {
		fmt.Println(err)
	}
}

func checkACL(tableId string, op string) (bool, error) {
	if acl[tableId] != nil {
		if acl[tableId][op] {
			return true, nil
		}
	}
	return false, errors.New("Access denied.")
}

type GlobalTokenInterceptor struct {
	*gorest.DefaultDataInterceptor
	Id string
}

func (this *GlobalTokenInterceptor) BeforeCreate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	if ok, err := checkACL(resourceId, "create"); !ok {
		return false, err
	}
	ctn, err := checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
	if ctn && err == nil {
		if context["meta"] != nil && context["meta"].(bool) {
			userToken := context["user_token"]
			if v, ok := userToken.(map[string]string); ok {
				data["CREATOR_ID"] = v["ID"]
				data["CREATOR_CODE"] = v["CODE"]
				data["CREATE_TIME"] = time.Now()
				data["UPDATER_ID"] = v["ID"]
				data["UPDATER_CODE"] = v["CODE"]
				data["UPDATE_TIME"] = time.Now()
			}
		}
	}
	return ctn, err
}
func (this *GlobalTokenInterceptor) AfterCreate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeLoad(resourceId string, db *sql.DB, field []string, context map[string]interface{}, id string) (bool, error) {
	if ok, err := checkACL(resourceId, "load"); !ok {
		return false, err
	}
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterLoad(resourceId string, db *sql.DB, field []string, context map[string]interface{}, data map[string]string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeUpdate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) (bool, error) {
	if ok, err := checkACL(resourceId, "update"); !ok {
		return false, err
	}
	ctn, err := checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
	if ctn && err == nil {
		if context["meta"] != nil && context["meta"].(bool) {
			userToken := context["user_token"]
			if v, ok := userToken.(map[string]string); ok {
				data["UPDATER_ID"] = v["ID"]
				data["UPDATER_CODE"] = v["CODE"]
				data["UPDATE_TIME"] = time.Now()
			}
		}
	}
	return ctn, err
}
func (this *GlobalTokenInterceptor) AfterUpdate(resourceId string, db *sql.DB, context map[string]interface{}, data map[string]interface{}) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeDuplicate(resourceId string, db *sql.DB, context map[string]interface{}, id string) (bool, error) {
	if ok, err := checkACL(resourceId, "duplicate"); !ok {
		return false, err
	}
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterDuplicate(resourceId string, db *sql.DB, context map[string]interface{}, id string, newId string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeDelete(resourceId string, db *sql.DB, context map[string]interface{}, id string) (bool, error) {
	if ok, err := checkACL(resourceId, "delete"); !ok {
		return false, err
	}
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterDelete(resourceId string, db *sql.DB, context map[string]interface{}, id string) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeListMap(resourceId string, db *sql.DB, field []string, context map[string]interface{}, filter *string, sort *string, group *string, start int64, limit int64, includeTotal bool) (bool, error) {
	if ok, err := checkACL(resourceId, "list"); !ok {
		return false, err
	}
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterListMap(resourceId string, db *sql.DB, field []string, context map[string]interface{}, data []map[string]string, total int64) error {
	return nil
}
func (this *GlobalTokenInterceptor) BeforeListArray(resourceId string, db *sql.DB, field []string, context map[string]interface{}, filter *string, sort *string, group *string, start int64, limit int64, includeTotal bool) (bool, error) {
	if ok, err := checkACL(resourceId, "list"); !ok {
		return false, err
	}
	return checkToken(db, context["api_token_id"].(string), context["api_token_key"].(string), context, resourceId)
}
func (this *GlobalTokenInterceptor) AfterListArray(resourceId string, db *sql.DB, field []string, context map[string]interface{}, data [][]string, total int64) error {
	return nil
}
