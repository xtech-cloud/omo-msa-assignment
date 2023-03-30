package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type MeetingService struct{}

func switchMeeting(info *cache.MeetingInfo) *pb.MeetingInfo {
	tmp := new(pb.MeetingInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Status = uint32(info.CheckStatus())
	tmp.Type = uint32(info.Type)
	tmp.Owner = info.Owner
	tmp.Group = info.Group
	tmp.Stopped = info.StopTime.Unix()
	tmp.Started = info.StartTime.Unix()
	tmp.Duration = uint32(info.Duration)
	tmp.Appointed = info.Appointed
	tmp.Location = info.Location
	tmp.Signs = info.Signs
	tmp.Submits = info.Submits
	tmp.Notifies = info.Notifies
	return tmp
}

func (mine *MeetingService) AddOne(ctx context.Context, in *pb.ReqMeetingAdd, out *pb.ReplyMeetingOne) error {
	path := "meeting.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateMeeting(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchMeeting(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMeetingOne) error {
	path := "meeting.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetMeeting(in.Uid)
	if er != nil {
		out.Status = outError(path, "the meeting not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchMeeting(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "meeting.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "meeting.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveMeeting(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyMeetingList) error {
	path := "meeting.getListByFilter"
	inLog(path, in)
	var list []*cache.MeetingInfo
	var err error
	if in.Key == "" {
		list = cache.Context().GetMeetingsByOwner(in.Owner)
	} else if in.Key == "type" {

	} else if in.Key == "group" {
		list = cache.Context().GetMeetingsByGroup(in.Value)
	} else if in.Key == "time" {
		if len(in.Values) > 1 {
			list, err = cache.Context().GetMeetingsByTime(in.Owner, in.Values[0], in.Values[1])
		} else {
			err = errors.New("the params is error")
		}
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.MeetingInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchMeeting(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *MeetingService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "meeting.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetMeeting(in.Uid)
	if er != nil {
		out.Status = outError(path, "the meeting not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "submit" {
		err = info.Submit(in.Value, in.Operator)
	} else if in.Key == "location" {
		err = info.UpdateLocation(in.Value, in.Operator, info.Type)
	} else if in.Key == "stop" {
		err = info.UpdateStop(in.Value, in.Operator)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "meeting.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetMeeting(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Flag == int32(cache.Close) {
		err = info.Close(in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) UpdateBase(ctx context.Context, in *pb.ReqMeetingUpdate, out *pb.ReplyInfo) error {
	path := "meeting.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetMeeting(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *MeetingService) Sign(ctx context.Context, in *pb.ReqMeetingSign, out *pb.ReplyInfo) error {
	path := "meeting.sign"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetMeeting(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.Sign(in.Member, in.Operator, in.Location)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
