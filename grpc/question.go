package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
	"strconv"
)

type QuestionService struct{}

func switchQuestion(info *cache.QuestionInfo) *pb.QuestionInfo {
	tmp := new(pb.QuestionInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cd = uint32(info.Cd)
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Quote = info.Quote
	tmp.Answers = info.Answers
	tmp.Assets = info.Assets
	tmp.Options = make([]*pb.QuestionOption, 0, len(info.Options))
	for _, option := range info.Options {
		id, _ := strconv.Atoi(option.Key)
		tmp.Options = append(tmp.Options, &pb.QuestionOption{
			Id:   uint32(id),
			Desc: option.Value,
		})
	}
	return tmp
}

func (mine *QuestionService) AddOne(ctx context.Context, in *pb.ReqQuestionAdd, out *pb.ReplyQuestionOne) error {
	path := "question.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().NewQuestion(in.Name, in.Remark, in.Category, in.Quote, in.Operator, int(in.Cd), in.Answers, in.Options)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchQuestion(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *QuestionService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyQuestionOne) error {
	path := "question.getOne"
	inLog(path, in)

	var info *cache.QuestionInfo
	var err error
	info, err = cache.Context().GetQuestion(in.Uid)

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchQuestion(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *QuestionService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "question.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *QuestionService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "question.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetQuestion(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	err = info.Delete(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *QuestionService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyQuestionList) error {
	path := "question.getListByFilter"
	inLog(path, in)
	var list []*cache.QuestionInfo
	var err error
	if in.Key == "entity" {
		list, err = cache.Context().GetQuestionsByEntity(in.Value)
	} else if in.Key == "name" {
		list, err = cache.Context().GetQuestionsByName(in.Value)
	} else if in.Key == "category" {
		list, err = cache.Context().GetQuestionsByCategory(in.Value)
	} else if in.Key == "name_kind" {
		list, err = cache.Context().GetQuestionsByNameAndKind(in.Value, in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.List = make([]*pb.QuestionInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchQuestion(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *QuestionService) UpdateBase(ctx context.Context, in *pb.ReqQuestionUpdate, out *pb.ReplyQuestionOne) error {
	path := "question.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetQuestion(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Operator, in.Category, uint16(in.Cd), in.Answers, in.Options)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchQuestion(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *QuestionService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "question.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetQuestion(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "assets" {
		err = info.UpdateAssets(in.Operator, in.Values)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
