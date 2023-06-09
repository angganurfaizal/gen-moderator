package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/logger"
)

// @Summary BTC Generate receive wallet address
// @Description BTC Generate receive wallet address
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param request body request.CreateInscribeBtcReq true "Create a btc wallet address request"
// @Success 200 {object} response.InscribeBtcResp{}
// @Router /inscribe/receive-address [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) btcCreateInscribeBTC(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var reqBody request.CreateInscribeBtcReq
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			reqUsecase := &structure.InscribeBtcReceiveAddrRespReq{}
			err = copier.Copy(reqUsecase, reqBody)
			if err != nil {
				return nil, err
			}

			if len(reqUsecase.FileName) == 0 {
				return nil, errors.New("Filename is required")
			}

			if len(reqUsecase.WalletAddress) == 0 {
				return nil, errors.New("WalletAddress is required")
			}

			if ok, _ := btc.ValidateAddress("btc", reqUsecase.WalletAddress); !ok {
				return nil, errors.New("WalletAddress is invalid")
			}

			/*if reqUsecase.FeeRate != 15 && reqUsecase.FeeRate != 20 && reqUsecase.FeeRate != 25 {
				return nil, errors.New("fee rate is invalid")
			}*/

			if len(reqUsecase.File) == 0 {
				return nil, errors.New("file is invalid")
			}
			userUuid, ok := ctx.Value(utils.SIGNED_USER_ID).(string)
			if ok {
				reqUsecase.SetFields(
					reqUsecase.WithUserUuid(userUuid),
				)
			}
			userWalletAddress, ok := ctx.Value(utils.SIGNED_WALLET_ADDRESS).(string)
			if ok {
				reqUsecase.SetFields(
					reqUsecase.WithUserWallerAddress(userWalletAddress),
				)
			}
			if reqUsecase.InvalidAuthentic() {
				return nil, errors.New("Access token is required")
			}
			btcWallet, err := h.Usecase.CreateInscribeBTC(ctx, *reqUsecase)
			if err != nil {
				logger.AtLog.Logger.Error("CreateInscribeBTC failed",
					zap.Any("payload", reqBody),
					zap.Error(err),
				)
				return nil, err
			}
			// logger.AtLog.Logger.Info("CreateInscribeBTC successfully", zap.Any("response", zap.Any("btcWallet)", btcWallet)))
			return h.inscribeBtcCreatedRespResp(btcWallet)
		},
	).ServeHTTP(w, r)
}

// @Summary compress-image
// @Description compress-image
// @Tags compress-image
// @Accept json
// @Produce json
// @Param request body request.CompressImageReq true "compress images"
// @Router /inscribe/compress-image [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) compressImage(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var reqBody request.CompressImageReq
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}

			if len(reqBody.ImageUrl) == 0 {
				return nil, errors.New("imageUrl is required")
			}

			if len(reqBody.CompressPercents) == 0 || len(reqBody.CompressPercents) > 3 {
				return nil, errors.New("compressPercents invalid")
			}

			response, err := h.Usecase.CompressNftImageFromMoralis(ctx, reqBody.ImageUrl, reqBody.CompressPercents)
			if err != nil {
				logger.AtLog.Logger.Error("CompressNftImageFromMoralis failed",
					zap.Any("payload", reqBody),
					zap.Error(err),
				)
				return nil, err
			}
			logger.AtLog.Logger.Info("CompressNftImageFromMoralis successfully", zap.Any("response", zap.Any("response)", response)))
			return response, nil
		},
	).ServeHTTP(w, r)
}

func (h *httpDelivery) inscribeBtcCreatedRespResp(input *entity.InscribeBTC) (*response.InscribeBtcResp, error) {
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
	resp.EstFeeInfo = input.EstFeeInfo
	return resp, nil
}

// @Summary BTC List Inscribe
// @Description BTC List Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Success 200 {object} entity.Pagination{}
// @Router /inscribe/list [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) btcListInscribeBTC(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			page := entity.GetPagination(r)
			req := &entity.FilterInscribeBT{
				BaseFilters: entity.BaseFilters{
					Limit: page.PageSize,
					Page:  page.Page,
				},
				Expired: true,
			}
			userUuid, ok := ctx.Value(utils.SIGNED_USER_ID).(string)
			if ok {
				req.UserUuid = &userUuid
			} else {
				return nil, errors.New("access-token is required")
			}
			return h.Usecase.ListInscribeBTC(req)
		},
	).ServeHTTP(w, r)
}

// @Summary BTC NFT Detail Inscribe
// @Description BTC NFT Detail Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param ID path string true "inscribe ID"
// @Success 200 {object} entity.InscribeBTCResp{}
// @Router /inscribe/nft-detail/{ID} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) btcDetailInscribeBTC(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uuid := vars["ID"]

	result, err := h.Usecase.DetailDeveloperInscribeBTC(uuid)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.DetailDeveloperInscribeBTC", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

}

// @Summary BTC Retry Inscribe
// @Description BTC Retry Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param ID path string true "inscribe ID"
// @Success 200
// @Router /inscribe/retry/{ID} [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) btcRetryInscribeBTC(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["ID"]

	err := h.Usecase.RetryInscribeBTC(id)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.RetryInscribeBTC", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")

}

// @Summary BTC Info Inscribe
// @Description BTC Info Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param ID path string true "inscribe ID"
// @Success 200 {object} response.InscribeInfoResp{}
// @Router /inscribe/info/{ID} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) getInscribeInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["ID"]
	inscribeInfo, err := h.Usecase.GetInscribeInfo(id)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetInscribeInfo", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.inscribeInfoToResp(inscribeInfo)
	if err != nil {
		logger.AtLog.Logger.Error("h.inscribeInfoToResp", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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
