package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)


func (r Repository) AggregateVolumn() ([]entity.AggregateWalleRespItem, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateWalleRespItem{}

	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": bson.A{
		bson.M{"status": entity.StatusMint_SentFundToMaster},
	}}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": 
			bson.M{ "projectID": "$projectID", "payType": "$payType" }, 
			"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
		}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}
	
	cursor, err := r.DB.Collection(entity.MintNftBtc{}.TableName()).Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, item := range results {
		res := &entity.AggregateWalletAddressItem{}
		err = helpers.Transform(item, res)
		if err != nil {
			return nil, err
		}
		tmp := entity.AggregateWalleRespItem{
			ProjectID: res.ID.ProjectID,
			Paytype: res.ID.Paytype,
			Amount: fmt.Sprintf("%d", int64(res.Amount)),
		}
		confs = append(confs, tmp)
	}
	
	return confs, nil
}

func (r Repository) AggregateAmount(walletAddress string, paymentType *string) ([]entity.AggregateAmount, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateAmount{}

	f :=  bson.A{
		bson.M{"creatorAddress": walletAddress},
	}

	if paymentType != nil && *paymentType != "" {
		f = append(f, bson.M{"amountType": paymentType})
	}

	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": f}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": 
			bson.M{"creatorAddress": "$creatorAddress"}, 
			"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
		}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}
	
	cursor, err := r.DB.Collection(entity.UserVolumn{}.TableName()).Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, item := range results {
		spew.Dump(item)
		tmp := &entity.AggregateAmount{}
		err = helpers.Transform(item, tmp)
		confs = append(confs, *tmp)
	}
	
	return confs, nil
}


func (r Repository) FindVolumn(projectID string, amountType string) (*entity.UserVolumn, error) {
	projectID = strings.ToLower(projectID)
	resp := &entity.UserVolumn{}
	usr, err := r.FilterOne(entity.UserVolumn{}.TableName(), bson.D{{"projectID", projectID}, {"amountType", amountType}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) CreateVolumn(input *entity.UserVolumn) error {
	err := r.InsertOne(input.TableName(), input)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateVolumn(ID string, data *entity.UserVolumn) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.UserVolumn{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r Repository) UpdateVolumnAmount(ID string, amount string) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	update := bson.M{"$set": bson.M{"amount": amount}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}