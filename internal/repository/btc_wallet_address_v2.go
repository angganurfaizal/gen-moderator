package repository

import (
	"context"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r Repository) FindBtcWalletAddressV2(key string) (*entity.BTCWalletAddressV2, error) {
	resp := &entity.BTCWalletAddressV2{}
	usr, err := r.FilterOne(entity.BTCWalletAddress{}.TableName(), bson.D{{utils.KEY_UUID, key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindBtcWalletAddressByOrdV2(ordAddress string) (*entity.BTCWalletAddressV2, error) {
	resp := &entity.BTCWalletAddressV2{}
	usr, err := r.FilterOne(entity.BTCWalletAddressV2{}.TableName(), bson.D{{"ordAddress", ordAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) DeleteBtcWalletAddressV2(id string) (*mongo.DeleteResult, error) {
	filter := bson.D{{utils.KEY_UUID, id}}
	result, err := r.DeleteOne(entity.BTCWalletAddressV2{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) InsertBtcWalletAddressV2(data *entity.BTCWalletAddressV2) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListBtcWalletAddressV2(filter entity.FilterBTCWalletAddress) (*entity.Pagination, error)  {
	confs := []entity.BTCWalletAddressV2{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.BTCWalletAddressV2{}.TableName(), filter.Page, filter.Limit, f, bson.D{},[]Sort{} , &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateBtcWalletAddressByOrdAddrV2(ordAddress string, conf *entity.BTCWalletAddressV2) (*mongo.UpdateResult, error) {
	filter := bson.D{{"ordAddress", ordAddress}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) ListProcessingWalletAddressV2() ([]entity.BTCWalletAddressV2, error)  {
	confs := []entity.BTCWalletAddressV2{}
	f := bson.M{}
	f["$or"] = []interface{}{
		bson.M{"isMinted": bson.M{"$not": bson.M{"$eq": true}}} ,
		bson.M{"isConfirm": bson.M{"$not": bson.M{"$eq": true}}} ,
	}
	
	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_WALLET_ADDRESS_V2).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) ListBTCAddressV2() ([]entity.BTCWalletAddressV2, error)  {
	confs := []entity.BTCWalletAddressV2{}
	
	f := bson.M{}
	f["mintResponse"] = bson.M{"$not": bson.M{"$eq": nil}}
	f["mintResponse.issent"] = false
	f["mintResponse.inscription"] = bson.M{"$not": bson.M{"$eq": ""}}
	
	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_WALLET_ADDRESS_V2).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

// func (r Repository) FilterActiveBitcoinAddressV2(filter entity.FilterBTCWalletAddress) (*entity.Pagination, error)  {
// 	confs := []entity.BTCWalletAddressV2{}
// 	resp := &entity.Pagination{}
// 	f := bson.M{}
// 	f["inscriptionID"] = bson.M{"$ne": ""}

// 	p, err := r.Paginate(entity.BTCWalletAddressV2{}.TableName(), filter.Page, filter.Limit, f, bson.D{},[]Sort{} , &confs)
// 	if err != nil {
// 		return nil, err
// 	}
	
// 	resp.Result = confs
// 	resp.Page = p.Pagination.Page
// 	resp.Total = p.Pagination.Total
// 	resp.PageSize = filter.Limit
// 	return resp, nil
// }
