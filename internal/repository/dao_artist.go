package repository

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_artist"
	"rederinghub.io/utils/constants/dao_artist_voted"
	"rederinghub.io/utils/logger"
)

func (s Repository) ListDAOArtist(ctx context.Context, request *request.ListDaoArtistRequest) ([]*entity.DaoArtist, int64, error) {
	limit := int64(100)
	filters := make(bson.M)
	filterIdOperation := "$lt"
	sorts := bson.M{
		"$sort": bson.D{{Key: "_id", Value: -1}},
	}
	matchFilters := bson.M{"$match": filters}
	lookupUser := bson.M{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "created_by",
			"foreignField": "wallet_address",
			"as":           "user",
		},
	}
	unwindUser := bson.M{"$unwind": "$user"}
	addFieldUserName := bson.M{
		"$addFields": bson.M{
			"user_name":          "$user.display_name",
			"collection_created": "$user.stats.collection_created",
		},
	}
	if len(request.Sorts) > 0 {
		sort := bson.D{}
		for _, srt := range request.Sorts {
			sort = append(sort, bson.E{Key: srt.Field, Value: srt.Type})
			if srt.Field == "_id" && srt.Type == entity.SORT_ASC {
				filterIdOperation = "$gt"
			}
		}
		sorts = bson.M{
			"$sort": sort,
		}
	}
	if request.PageSize > 0 && request.PageSize <= limit {
		limit = request.PageSize
	}
	filters["$or"] = bson.A{
		bson.M{"expired_at": bson.M{"$gt": time.Now()}},
		bson.M{"status": dao_artist.Verified},
	}
	if request.Id != nil {
		id, err := primitive.ObjectIDFromHex(*request.Id)
		if err != nil {
			return nil, 0, err
		}
		filters["_id"] = id
	}
	if request.SeqId != nil {
		filters["seq_id"] = *request.SeqId
	}
	if request.Status != nil {
		filters["status"] = *request.Status
	}
	if request.Cursor != "" {
		if id, err := primitive.ObjectIDFromHex(request.Cursor); err == nil {
			filters["_id"] = bson.M{filterIdOperation: id}
		}
	}
	filterSearch := make(bson.M)
	matchSearch := bson.M{"$match": filterSearch}
	if request.Keyword != nil {
		search := bson.A{
			bson.M{"user_name": primitive.Regex{
				Pattern: *request.Keyword,
				Options: "i",
			}},
		}
		if seqId, err := strconv.Atoi(*request.Keyword); err == nil {
			search = append(search, bson.M{"seq_id": seqId})
		} else {
			if id, err := primitive.ObjectIDFromHex(*request.Keyword); err == nil {
				search = append(search, bson.M{"_id": id})
			}
		}
		filterSearch["$or"] = search
	}
	lookupDaoArtistVoted := bson.M{
		"$lookup": bson.M{
			"from":         "dao_artist_voted",
			"localField":   "_id",
			"foreignField": "dao_artist_id",
			"as":           "dao_artist_voted",
		},
	}
	addFieldsCount := bson.M{
		"$addFields": bson.M{
			"verify": bson.M{
				"$filter": bson.M{
					"input": "$dao_artist_voted",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", 1},
					},
				},
			},
			"report": bson.M{
				"$filter": bson.M{
					"input": "$dao_artist_voted",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.status", 0},
					},
				},
			},
		},
	}
	projectAgg := bson.M{
		"$project": bson.M{
			"_id":              1,
			"uuid":             1,
			"created_at":       1,
			"seq_id":           1,
			"created_by":       1,
			"user":             1,
			"expired_at":       1,
			"status":           1,
			"dao_artist_voted": 1,
			"user_name":        1,
			"total_verify":     bson.M{"$size": "$verify"},
			"total_report":     bson.M{"$size": "$report"},
		},
	}
	projects := []*entity.DaoArtist{}
	total, err := s.Aggregation(ctx,
		entity.DaoArtist{}.TableName(),
		0,
		limit,
		&projects,
		matchFilters,
		lookupUser,
		unwindUser,
		addFieldUserName,
		matchSearch,
		lookupDaoArtistVoted,
		addFieldsCount,
		sorts,
		projectAgg)
	if err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}

func (s Repository) CheckDAOArtistAvailableByUser(ctx context.Context, userWallet string) (*entity.DaoArtist, bool) {
	daoArtist := &entity.DaoArtist{}
	if err := s.FindOneBy(ctx, daoArtist.TableName(), bson.M{
		"created_by": userWallet,
		"$or": bson.A{
			bson.M{"expired_at": bson.M{"$gt": time.Now()}},
			bson.M{"status": dao_artist.Verified},
		},
	}, daoArtist); err != nil {
		return nil, false
	}
	return daoArtist, true
}
func (s Repository) SetExpireYourProposalArtist(ctx context.Context, userWallet string) error {
	filter := bson.M{
		"created_by": userWallet,
		"$or": bson.A{
			bson.M{"expired_at": bson.M{"$gt": time.Now()}},
			bson.M{"status": dao_artist.Verified},
		},
	}
	update := bson.M{"$set": bson.M{"expired_at": time.Now(), "status": dao_artist.Verifying}}
	count, err := s.UpdateMany(ctx, entity.DaoArtist{}.TableName(), filter, update)
	if err != nil {
		return err
	}
	logger.AtLog.Logger.Info("SetExpireYourProposalArtist success", zap.Int64("count", count))
	return nil
}
func (s Repository) CountDAOArtistVoteByStatus(ctx context.Context, daoArtistId primitive.ObjectID, status dao_artist_voted.Status) int {
	match := bson.M{"$match": bson.M{
		"dao_artist_id": daoArtistId,
		"status":        status,
	}}
	group := bson.M{
		"$group": bson.M{
			"_id":   "$dao_artist_id",
			"count": bson.M{"$sum": 1},
		},
	}
	cur, err := s.DB.Collection(entity.DaoArtistVoted{}.TableName()).Aggregate(ctx, bson.A{match, group})
	if err != nil {
		return 0
	}
	var results []*Count
	if err := cur.All(ctx, &results); err != nil {
		return 0
	}
	if len(results) > 0 {
		return results[0].Count
	}
	return 0
}
