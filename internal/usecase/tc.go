package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (u Usecase) GetNftsByAddress(address string) (interface{}, error) {
	url := fmt.Sprintf("https://dapp.trustless.computer/dapp/api/nft-explorer/owner-address/%s/nfts", address)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result struct {
		Data []*struct {
			Collection        string `json:"collection"`
			CollectionAddress string `json:"collection_address"`
			TokenID           string `json:"token_id"`
			Name              string `json:"name"`
			ContentType       string `json:"content_type"`
			Image             string `json:"image"`
			Explorer          string `json:"explorer"`
			ArtistName        string `json:"artist_name"`
		} `json:"data"`
	}

	// parse:
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// var listContractID []string

	for _, nft := range result.Data {

		// listContractID = append(listContractID, nft.CollectionAddress)

		// if len(nft.Image) > 0 {
		// 	nft.Image += "/content"
		// }
		nft.Explorer = fmt.Sprintf("https://trustless.computer/inscription?contract=%s&id=%s", nft.CollectionAddress, nft.TokenID)

		p, _ := u.Repo.FindProjectByGenNFTAddr(nft.CollectionAddress)
		if p != nil {
			nft.ArtistName = p.Name
		}

	}
	return result.Data, err
}