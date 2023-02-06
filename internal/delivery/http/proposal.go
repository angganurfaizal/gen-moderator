package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary DAO list proposal
// @Description DAO list proposal
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param proposer query string false "filter by proposer"
// @Param proposalID query string false "filter by proposalID"
// @Param state query string false "filter by state"
// @Param sort query string false "newest, minted-newest, token-price-asc, token-price-desc"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals [GET]
func (h *httpDelivery) proposals(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.proposals", r)
	defer h.Tracer.FinishSpan(span, log )

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	f := structure.FilterProposal{}
	f.BaseFilters = *baseF

	proposer := r.URL.Query().Get("proposer")
	proposalID := r.URL.Query().Get("proposalID")
	state := r.URL.Query().Get("state")

	if proposer != "" {
		f.Proposer = &proposer
	}
	
	if proposalID != "" {
		f.ProposalID = &proposalID
	}
	
	if state != "" {
		stateINT, err := strconv.Atoi(state)
		if err ==   nil {
			f.State = &stateINT
		}
	}

	uProposals, err := h.Usecase.GetProposals(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	pResp :=  []response.ProposalResp{}
	iPro := uProposals.Result
	pro := iPro.([]entity.Proposal)
	for _, proItem := range pro {

		p, err := h.proposalToResp(&proItem)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}


	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success,  h.PaginationResp(uProposals, pResp), "")
}

// UserCredits godoc
// @Summary DAO proposal's detail
// @Description DAO proposal's detail
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param proposalID path string true "proposalID: the onchain ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals/{proposalID} [GET]
func (h *httpDelivery) getProposal(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getProposal", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	span.SetTag("proposalID", proposalID)

	proposal, err := h.Usecase.GetProposal(span, proposalID)
	if err != nil {
		log.Error("h.Usecase.GetProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp, err := h.proposalToResp(proposal)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	
	log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp , "")
}

// UserCredits godoc
// @Summary DAO create a draft proposal
// @Description DAO create a draft proposal
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param request body request.CreateProposalReq true "Create a draft proposal request"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals [POST]
func (h *httpDelivery) createDraftProposals(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.createDraftProposals", r)
	defer h.Tracer.FinishSpan(span, log )

	var reqBody request.CreateProposalReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProposaltReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	uProposals, err := h.Usecase.CreateDraftProposal(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	log.SetData("uProposals", uProposals)
	resp, err := h.proposalDetailToResp(uProposals)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}


// UserCredits godoc
// @Summary DAO map off and onchain proposal
// @Description DAO off and onchain proposal
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param ID path string true "ID: the offChain ID"
// @Param proposalID path string true "proposalID: the onchain ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals/{ID}/{proposalID} [PUT]
func (h *httpDelivery) mapOffAndOnChainProposal(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.createDraftProposals", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	iD := vars["ID"]
	span.SetTag("iD", iD)

	proposalID := vars["proposalID"]
	span.SetTag("proposalID", proposalID)

	uProposals, err := h.Usecase.MapOffToOnChainProposal(span, iD, proposalID)
	if err != nil {
		log.Error("h.Usecase.MapOffToOnChainProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	log.SetData("uProposals", uProposals)
	resp, err := h.proposalDetailToResp(uProposals)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}


func (h *httpDelivery) proposalToResp(input *entity.Proposal) (*response.ProposalResp, error) {
	resp := &response.ProposalResp{}
	err := response.CopyEntityToRes(resp, input)
	if err != nil {
		return nil, err
	}

	resp.Amount = input.ProposalDetail.Amount
	resp.Title = input.ProposalDetail.Title
	resp.Description = input.ProposalDetail.Description
	resp.TokenType = input.ProposalDetail.TokenType
	resp.ReceiverAddress = input.ProposalDetail.ReceiverAddress

	vote := response.ProposalVote{
		For: 10,
		Against: 5,
		Total:  15,
		PercentFor: (10 / 15) * 100,
		PercentAgainst: (5 / 15) * 100,
	}

	resp.Vote = vote
	return resp, nil
}

func (h *httpDelivery) proposalDetailToResp(input *entity.ProposalDetail) (*response.ProposalResp, error) {
	resp := &response.ProposalResp{}
	err := response.CopyEntityToRes(resp, input)
	if err != nil {
		return nil, err
	}

	vote := response.ProposalVote{
		For: 10,
		Against: 5,
		Total:  15,
		PercentFor: (10 / 15) * 100,
		PercentAgainst: (5 / 15) * 100,
	}

	resp.Vote = vote
	return resp, nil
}
