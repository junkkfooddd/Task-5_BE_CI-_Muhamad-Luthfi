package delivery

import (
	"self-payrol/helper"
	"self-payrol/model"
	"self-payrol/request"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type companyDelivery struct {
	companyUsecase model.CompanyUsecase
}

type CompanyDelivery interface {
	Mount(group *echo.Group)
}

func NewCompanyDelivery(companyUsecase model.CompanyUsecase) CompanyDelivery {
	return &companyDelivery{companyUsecase: companyUsecase}
}

func (comp *companyDelivery) Mount(group *echo.Group) {

	group.GET("/company", comp.GetDetailCompanyHandler)
	group.POST("/company", comp.UpdateOrCreateCompanyHandler)
	group.POST("/topup", comp.TopupBalanceHandler)

}

func (comp *companyDelivery) GetDetailCompanyHandler(e echo.Context) error {
	ctx := e.Request().Context()

	info, i, err := comp.companyUsecase.GetCompanyInfo(ctx)
	if err != nil {
		return helper.ResponseErrorJson(e, i, err)
	}

	return helper.ResponseSuccessJson(e, "success", info)

}

func (comp *companyDelivery) UpdateOrCreateCompanyHandler(e echo.Context) error {
	ctx := e.Request().Context()

	var req request.CompanyRequest

	if err := e.Bind(&req); err != nil {
		return helper.ResponseValidationErrorJson(e, "Error binding struct", err.Error())
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return helper.ResponseValidationErrorJson(e, "Error validation", errVal)
	}

	company, i, err := comp.companyUsecase.CreateOrUpdateCompany(ctx, req)
	if err != nil {
		return helper.ResponseErrorJson(e, i, err)
	}

	return helper.ResponseSuccessJson(e, "success", company)
}

func (comp *companyDelivery) TopupBalanceHandler(e echo.Context) error {
	ctx := e.Request().Context()

	type TopupRequest struct {
		Amount float64 `json:"amount"`
	}

	// TODO: Perform request validation here
	var req TopupRequest
	if err := e.Bind(&req); err != nil {
		return helper.ResponseValidationErrorJson(e, "Error binding struct", err.Error())
	}

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required),
	); err != nil {
		errVal := err.(validation.Errors)
		return helper.ResponseValidationErrorJson(e, "Error validation", errVal)
	}

	topupReq := request.TopupCompanyBalance{Balance: req.Amount} // Create instance of TopupCompanyBalance

	company, i, err := comp.companyUsecase.TopupBalance(ctx, topupReq)

	if err != nil {
		return helper.ResponseErrorJson(e, i, err)
	}

	return helper.ResponseSuccessJson(e, "success", company)
}
