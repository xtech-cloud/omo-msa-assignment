package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type ApplyService struct{}

func switchApply(info *cache.ApplyInfo) *pb.ApplyInfo {
	tmp := new(pb.ApplyInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator

	tmp.Type = uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Applicant = info.Applicant
	tmp.Inviter = info.Inviter
	tmp.Owner = info.Scene
	tmp.Group = info.Group
	tmp.Remark = info.Remark
	tmp.Reason = info.Reason
	return tmp
}

func (mine *ApplyService) AddOne(ctx context.Context, in *pb.ReqApplyAdd, out *pb.ReplyApplyOne) error {
	path := "apply.add"
	inLog(path, in)
	if len(in.Applicant) < 1 {
		out.Status = outError(path, "the application is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateApply(in.Operator, in.Owner, in.Group, in.Applicant, in.Inviter, in.Remark, uint8(in.Type))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchApply(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ApplyService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyApplyOne) error {
	path := "apply.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetApply(in.Uid)
	if er != nil {
		out.Status = outError(path, "the apply not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchApply(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ApplyService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "apply.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ApplyService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "apply.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveApply(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *ApplyService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyApplyList) error {
	path := "apply.getListByFilter"
	inLog(path, in)
	var list []*cache.ApplyInfo
	var err error
	if in.Key == "" {
		list = cache.Context().GetAppliesByOwner(in.Owner, -1)
	} else if in.Key == "creator" {
		list = cache.Context().GetAppliesByCreator(in.Value)
	} else if in.Key == "application" {
		list = cache.Context().GetAppliesByUser(in.Value)
	} else if in.Key == "group" {
		list = cache.Context().GetAppliesByGroup(in.Value)
	} else if in.Key == "scene" {
		tp := parseStringToInt(in.Value)
		list = cache.Context().GetAppliesByOwner(in.Owner, int32(tp))
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ApplyInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchApply(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ApplyService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "apply.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var err error
	//info, er := cache.Context().GetApply(in.Uid)
	//if er != nil {
	//	out.Status = outError(path, "the apply not found ", pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	//
	//if in.Key == "master" {
	//
	//} else if in.Key == "assistants" {
	//
	//}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ApplyService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "apply.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetApply(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.SetStatus(uint8(in.Flag), in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
