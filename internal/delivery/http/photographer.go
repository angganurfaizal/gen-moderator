package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils/logger"
	_req "rederinghub.io/utils/request"
)

// Capture
// @Summary captures url content as image
// @Description captures url content as image
// @Tags Photographer
// @Accept  json
// @Produce  json
// @Param json body request.CaptureRequest true "capture request"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.CaptureResponse}
// @Router /photo/capture [POST]
func (h *httpDelivery) Capture(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			var req request.CaptureRequest
			err := _req.BindJson(r, &req)
			if err != nil {
				return nil, err
			}
			url, err := h.Usecase.CaptureContent(req.ID, req.Url)
			if err != nil {
				return nil, err
			}

			renderURL := fmt.Sprintf("%s/api/v1/device/%s/renderer-set-image", os.Getenv("RENDER_DOMAIN"), req.ID)
			code, result, err := _req.PostToRenderer(renderURL, map[string]string{
				"image_url": url,
			})
			logger.AtLog.Info("call to renderer-set-image ", zap.Error(err), zap.String("renderURL", renderURL), zap.String("device_id", req.ID), zap.Int("code", code), zap.String("response", string(result)))

			return response.CaptureResponse{
				ImageUrl: url,
			}, nil

		},
	).ServeHTTP(w, r)
}
