package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary BTC Generate receive wallet address
// @Description Generate receive wallet address
// @Tags BTC
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReqV2 true "Create a btc wallet address request"
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
	if len(reqUsecase.Name) == 0 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("name is invalid"))
		return
	}
	if reqUsecase.FeeRate != 5 && reqUsecase.FeeRate != 10 && reqUsecase.FeeRate != 15 {
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
	resp.TimeoutAt = fmt.Sprintf("%d", input.ExpiredAt.Unix())
	resp.SegwitAddress = input.SegwitAddress
	return resp, nil
}
