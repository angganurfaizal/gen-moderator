package request

type CreateDexBTCListing struct {
	RawPSBT       string `json:"raw_psbt"`
	InscriptionID string `json:"inscription_id"`
}

type CancelDexBTCListing struct {
	Txhash        string `json:"txhash"`
	InscriptionID string `json:"inscription_id"`
	OrderID       string `json:"order_id"`
}