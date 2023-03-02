package usecase

import (
	"fmt"
	"strconv"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils/algolia"
)

func (uc *Usecase) AlgoliaSearchProject(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "project" {
		return nil, 0, 0, nil
	}
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("projects", filter)
	if err != nil {
		return nil, 0, 0, err
	}

	projects := []*response.SearchProject{}
	resp.UnmarshalHits(&projects)
	dataResp := []*response.SearchResponse{}
	for _, i := range projects {
		mintPriceInt, err := strconv.ParseInt(i.MintPrice, 10, 64)
		if err == nil {
			i.MintPrice = fmt.Sprintf("%f", float64(mintPriceInt)/1e8)
			p := float64(mintPriceInt) / 1e8
			if p == float64(0) {
				i.MintPrice = "0"
			}
		}

		obj := &response.SearchResponse{
			ObjectType: "project",
			Project:    i,
		}
		dataResp = append(dataResp, obj)
	}

	return dataResp, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchInscription(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "inscription" {
		return nil, 0, 0, nil
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("inscriptions", filter)
	if err != nil {
		return nil, 0, 0, err
	}

	inscriptions := []*response.SearhcInscription{}
	for _, h := range resp.Hits {
		i := &response.SearhcInscription{
			ObjectId:      h["objectID"].(string),
			InscriptionId: h["inscription_id"].(string),
			Number:        int64(h["number"].(float64)),
			Sat:           fmt.Sprintf("%d", int64(h["sat"].(float64))),
			Chain:         h["chain"].(string),
			GenesisFee:    int64(h["genesis_fee"].(float64)),
			GenesisHeight: int64(h["genesis_height"].(float64)),
			Timestamp:     h["timestamp"].(string),
			ContentType:   h["content_type"].(string),
		}
		inscriptions = append(inscriptions, i)
	}
	resp.UnmarshalHits(&inscriptions)

	dataResp := []*response.SearchResponse{}
	for _, i := range inscriptions {
		obj := &response.SearchResponse{
			ObjectType:  "inscription",
			Inscription: i,
		}
		dataResp = append(dataResp, obj)
	}

	return dataResp, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchArtist(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "artist" {
		return nil, 0, 0, nil
	}
	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("users", filter)
	if err != nil {
		return nil, 0, 0, err
	}
	artists := []*response.SearchArtist{}
	resp.UnmarshalHits(&artists)

	dataResp := []*response.SearchResponse{}
	for _, i := range artists {
		obj := &response.SearchResponse{
			ObjectType: "artist",
			Artist:     i,
		}
		dataResp = append(dataResp, obj)
	}
	return dataResp, resp.NbHits, resp.NbPages, nil
}

func (uc *Usecase) AlgoliaSearchTokenUri(filter *algolia.AlgoliaFilter) ([]*response.SearchResponse, int, int, error) {
	if filter.ObjType != "" && filter.ObjType != "token" {
		return nil, 0, 0, nil
	}

	algoliaClient := algolia.NewAlgoliaClient(uc.Config.AlgoliaApplicationId, uc.Config.AlgoliaApiKey)

	resp, err := algoliaClient.Search("token-uris", filter)
	if err != nil {
		return nil, 0, 0, err
	}
	inscriptions := []*response.SearchTokenUri{}
	resp.UnmarshalHits(&inscriptions)

	dataResp := []*response.SearchResponse{}
	for _, i := range inscriptions {
		obj := &response.SearchResponse{
			ObjectType: "token",
			TokenUri:   i,
		}
		dataResp = append(dataResp, obj)
	}

	return dataResp, resp.NbHits, resp.NbPages, nil
}