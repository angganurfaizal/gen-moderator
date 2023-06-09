package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"

	"rederinghub.io/internal/delivery/http/request"
)

// UserCredits godoc
// @Summary Generate a message
// @Description Generate a message for user's wallet
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body request.GenerateMessageRequest true "Generate message request"
// @Success 200 {object} response.JsonResponse{data=response.GeneratedMessage}
// @Router /auth/nonce [POST]
func (h *httpDelivery) generateMessage(w http.ResponseWriter, r *http.Request) {

	var reqBody request.GenerateMessageRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.SelfValidate()
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	logger.AtLog.Logger.Info("generateMessage", zap.Any("reqBody", zap.Any("reqBody)", reqBody)))
	message, err := h.Usecase.GenerateMessage(structure.GenerateMessage{
		Address:    *reqBody.Address,
		WalletType: reqBody.WalletType,
	})

	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	logger.AtLog.Logger.Info("resp.message", zap.Any("message", message))
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.GeneratedMessage{Message: *message}, "")
}

// UserCredits godoc
// @Summary Verified the generated message
// @Description Verified the generated message
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body request.VerifyMessageRequest true "Verify message request"
// @Success 200 {object} response.JsonResponse{data=response.VerifyResponse}
// @Router /auth/nonce/verify [POST]
func (h *httpDelivery) verifyMessage(w http.ResponseWriter, r *http.Request) {

	var reqBody request.VerifyMessageRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.SelfValidate()
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	logger.AtLog.Logger.Info("request.decoder", zap.Any("decoder", decoder))
	verifyMessage := structure.VerifyMessage{
		ETHSignature:     *reqBody.ETHSinature,
		Signature:        *reqBody.Sinature,
		Address:          *reqBody.Address,         //eth
		AddressBTC:       reqBody.AddressBTC,       //btc taproot addree -> use for transfer nft
		AddressBTCSegwit: reqBody.AddressBTCSegwit, //btc segwit address -> use for verify signature
		MessagePrefix:    reqBody.MessagePrefix,    //btc prefix message
		AddressPayment:   reqBody.AddressPayment,   //,
	}
	verified, err := h.Usecase.VerifyMessage(verifyMessage)

	logger.AtLog.Logger.Info("verified", zap.Any("verified", verified))
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.VerifyResponse{
		IsVerified:   verified.IsVerified,
		Token:        verified.Token,
		RefreshToken: verified.RefreshToken,
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
