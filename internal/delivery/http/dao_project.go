package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

// @Summary List DAO Project
// @Description List DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword"
// @Param status query int false "Status"
// @Param cursor query string false "Last Id"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.Pagination{}
// @Router /dao-project [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) listDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			req := &request.ListDaoProjectRequest{}
			if err := utils.QueryParser(r, req); err != nil {
				return nil, err
			}
			req.Pagination = entity.GetPagination(r)
			userWallet := muxVars[utils.SIGNED_WALLET_ADDRESS]
			resp, err := h.Usecase.ListDAOProject(ctx, userWallet, req)
			if err != nil {
				return &entity.Pagination{Result: make([]*response.DaoProject, 0)}, nil
			}
			return resp, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Create DAO Project
// @Description Create DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param request body request.CreateDaoProjectRequest true "Create Dao Project Request"
// @Success 200
// @Router /dao-project [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			var reqBody request.CreateDaoProjectRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			if err := h.Validator.Struct(reqBody); err != nil {
				return nil, err
			}
			reqBody.CreatedBy = muxVars[utils.SIGNED_WALLET_ADDRESS]
			return h.Usecase.CreateDAOProject(ctx, &reqBody)
		},
	).ServeHTTP(w, r)
}

// @Summary Get DAO Project
// @Description Get DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param id path string true "Dao Project Id"
// @Success 200 {object} response.DaoProject{}
// @Router /dao-project/{id} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) getDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return h.Usecase.GetDAOProject(ctx, muxVars["id"], muxVars[utils.SIGNED_WALLET_ADDRESS])
		},
	).ServeHTTP(w, r)
}

// @Summary Vote DAO Project
// @Description Vote DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param request body request.VoteDaoProjectRequest true "Vote Dao Project Request"
// @Param id path string true "Dao Project Id"
// @Success 200
// @Router /dao-project/{id} [PUT]
// @Security ApiKeyAuth
func (h *httpDelivery) voteDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			var reqBody request.VoteDaoProjectRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			return nil, h.Usecase.VoteDAOProject(ctx, muxVars["id"], muxVars[utils.SIGNED_WALLET_ADDRESS], &reqBody)
		},
	).ServeHTTP(w, r)
}

// @Summary List Projects Is Hidden
// @Description List Projects Is Hidden
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword"
// @Param cursor query string false "Last Id"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.Pagination{}
// @Router /dao-project/me/projects-hidden [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) listYourProjectsIsHidden(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			req := &request.ListProjectHiddenRequest{}
			if err := utils.QueryParser(r, req); err != nil {
				return nil, err
			}
			req.Pagination = entity.GetPagination(r)
			userWallet := muxVars[utils.SIGNED_WALLET_ADDRESS]
			if userWallet == "" {
				return nil, errors.New("token is empty")
			}
			return h.Usecase.ListYourProjectsIsHidden(ctx, userWallet, req)
		},
	).ServeHTTP(w, r)
}
