package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type TeamService struct{}

func switchTeam(info *cache.TeamInfo) *pb.TeamInfo {
	tmp := new(pb.TeamInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Master = info.Master
	tmp.Status = uint32(info.Status)
	tmp.Region = info.Region
	tmp.Assistants = info.Assistants
	tmp.Owner = info.Owner
	tmp.Members = info.Members
	tmp.Tags = info.Tags
	return tmp
}

func (mine *TeamService) AddOne(ctx context.Context, in *pb.ReqTeamAdd, out *pb.ReplyTeamInfo) error {
	path := "team.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	if cache.Context().HadTeamByName(in.Owner, in.Name) {
		out.Status = outError(path, "the name is repeated in scene ", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info, err := cache.Context().CreateTeam(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTeam(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeamInfo) error {
	path := "team.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTeam(in.Uid)
	if er != nil {
		out.Status = outError(path, "the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchTeam(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Team.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "team.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveTeam(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeamList) error {
	path := "team.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeamService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyTeamList) error {
	path := "team.getListByFilter"
	inLog(path, in)
	var list []*cache.TeamInfo
	var err error
	if in.Key == "" {
		list = cache.Context().GetTeamsByOwner(in.Owner)
	} else if in.Key == "type" {

	} else if in.Key == "user" {
		list = cache.Context().GetTeamsByUser(in.Value)
	} else if in.Key == "array" {
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeamInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchTeam(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeamService) UpdateBase(ctx context.Context, in *pb.ReqTeamUpdate, out *pb.ReplyInfo) error {
	path := "team.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTeam(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		if in.Name != info.Name && cache.Context().HadTeamByName(info.Owner, in.Name) {
			out.Status = outError(path, "the team name repeated ", pbstatus.ResultStatus_Repeated)
			return nil
		}
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "team.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTeam(in.Uid)
	if er != nil {
		out.Status = outError(path, "the team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "master" {
		err = info.UpdateMaster(in.Value, in.Operator)
	}else if in.Key == "assistants" {
		err = info.UpdateAssistants(in.Operator, in.Values)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "Team.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTeam(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(in.Operator, uint8(in.Flag))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) AppendMember(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "team.appendMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTeam(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMembers(in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Members
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) SubtractMember(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "Team.subtractMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er:= cache.Context().GetTeam(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMembers(in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Members
	out.Status = outLog(path, out)
	return nil
}
