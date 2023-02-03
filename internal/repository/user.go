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

func (r Repository) FindUserByWalletAddress(walletAddress string) (*entity.Users, error) {
	resp := &entity.Users{}

	cached, err := r.GetCache(utils.COLLECTION_USERS, walletAddress)
	if err ==  nil && cached != nil {
		err = helpers.Transform(cached, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	usr, err := r.FilterOne(utils.COLLECTION_USERS, bson.D{{utils.KEY_WALLET_ADDRESS, walletAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	r.CreateCache(utils.COLLECTION_USERS, walletAddress, resp)
	return resp, nil
}

func (r Repository) FindUserByID(userID string) (*entity.Users, error) {
	resp := &entity.Users{}

	usr, err := r.FindOne(utils.COLLECTION_USERS, userID)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindUserByEmail(email string) (*entity.Users, error) {
	resp := &entity.Users{}

	f := bson.D{
		{"$and", 
			bson.A{
				bson.D{{"email", email}},
				bson.D{{utils.KEY_DELETED_AT, nil}},
			},
		},
	}

	usr, err := r.FilterOne(utils.COLLECTION_USERS, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateUser(insertedUser *entity.Users) error {

	err := r.InsertOne(insertedUser.TableName(), insertedUser)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateUserByWalletAddress(walletAdress string, updateddUser *entity.Users) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_WALLET_ADDRESS, walletAdress}}
	result, err := r.UpdateOne(updateddUser.TableName(), filter, updateddUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateUserByID(userID string, updateddUser *entity.Users) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, userID}}
	result, err := r.UpdateOne(updateddUser.TableName(), filter, updateddUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateUserStats(walletAddress string, stats entity.UserStats) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: utils.KEY_WALLET_ADDRESS, Value: walletAddress}}
	update := bson.M{
		"$set": bson.M {
			"stats": stats,
		},
	}
	result, err := r.DB.Collection(utils.COLLECTION_USERS).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) FindUserByAutoUserID(autoUserID int32) (*entity.Users, error) {
	resp := &entity.Users{}

	usr, err := r.FilterOne(utils.COLLECTION_USERS, bson.D{{utils.KEY_AUTO_USERID, autoUserID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) ListUsers(filter entity.FilterUsers) (*entity.Pagination, error)  {
	users := []entity.Users{}
	resp := &entity.Pagination{}

	filter1 := bson.M{}
	filter1[utils.KEY_DELETED_AT] = nil

	if filter.Email != nil {
		if *filter.Email != "" {
			filter1["email"] = *filter.Email
		}
	}
	
	if filter.WalletAddress != nil {
		if *filter.WalletAddress != "" {
			filter1["wallet_address"] = *filter.WalletAddress
		}
	}
	
	if filter.UserType != nil {
		filter1["user_type"] = *filter.UserType
	}
	
	p, err := r.Paginate(utils.COLLECTION_USERS, filter.Page, filter.Limit, filter1,bson.D{}, []Sort{}, &users)
	if err != nil {
		return nil, err
	}
	
	resp.Result = users
	//resp.Limit = p.Pagination.PerPage
	resp.Page = p.Pagination.Page
	// resp.Next = p.Pagination.Next
	// resp.Prev = p.Pagination.Prev
	// resp.TotalPage = p.Pagination.TotalPage
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetAllUsers(filter entity.FilterUsers) ([]entity.Users, error)  {
	users := []entity.Users{}
	f := bson.M{}
	if filter.IsUpdatedAvatar != nil {
		f["is_updated_avatar"] = *filter.IsUpdatedAvatar
	}else{
		f["is_updated_avatar"] = nil
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_USERS).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

// find all user profile, exclude avatar field for ease of optimization
func (r Repository) GetAllUserProfiles() ([]entity.Users, error)  {
	users := []entity.Users{}
	f := bson.M{}

	opts := options.Find().SetProjection(bson.D{{Key: "avatar", Value: 0}})

	cursor, err := r.DB.Collection(utils.COLLECTION_USERS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
