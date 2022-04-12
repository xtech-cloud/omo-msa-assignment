package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type AgentService struct{}

func switchAgent(info *cache.AgentInfo) *pb.AgentInfo {
	tmp := new(pb.AgentInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Owner = info.Owner
	tmp.Remark = info.Remark
	tmp.Way = info.Way
	tmp.User = info.User
	tmp.Entity = info.Entity
	tmp.Status = uint32(info.Status)
	tmp.Type = uint32(info.Type)
	tmp.Regions = info.Regions
	tmp.Tags = info.Tags
	return tmp
}

func (mine *AgentService) AddOne(ctx context.Context, in *pb.ReqAgentAdd, out *pb.ReplyAgentOne) error {
	path := "agent.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateAgent(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchAgent(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAgentOne) error {
	path := "agent.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,err := cache.Context().GetAgent(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchAgent(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAgentList) error {
	path := "agent.search"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Agent.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "agent.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveAgent(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyAgentList) error {
	path := "agent.getListByFilter"
	inLog(path, in)
	var list []*cache.AgentInfo
	if in.Key == "" {
		list = cache.Context().GetAgentsByOwner(in.Owner)
	} else if in.Key == "way" {
		list = cache.Context().GetAgentsByWay(in.Owner, in.Value)
	} else {
		list = make([]*cache.AgentInfo, 0, 1)
	}

	out.List = make([]*pb.AgentInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchAgent(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AgentService) UpdateBase(ctx context.Context, in *pb.ReqAgentUpdate, out *pb.ReplyInfo) error {
	path := "agent.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetAgent(in.Uid)
	if er == nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "agent.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetAgent(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if in.Key == "entity" {
		err = info.UpdateEntity(in.Value, in.Operator)
	}else if in.Key == "tags" {
		err = info.UpdateTags(in.Operator, in.Values)
	}else if in.Key == "regions" {
		err = info.UpdateRegions(in.Operator, in.Values)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "agent.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetAgent(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(in.Operator, uint32(in.Flag))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
