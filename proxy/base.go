package proxy

import (
	"time"
)

type DateInfo struct {
	Begin string `json:"begin" bson:"begin"`
	End   string `json:"end" bson:"end"`
}

type RecordInfo struct {
	UID         string    `json:"uid" bson:"uid"`
	CreatedTime time.Time `json:"createdAt" bson:"createdAt"`
	Creator     string    `json:"creator" bson:"creator"`

	Name   string   `json:"name" bson:"name"`
	Remark string   `json:"remark" bson:"remark"`
	Executor  string   `json:"executor" bson:"executor"`
	Status uint8 `json:"status" bson:"status"`
	Tags   []string `json:"tags" bson:"tags"`
	Assets []string `json:"assets" bson:"assets"`
}

type CustodianInfo struct {
	User   string `json:"user" bson:"user"`
	Identifies []IdentifyInfo `json:"identify" bson:"identify"`
}

type IdentifyInfo struct {
	Child   string `json:"child" bson:"child"`
	Remark string `json:"remark" bson:"remark"`
}

type MemberInfo struct {
	User   string `json:"user" bson:"user"`
	Remark string `json:"remark" bson:"remark"`
}
