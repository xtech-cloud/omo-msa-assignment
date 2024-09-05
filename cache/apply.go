package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

const (
	/*
		待审核
	*/
	ApplyStatusPending = 0
	/*
		通过
	*/
	ApplyStatusPass = 1
	/*
		被拒
	*/
	ApplyStatusRefused = 2

	/*
		不存在
	*/
	ApplyStatusNull = 5

	/**
	搁置
	**/
	ApplyStatusStay = 6
)

type ApplyInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	SubmitTime time.Time
	/**
	邀请人
	*/
	Inviter string
	/**
	申请人，被邀请人
	*/
	Applicant string
	/**
	所属场景
	*/
	Scene string
	/**
	所属部门或者小组
	*/
	Group  string
	Reason string
	Remark string
}

func (mine *cacheContext) GetAppliesByUser(uid string) []*ApplyInfo {
	list := make([]*ApplyInfo, 0, 5)
	array, err := nosql.GetAppliesByApplicant(uid)
	if err == nil {
		for _, item := range array {
			info := new(ApplyInfo)
			info.initInfo(item)

			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetAppliesByCreator(uid string) []*ApplyInfo {
	list := make([]*ApplyInfo, 0, 5)
	array, err := nosql.GetAppliesByCreator(uid)
	if err == nil {
		for _, item := range array {
			info := new(ApplyInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetAppliesByGroup(uid string) []*ApplyInfo {
	list := make([]*ApplyInfo, 0, 5)
	array, err := nosql.GetAppliesByGroup(uid)
	if err == nil {
		for _, item := range array {
			info := new(ApplyInfo)
			info.initInfo(item)

			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetAppliesByOwner(scene string, tp int32) []*ApplyInfo {
	var array []*nosql.Apply
	var err error
	if tp < 0 {
		array, err = nosql.GetAppliesByScene1(scene)
	} else {
		array, err = nosql.GetAppliesByScene(scene, uint8(tp))
	}

	list := make([]*ApplyInfo, 0, len(array))
	if err == nil {
		for _, item := range array {
			info := new(ApplyInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetApply(uid string) (*ApplyInfo, error) {
	db, err := nosql.GetApply(uid)
	if err != nil {
		return nil, err
	}
	info := new(ApplyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveApply(uid, operator string) error {
	err := nosql.RemoveApply(uid, operator)
	if err != nil {
		return err
	}
	return nil
}

func (mine *cacheContext) CreateApply(creator, scene, group, applicant, inviter, remark string, tp uint8) (*ApplyInfo, error) {
	var db = new(nosql.Apply)
	db.UID = primitive.NewObjectID()
	db.CreatedTime = time.Now()
	db.ID = nosql.GetApplyNextID()
	db.Creator = creator
	db.Applicant = applicant
	db.Inviter = inviter
	db.SubmitTime = time.Now()
	db.Scene = scene
	db.Group = group
	db.Type = tp
	db.Remark = remark
	db.Status = ApplyStatusPending

	info := new(ApplyInfo)
	info.initInfo(db)
	err := nosql.CreateApply(db)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (mine *ApplyInfo) initInfo(db *nosql.Apply) bool {
	if db == nil {
		return false
	}
	mine.Scene = db.Scene
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Type = db.Type
	mine.Inviter = db.Inviter
	mine.SubmitTime = db.SubmitTime
	mine.Applicant = db.Applicant
	mine.Status = db.Status
	mine.Group = db.Group
	return true
}

func (mine *ApplyInfo) SetStatus(dist uint8, reason, operator string) error {
	if mine.Status != ApplyStatusPending {
		return errors.New("the apply now status not pending")
	}
	if dist == ApplyStatusPending {
		return errors.New("the apply dist status is pending")
	}
	err := nosql.UpdateApply(mine.UID, reason, operator, dist)
	if err == nil {
		mine.Status = dist
	}
	return err
}
