package cache

import (
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"omo.msa.assignment/config"
	"omo.msa.assignment/proxy/nosql"
	"reflect"
	"time"
)

type baseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}
	//dbs, _ := nosql.GetAllFamilies()
	//for _, db := range dbs {
	//	fmt.Printf(db.Name)
	//}
	return nil
}

func Context() *cacheContext {
	return cacheCtx
}

func checkPage(page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total / number
	if total%number != 0 {
		maxPage = total/number + 1
	}
	if page < 1 {
		return total, maxPage, all
	}
	if page > maxPage {
		page = maxPage
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}

func switchOldFamilyToCoterie() {
	dbs, _ := nosql.GetAllFamilies()
	for _, db := range dbs {
		if len(db.Children) > 0 {
			in := new(pb.ReqCoterieAdd)
			in.Centre = db.Children[0]
			in.Name = db.Name
			in.Remark = db.Remark
			in.Type = 0
			in.Passwords = db.Passwords
			in.Operator = db.Creator
			in.Cover = db.Cover
			in.Master = db.Master
			in.Members = make([]*pb.IdentifyInfo, 0, len(db.Custodians))
			for _, custodian := range db.Custodians {
				if len(custodian.Identifies) > 0 {
					in.Members = append(in.Members, &pb.IdentifyInfo{User: custodian.User, Name: "", Remark: custodian.Identifies[0].Remark})
				}else{
					in.Members = append(in.Members, &pb.IdentifyInfo{User: custodian.User, Name: "", Remark: ""})
				}

			}
			cacheCtx.CreateCoterie(in)
		}
	}
}