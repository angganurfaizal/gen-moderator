package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary BTC/ETH Generate receive wallet address
// @Description Generate receive wallet address
// @Tags BTC/ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc/eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /mint-nft-btc/receive-address [POST]
func (h *httpDelivery) createMintReceiveAddress(w http.ResponseWriter, r *http.Request) {

	// verify user:
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.CreateMintReceiveAddressReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("httpDelivery.MintNftBtc.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// if !strings.EqualFold(profile.WalletAddressBTCTaproot, reqBody.WalletAddress) {
	// 	err = errors.New("permission dined")
	// 	logger.AtLog.Logger.Error("h.Usecase.createMintReceiveAddress", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	reqUsecase := &structure.MintNftBtcData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase.UserID = profile.UUID
	reqUsecase.UserAddress = profile.WalletAddressBTCTaproot

	mintNftBtcWallet, err := h.Usecase.CreateMintReceiveAddress(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.createMintReceiveAddress", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := h.MintNftBtcToResp(mintNftBtcWallet)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary cancel the mint request
// @Description cancel the mint request
// @Tags BTC/ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc/eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /mint-nft-btc/receive-address [DELETE]
func (h *httpDelivery) cancelMintNftBt(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err = h.Usecase.CancelMintNftBtc(profile.WalletAddressBTCTaproot, uuid)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.CancelMintNftBt", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}

func (h *httpDelivery) getDetailMintNftBtc(w http.ResponseWriter, r *http.Request) {

	// verify user:
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	item, err := h.Usecase.GetDetalMintNftBtc(uuid)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.CancelMintNftBt", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	// if !strings.EqualFold(item.OriginUserAddress, profile.WalletAddressBTCTaproot) {
	// 	err := errors.New("permission dined")
	// 	logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }
	if item.UserID != profile.UUID {
		err := errors.New("permission dined")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, item, "")
}

func (h *httpDelivery) MintNftBtcToResp(input *entity.MintNftBtc) *response.MintNftBtcReceiveWalletResp {
	return &response.MintNftBtcReceiveWalletResp{
		Address:             input.ReceiveAddress,
		Price:               input.Amount,
		PayType:             input.PayType,
		NetworkFeeByPayType: input.NetworkFeeByPayType,
		MintPriceByPayType:  input.MintPriceByPayType,
		Quantity:            input.Quantity,
		MintFeeInfos:        input.EstFeeInfo,
	}

}

func (h *httpDelivery) MintNftBtcResp(input *entity.MintNftBtc) (*response.MintNftBtcResp, error) {
	resp := &response.MintNftBtcResp{}
	err := copier.Copy(resp, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *httpDelivery) getMintFeeRateInfos(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fileSize := vars["fileSize"]
	customRate := vars["customRate"]
	mintPrice := vars["mintPrice"]

	fileSizeInt, err := strconv.Atoi(fileSize)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetLevelFeeInfo", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	customRateInt, err := strconv.Atoi(customRate)
	if err != nil {
		customRateInt = 0
	}

	mintPriceInt, err := strconv.Atoi(mintPrice)
	if err != nil {
		mintPriceInt = 0
	}

	item, err := h.Usecase.GetLevelFeeInfo(int64(fileSizeInt), int64(customRateInt), int64(mintPriceInt))
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetLevelFeeInfo", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, item, "")
}
