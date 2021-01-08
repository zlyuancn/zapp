/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/3
   Description :
-------------------------------------------------
*/

package mysql_binlog

import (
	jsoniter "github.com/json-iterator/go"
)

type Record struct {
	Action    string                 `json:"action"`
	Old       map[string]interface{} `json:"old"`
	New       map[string]interface{} `json:"new"`
	DbName    string                 `json:"db_name"`
	TableName string                 `json:"table_name"`
	Timestamp uint32                 `json:"timestamp"`
}

func (r *Record) OldString() string {
	text, _ := jsoniter.MarshalToString(r.Old)
	return text
}
func (r *Record) NewString() string {
	text, _ := jsoniter.MarshalToString(r.New)
	return text
}
func (r *Record) String() string {
	text, _ := jsoniter.MarshalToString(r)
	return text
}
