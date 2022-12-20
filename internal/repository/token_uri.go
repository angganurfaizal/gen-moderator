package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

func (r Repository) FindTokenBy(contractAddress string, tokenID string) (*entity.TokenUri, error) {
	resp := &entity.TokenUri{}

	usr, err := r.FilterOne(entity.TokenUri{}.TableName(), bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateTokenURI(data *entity.TokenUri) error {

	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}

	return nil
}
