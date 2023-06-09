package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type DexVolumeInscription struct {
	BaseEntity `bson:",inline"`
	Timestamp  *time.Time                   `bson:"timestamp"`
	Metadata   DexVolumeInscriptionMetadata `bson:"metadata"`
	Amount     uint64                       `bson:"amount"`
}

type DexVolumeInscriptionSumary struct {
	DexVolumeInscription *DexVolumeInscription `json:"dex_volume_inscription" bson:"dex_volume_inscription"`
	DexBTCListings       []*DexBTCListing      `json:"dex_btc_listings" bson:"dex_btc_listings"`
	TotalVolume          uint64                `json:"total_volume" bson:"total_volume"`
	Volume1h             uint64                `json:"volume_1h" bson:"volume_1h"`
	Volume1d             uint64                `json:"volume_1d" bson:"volume_1d"`
	Volume7d             uint64                `json:"volume_7d" bson:"volume_7d"`
	Volume30d            uint64                `json:"volume_30d" bson:"volume_30d"`
	InscriptionId        string                `json:"inscription_id" bson:"inscription_id"`
}

type AggerateChartForProject struct {
	ProjectID *string    `json:"projectID"`
	FromDate  *time.Time `json:"fromDate"`
	ToDate    *time.Time `json:"toDate"`
}

type AggerateChartForToken struct {
	TokenID  *string    `json:"tokenID"`
	FromDate *time.Time `json:"fromDate"`
	ToDate   *time.Time `json:"toDate"`
}

type AggragetedProject struct {
	ID     AggragetedProjectID `bson:"_id"`
	Amount int64               `bson:"amount"`
}
type AggragetedProjectID struct {
	ProjectID    string `bson:"projectID"`
	ProjectName  string `bson:"projectName"`
	Timestamp    string `bson:"timestamp"`
	Erc20Address string `bson:"erc20Address"`
}

type DexVolumeInscriptionMetadata struct {
	InscriptionId string `bson:"inscription_id"`
	MatchedTx     string `bson:"matched_tx"`
}

func (u DexVolumeInscription) TableName() string {
	return "dex_volume_inscription"
}

func (u DexVolumeInscription) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type AggragetedToken struct {
	ID     AggragetedTokenID `bson:"_id"`
	Amount int64             `bson:"amount"`
}
type AggragetedTokenID struct {
	TokenID     string `bson:"tokenID"`
	ProjectName string `bson:"projectName"`
	Timestamp   string `bson:"timestamp"`
}
