package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) InsertReferral(data *entity.Referral) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FilterReferrals(filter entity.FilterReferrals) bson.M {
	f := bson.M{}
	if filter.ReferreeID != nil {
		f["referree_id"] = filter.ReferreeID
	}
	if filter.ReferrerID != nil {
		f["referrer_id"] = filter.ReferrerID
	}
	return f
}

func (r Repository) GetReferrals(filter entity.FilterReferrals) (*entity.Pagination, error) {
	confs := []entity.Referral{}
	resp := &entity.Pagination{}
	f := r.FilterReferrals(filter)
	p, err := r.Paginate(utils.COLLECTION_REFERRALS, filter.Page, filter.Limit, f, bson.M{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

