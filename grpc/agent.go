package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
	"strconv"
)

type AgentService struct {}

func switchAgent(info *cache.AgentInfo) *pb.AgentInfo {
	tmp := new(pb.AgentInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Owner = info.Scene
	tmp.Remark = info.Remark
	tmp.Quotes = info.Quotes
	tmp.Devices = info.Products()
	return tmp
}

func (mine *AgentService)AddOne(ctx context.Context, in *pb.ReqAgentAdd, out *pb.ReplyAgentInfo) error {
	path := "Agent.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Owner)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	if scene.HadAgentByName(in.Name) {
		out.Status = outError(path,"the name repeated ", pbstatus.ResultStatus_Repeated)
		return nil
	}

	Agent, err := scene.CreateAgent(in)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchAgent(Agent)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAgentInfo) error {
	path := "Agent.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.AgentInfo
	if len(in.Parent) > 0 {
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info = scene.GetAgent(in.Uid)
	}else{
		info = cache.Context().GetAgent(in.Uid)
	}

	if info == nil {
		out.Status = outError(path,"the Agent not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchAgent(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Agent.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path,"the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "Agent.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveAgent(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService)GetList(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyAgentList) error {
	path := "Agent.getList"
	inLog(path, in)
	var list []*cache.AgentInfo
	if in.Parent == "" {
		if in.Key == "device" {
			list = cache.Context().GetAgentsByDevice(in.Value)
		}else if in.Key == "quote" {
			list = cache.Context().GetAgentsByQuote(in.Value)
		}else{
			list = make([]*cache.AgentInfo, 0, 1)
		}
	}else{
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Key == "" {
			list = scene.GetAgents()
		}else if in.Key == "product" {
			tp,er := strconv.ParseUint(in.Value, 10, 32)
			if er != nil {
				out.Status = outError(path,er.Error(), pbstatus.ResultStatus_DBException)
				return nil
			}
			list = scene.GetAgentsByType(uint8(tp))
		}else if in.Key == "quote" {
			list = scene.GetAgentsByQuote(in.Value)
		}else if in.Key == "device" {
			list = scene.GetAgentsByDevice(in.Value)
		}else{
			list = make([]*cache.AgentInfo, 0, 1)
		}
	}


	out.List = make([]*pb.AgentInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchAgent(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AgentService) UpdateBase (ctx context.Context, in *pb.ReqAgentUpdate, out *pb.ReplyInfo) error {
	path := "Agent.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAgent(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Agent not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if len(in.Name) > 0 || len(in.Remark) > 0 {
		scene := cache.Context().GetScene(info.Scene)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Name != info.Name && scene.HadAgentByName(in.Name) {
			out.Status = outError(path,"the Agent name repeated ", pbstatus.ResultStatus_Repeated)
			return nil
		}
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) UpdateQuotes (ctx context.Context, in *pb.ReqAgentQuotes, out *pb.ReplyInfo) error {
	path := "Agent.updateQuotes"
	inLog(path, in)
	if len(in.Scene) < 1 || len(in.Agent) < 1 {
		out.Status = outError(path,"the scene or Agent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	info := scene.GetAgent(in.Agent)
	if info == nil {
		out.Status = outError(path,"the Agent not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	scene.ClearQuotes(in.Operator, in.List)
	err := info.UpdateQuotes(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) UpdateDisplays (ctx context.Context, in *pb.ReqAgentDisplays, out *pb.ReplyInfo) error {
	path := "Agent.updateDisplays"
	inLog(path, in)
	if len(in.Scene) < 1 {
		out.Status = outError(path,"the scene or Agent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"the scene not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	for _, Agent := range in.Agents {
		info := scene.GetAgent(Agent.Agent)
		if info == nil {
			out.Status = outError(path,"the Agent not found ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		err := info.UpdateDisplays(Agent.Sn, Agent.Group, in.Operator, Agent.Showing, Agent.List)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) AppendDevice (ctx context.Context, in *pb.ReqAgentDevice, out *pb.ReplyAgentDevices) error {
	path := "Agent.appendDevice"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	//if cache.Context().HadBindDeviceInAgent(in.Device) {
	//	out.Status = outError(path,"the device had bind by other Agent", pbstatus.ResultStatus_Repeated)
	//	return nil
	//}
	info := cache.Context().GetAgent(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Agent not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if info.HadDeviceByType(uint8(in.Type)) {
		out.Status = outError(path,"the Agent had the device of type", pbstatus.ResultStatus_Repeated)
		return nil
	}

	err := info.AppendDevice(in.Device, in.Remark, in.Operator, in.Type)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Products()
	out.Status = outLog(path, out)
	return nil
}

func (mine *AgentService) SubtractDevice (ctx context.Context, in *pb.ReqAgentDevice, out *pb.ReplyAgentDevices) error {
	path := "Agent.subtractDevice"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAgent(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Agent not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	device,err := cache.Context().GetDevice(in.Device)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	err = device.UpdateAgent("", in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Products()
	out.Status = outLog(path, out)
	return nil
}


