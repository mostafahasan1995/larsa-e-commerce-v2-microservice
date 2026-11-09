package order_handlers

import (
	"encoding/json"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func (h *OrderHandler) NewOrderReportHandler(i *do.Injector, r chi.Router) {
	r.Post("/report", helpers.Make(h.GetOrdersByFilterV2))
	r.Post("/v2report", helpers.Make(h.GetOrdersByFilter))
	r.Post("/report/salesTable", helpers.Make(h.GetFilteredOrdersStatsV2))
	// r.Post("/v2report/salesTable", helpers.Make(h.GetFilteredOrdersStats))
	r.Post("/v2report/branchSales", helpers.Make(h.GetBranchSalesV2))
	r.Post("/report/branchSales", helpers.Make(h.GetBranchSalesV3))

	r.Post("/v2report/branchSales", helpers.Make(h.GetBranchSales))
	r.Post("/report/topProduct", helpers.Make(h.GetProductSalesV3))

	r.Post("/v1report/topProduct", helpers.Make(h.GetProductSalesV2))
	r.Post("/v2report/topProduct", helpers.Make(h.GetProductSales))

	r.Post("/v2report/productReport", helpers.Make(h.GetReportDependProductV2))
	r.Post("/report/productReport", helpers.Make(h.GetReportDependProductV3))

	r.Post("/report/productSalesReport", helpers.Make(h.GetProductSalesReportV2))
	r.Post("/v2report/productSalesReport", helpers.Make(h.GetProductSalesReport))

	// not use
	r.Post("/reportv2", helpers.Make(h.GetOrdersByItemFilter))
}

func (h *OrderHandler) GetBranchSales(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetBranchSales(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}
func (h *OrderHandler) GetBranchSalesV2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetBranchSalesV2(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetBranchSalesV3(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetBranchSalesV3(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetProductSales(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetProductSales(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}
func (h *OrderHandler) GetProductSalesV2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetProductSalesV2(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}
func (h *OrderHandler) GetProductSalesV3(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetProductSalesV3(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetReportDependProductV2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetReportDependProductV2(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}
func (h *OrderHandler) GetReportDependProductV3(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetReportDependProductV3(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrdersByFilter(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetOrdersByFilter(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrdersByFilterV2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	cfg, _ := util.GetReqAppCfg(ctx)
	fmt.Println(cfg.Db)
	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetOrdersByFilterV2(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetProductSalesReport(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetProductSalesReport(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}
func (h *OrderHandler) GetProductSalesReportV2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetProductSalesReportV2(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrdersByItemFilter(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderItemFilterDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	orders, apiErr := h.orderReportsService.GetOrdersByItemFilter(ctx, &filterDto)

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetFilteredOrdersStatsV2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var filterDto order_dto.OrderReportDTO
	if err := json.NewDecoder(r.Body).Decode(&filterDto); err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	//f
	orders, apiErr := h.orderReportsService.GetFilteredOrdersStatsV2(ctx, &filterDto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, orders)
}
