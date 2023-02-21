package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type StatusMint int

const (
	StatusMint_Pending      StatusMint = iota // 0: pending: waiting for fund
	StatusMint_ReceivedFund                   // 1: received fund from user

	StatusMint_Minting // 2: minting
	StatusMint_Minted  // 3: mint success

	StatusMint_SendingNFTToUser // 4: sending nft to user
	StatusMint_SentNFTToUser    // 5: send nft to user success:

	StatusMint_SendingFundToMaster // 6: sending nft to user
	StatusMint_SentFundToMaster    // 7: send btc to master:

	StatusMint_TxMintFailed // 8: tx mint failed

	StatusMint_NotEnoughBalance // 9: balance not enough

	StatusMint_Refunding // 10: refunding
	StatusMint_Refunded  // 11: refunding

	StatusMint_TxRefundFailed // 12: tx refund fund failed

)

type MintNftBtc struct {
	BaseEntity  `bson:",inline"`
	UserAddress string `bson:"user_address"` //user's wallet address from FE

	OriginUserAddress string `bson:"origin_user_address"` //user's wallet address from FE

	Amount string `bson:"amount"` // amount required

	PayType string `bson:"payType"` // eth/btc...

	ReceiveAddress string `bson:"receiveAddress"` // address is generated for receive user fund.
	PrivateKey     string `bson:"privateKey"`     // private key of the receive wallet.

	Balance string `bson:"balance"` // balance after check

	ExpiredAt time.Time `bson:"expired_at"`

	Status StatusMint `bson:"status"` // status for record

	TxMintNft    string `bson:"tx_mint_nft"`
	TxSendNft    string `bson:"tx_send_nft"`
	TxSendMaster string `bson:"tx_send_master"`

	FileURI string `bson:"fileURI"` // FileURI will be mount if OrdAddress get all amount

	InscriptionID string `bson:"inscriptionID"`

	ProjectID string `bson:"projectID"` //projectID

	// just log for users, not using for the job checking.
	IsConfirm    bool `bson:"isConfirm"`
	IsMinted     bool `bson:"isMinted"`
	IsSentUser   bool `bson:"isSentUser"`
	IsSentMaster bool `bson:"isSentMaster"`

	OutputMintNFT interface{} `bson:"output_mint_nft"`
	OutputSendNFT interface{} `bson:"output_send_nft"`
}

func (u MintNftBtc) TableName() string {
	return utils.MINT_NFT_BTC
}

func (u MintNftBtc) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type MintNftBtcLogs struct {
	BaseEntity  `bson:",inline"`
	RecordID    string      `bson:"record_id"`
	Table       string      `bson:"table"`
	Name        string      `bson:"name"`
	Status      interface{} `bson:"status"`
	RequestMsg  interface{} `bson:"request_msg"`
	ResponseMsg interface{} `bson:"response_msg"`
}

func (u MintNftBtcLogs) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

func (u MintNftBtcLogs) TableName() string {
	return "mint_nft_btc_logs"
}

type MintNftBtcResp struct {
	UUID string `bson:"uuid"`

	InscriptionID  string `bson:"inscriptionID"`  // tokenID in ETH
	ReceiveAddress string `bson:"receiveAddress"` // address is generated for receive user fund.

	Amount string `bson:"amount"`

	PayType string `bson:"payType"` // eth/btc...

	ExpiredAt time.Time `bson:"expired_at"`

	Status StatusMint `bson:"status"` // status for record

	TxSendNft string `bson:"tx_send_nft"`
	TxMintNft string `bson:"tx_mint_nft"`
}