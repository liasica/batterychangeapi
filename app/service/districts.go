package service

import (
	"battery/app/dao"
	"battery/app/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"io/ioutil"
)

var DistrictsService = districtsService{}

type districtsService struct {
}

// Child 获取下级地区
func (s *districtsService) Child(parentId uint64) []model.DistrictsChildRep {
	var rep []model.DistrictsChildRep
	_ = dao.Districts.Where(dao.Districts.Columns.ParentId, parentId).Fields(model.DistrictsChildRep{}).Scan(&rep)
	return rep
}

// CurrentCity 获取当前城市
func (s *districtsService) CurrentCity(ctx context.Context, req model.DistrictsCurrentCityReq) (rep model.DistrictsCurrentCityRep, err error) {
	geoRep, err := g.Client().Get(fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%s,%s&radius=1000&extensions=base", g.Cfg().GetString("amap.key"), req.Lng, req.Lat))
	if err != nil {
		return
	}
	//默认城市北京
	cityCode := model.DistrictsDefaultCityCode
	if bodyText, err := ioutil.ReadAll(geoRep.Body); err == nil {
		var Result struct {
			Status    string `json:"status"`
			Regeocode struct {
				AddressComponent struct {
					Province string `json:"province"`
					Adcode   string `json:"adcode"`
					District string `json:"district"`
					Country  string `json:"country"`
					Township string `json:"township"`
					Citycode string `json:"citycode"`
				} `json:"addressComponent"`
			} `json:"regeocode"`
			Info     string `json:"info"`
			Infocode string `json:"infocode"`
		}
		if err := json.Unmarshal(bodyText, &Result); err == nil && Result.Status == "1" {
			cityCode = Result.Regeocode.AddressComponent.Citycode
		}
	}
	err = dao.Districts.Ctx(ctx).Where(dao.Districts.Columns.CityCode, cityCode).Where(dao.Districts.Columns.Level, "city").Fields(rep).Limit(1).Scan(&rep)
	return
}

func (s *districtsService) Detail(ctx context.Context, id uint) (d model.Districts, err error) {
	err = dao.Districts.Ctx(ctx).WherePri(id).Scan(&d)
	return
}

// MapIdName 获取地区名IDMap
func (s *districtsService) MapIdName(ctx context.Context, ids []uint) map[uint]string {
	var list []model.Shop
	rep := map[uint]string{}
	_ = dao.Districts.Ctx(ctx).WhereIn(dao.Shop.Columns.Id, ids).Fields(dao.Shop.Columns.Id, dao.Shop.Columns.Name).Scan(&list)
	for _, d := range list {
		rep[d.Id] = d.Name
	}
	return rep
}

// GetByIds 获取地区名IDMap
func (s *districtsService) GetByIds(ctx context.Context, ids []uint) (list []model.Districts) {
	_ = dao.Districts.Ctx(ctx).WhereIn(dao.Shop.Columns.Id, ids).Scan(&list)
	return
}
