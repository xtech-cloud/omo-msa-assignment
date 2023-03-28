package cache

import (
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"math"
	"omo.msa.assignment/config"
	"omo.msa.assignment/proxy/nosql"
	"reflect"
	"strconv"
	"strings"
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

type Vector struct {
	X float64 // 纬度latitude
	Y float64 // 经度longitude
	Z float64
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
				} else {
					in.Members = append(in.Members, &pb.IdentifyInfo{User: custodian.User, Name: "", Remark: ""})
				}

			}
			cacheCtx.CreateCoterie(in)
		}
	}
}

func (mine *cacheContext) checkDistance(center, loc string) bool {
	if len(center) < 1 {
		return true
	}
	if len(loc) < 1 {
		return false
	}
	from := parseLocation(loc)
	to := parseLocation(center)
	dis := geoDistance(from, to, "M")
	if dis < 100 {
		return true
	} else {
		return false
	}
}

func (mine *cacheContext) formatTime(from string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04", from, time.Local)
	if err == nil {
		return t, nil
	} else {
		return time.Now(), err
	}
}

func (mine *cacheContext) formatDate(from string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", from, time.Local)
	if err == nil {
		return t, nil
	} else {
		return time.Now(), err
	}
}

// GeoDistance 计算地理距离，依次为两个坐标的纬度、经度、单位（默认：英里，K => 公里，N => 海里）
func geoDistance(from, to Vector, unit string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := PI * from.X / 180
	radlat2 := PI * to.X / 180

	theta := to.Y - from.Y
	radTheta := PI * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radTheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if unit == "K" {
		dist = dist * 1.609344
	} else if unit == "N" {
		dist = dist * 0.8684
	} else if unit == "M" {

	}

	return dist
}

func parseLocation(local string) Vector {
	if len(local) < 1 {
		return Vector{X: 0, Y: 0, Z: 0}
	}
	arr := strings.Split(local, "|")
	if arr == nil || len(arr) < 3 {
		return Vector{X: 0, Y: 0, Z: 0}
	}
	x, _ := strconv.ParseFloat(arr[2], 64)
	y, _ := strconv.ParseFloat(arr[1], 64)
	return Vector{X: x, Y: y, Z: 0}
}
