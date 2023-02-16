package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/btc"
)

// UserCredits godoc
// @Summary BTC Generate receive wallet address
// @Description Generate receive wallet address
// @Tags BTC
// @Accept  json
// @Produce  json
// @Param request body request.CreateInscribeBtcReq true "Create a btc wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /inscribe/receive-address [POST]
func (h *httpDelivery) btcCreateInscribeBTC(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.btcCreateInscribeBTC", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateInscribeBtcReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcCreateInscribeBTC.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.InscribeBtcReceiveAddrRespReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if reqUsecase == nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("invalid param"))
		return
	}

	if len(reqUsecase.FileName) == 0 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Filename is required"))
		return
	}

	if len(reqUsecase.WalletAddress) == 0 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("WalletAddress is required"))
		return
	}

	if ok, _ := btc.ValidateAddress("btc", reqUsecase.WalletAddress); !ok {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("WalletAddress is invalid"))
		return
	}

	if reqUsecase.FeeRate != 15 && reqUsecase.FeeRate != 20 && reqUsecase.FeeRate != 25 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("fee rate is invalid"))
		return
	}

	if len(reqUsecase.File) == 0 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("file is invalid"))
		return
	}

	btcWallet, err := h.Usecase.CreateInscribeBTC(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.btcCreateInscribeBTC", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("btcCreateInscribeBTC", btcWallet)
	resp, err := h.InscribeBtcCreatedRespResp(btcWallet)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) InscribeBtcCreatedRespResp(input *entity.InscribeBTC) (*response.InscribeBtcResp, error) {
	resp := &response.InscribeBtcResp{}
	resp.UserAddress = input.UserAddress
	resp.Amount = input.Amount
	resp.MintFee = input.MintFee
	resp.SentTokenFee = input.SentTokenFee
	resp.OrdAddress = input.OrdAddress
	resp.FileURI = input.FileURI
	resp.IsConfirm = input.IsConfirm
	resp.InscriptionID = input.InscriptionID
	resp.Balance = input.Balance
	resp.TimeoutAt = fmt.Sprintf("%d", time.Now().Add(time.Hour*1).Unix()) // return FE in 1h. //TODO: need update
	resp.SegwitAddress = input.SegwitAddress
	return resp, nil
}

func (h *httpDelivery) btcListInscribeBTC(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("btcListInscribeBTC", r)
	defer h.Tracer.FinishSpan(span, log)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	result, err := h.Usecase.ListInscribeBTC(span, int64(limit), int64(page))
	if err != nil {
		log.Error("h.Usecase.ListInscribeBTC", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

// detail:
func (h *httpDelivery) btcDetailInscribeBTC(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("btcDetailInscribeBTC", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	uuid := vars["ID"]

	result, err := h.Usecase.DetailInscribeBTC(uuid)
	if err != nil {
		log.Error("h.Usecase.DetailInscribeBTC", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

}

func (h *httpDelivery) btcRetryInscribeBTC(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("btcRetryInscribeBTC", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	id := vars["ID"]

	err := h.Usecase.RetryInscribeBTC(id)
	if err != nil {
		log.Error("h.Usecase.RetryInscribeBTC", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")

}

// UserCredits godoc
// @Summary get inscribe info
// @Description get inscribe info
// @Tags Inscribe
// @Accept  json
// @Produce  json
// @Param ID path string true "inscribe ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /inscribe/info/{ID} [GET]
func (h *httpDelivery) getInscribeInfo(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getInscribeInfo", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	id := vars["ID"]
	inscribeInfo, err := h.Usecase.GetInscribeInfo(span, id)
	if err != nil {
		log.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.inscribeInfoToResp(inscribeInfo)
	if err != nil {
		log.Error("h.inscribeInfoToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) inscribeInfoToResp(input *entity.InscribeInfo) (*response.InscribeInfoResp, error) {
	resp := &response.InscribeInfoResp{}
	resp.ID = input.ID
	resp.Index = input.Index
	resp.Address = input.Address
	resp.OutputValue = input.OutputValue
	resp.Sat = input.Sat
	resp.Preview = input.Preview
	resp.Content = input.Content
	resp.ContentLength = input.ContentLength
	resp.ContentType = input.ContentType
	resp.Timestamp = input.Timestamp
	resp.GenesisHeight = input.GenesisHeight
	resp.GenesisTransaction = input.GenesisTransaction
	resp.Location = input.Location
	resp.Output = input.Output
	resp.Offset = input.Offset
	return resp, nil
}
