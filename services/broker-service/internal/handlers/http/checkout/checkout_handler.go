package checkout

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"broker-service/internal/util"
	checkoutpb "broker-service/proto/checkout"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	CheckoutClient checkoutpb.CheckoutServiceClient
}

// CreateOrder creates a new order
func (c *Config) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var requestPayload struct {
		Items       []OrderItemRequest `json:"items"`
		BillingInfo struct {
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			CompanyName string `json:"companyName,omitempty"`
			Address     string `json:"address"`
			Country     string `json:"country"`
			Region      string `json:"region"`
			City        string `json:"city"`
			ZipCode     string `json:"zipCode"`
			Email       string `json:"email"`
			Phone       string `json:"phone"`
		} `json:"billingInfo"`
		ShippingInfo struct {
			ShipToDifferentAddress bool   `json:"shipToDifferentAddress"`
			FirstName              string `json:"firstName,omitempty"`
			LastName               string `json:"lastName,omitempty"`
			CompanyName            string `json:"companyName,omitempty"`
			Address                string `json:"address,omitempty"`
			Country                string `json:"country,omitempty"`
			Region                 string `json:"region,omitempty"`
			City                   string `json:"city,omitempty"`
			ZipCode                string `json:"zipCode,omitempty"`
			ShippingMethod         string `json:"shippingMethod"`
		} `json:"shippingInfo"`
		PaymentMethod string `json:"paymentMethod"`
		CouponCode    string `json:"couponCode,omitempty"`
		Notes         string `json:"notes,omitempty"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Convert request items to proto message format
	items := make([]*checkoutpb.OrderItem, 0)
	for _, item := range requestPayload.Items {
		orderItem := &checkoutpb.OrderItem{
			ProductId: strconv.FormatInt(item.ProductID, 10),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			ImageUrl:  item.ImageURL,
		}
		items = append(items, orderItem)
	}

	// Create billing info proto message
	billingInfo := &checkoutpb.BillingInfo{
		FirstName:   requestPayload.BillingInfo.FirstName,
		LastName:    requestPayload.BillingInfo.LastName,
		CompanyName: requestPayload.BillingInfo.CompanyName,
		Address:     requestPayload.BillingInfo.Address,
		Country:     requestPayload.BillingInfo.Country,
		Region:      requestPayload.BillingInfo.Region,
		City:        requestPayload.BillingInfo.City,
		ZipCode:     requestPayload.BillingInfo.ZipCode,
		Email:       requestPayload.BillingInfo.Email,
		Phone:       requestPayload.BillingInfo.Phone,
	}

	// Create shipping info proto message
	shippingInfo := &checkoutpb.ShippingInfo{
		ShipToDifferentAddress: requestPayload.ShippingInfo.ShipToDifferentAddress,
		FirstName:              requestPayload.ShippingInfo.FirstName,
		LastName:               requestPayload.ShippingInfo.LastName,
		CompanyName:            requestPayload.ShippingInfo.CompanyName,
		Address:                requestPayload.ShippingInfo.Address,
		Country:                requestPayload.ShippingInfo.Country,
		Region:                 requestPayload.ShippingInfo.Region,
		City:                   requestPayload.ShippingInfo.City,
		ZipCode:                requestPayload.ShippingInfo.ZipCode,
		ShippingMethod:         requestPayload.ShippingInfo.ShippingMethod,
	}

	// Call checkout service via gRPC
	res, err := c.CheckoutClient.CreateOrder(r.Context(), &checkoutpb.CreateOrderRequest{
		UserId:        userID,
		Items:         items,
		BillingInfo:   billingInfo,
		ShippingInfo:  shippingInfo,
		PaymentMethod: requestPayload.PaymentMethod,
		CouponCode:    requestPayload.CouponCode,
		Notes:         requestPayload.Notes,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Order created successfully",
		Data:    formatOrderResponse(res),
	}

	util.WriteJSON(w, http.StatusCreated, responseData)
}

// GetOrder retrieves an order by ID
func (c *Config) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Get orderID from URL params
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		util.ErrorJSON(w, errors.New("order ID is required"), http.StatusBadRequest)
		return
	}

	// Call checkout service via gRPC
	res, err := c.CheckoutClient.GetOrder(r.Context(), &checkoutpb.GetOrderRequest{
		Id:     orderID,
		UserId: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Order retrieved successfully",
		Data:    formatOrderResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// ListOrders retrieves all orders for a user
func (c *Config) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Check the path to determine if this is an admin request
	// The admin-specific endpoint is routed differently and has middleware
	isAdminRequest := false
	if r.URL.Path == "/api/v1/checkout/admin/orders" {
		// Double check role for safety
		role, _ := r.Context().Value("role").(string)
		if role != "admin" {
			util.ErrorJSON(w, errors.New("admin privileges required"), http.StatusForbidden)
			return
		}
		isAdminRequest = true
	}

	// Log debug info
	log.Printf("ListOrders: userID=%s, isAdminRequest=%v, path=%s", userID, isAdminRequest, r.URL.Path)

	// Get pagination parameters
	page, pageSize := 1, 10 // Default values
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := r.URL.Query().Get("limit"); pageSizeStr != "" {
		if l, err := strconv.Atoi(pageSizeStr); err == nil && l > 0 {
			pageSize = l
		}
	}

	// Call checkout service via gRPC
	var res *checkoutpb.ListOrdersResponse
	var err error

	if isAdminRequest {
		// Admin user - retrieving all orders
		log.Printf("Admin request - retrieving all orders (ignoring user filter)")
		res, err = c.CheckoutClient.ListOrders(r.Context(), &checkoutpb.ListOrdersRequest{
			UserId:   "", // Empty user ID signals backend to ignore user filter for admins
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
	} else {
		// Regular user - retrieving only their orders
		log.Printf("Regular user request - retrieving orders for user ID: %s", userID)
		res, err = c.CheckoutClient.ListOrders(r.Context(), &checkoutpb.ListOrdersRequest{
			UserId:   userID,
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
	}

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Format orders
	orders := make([]map[string]interface{}, 0)
	for _, order := range res.Orders {
		orders = append(orders, formatOrderResponse(order))
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Orders retrieved successfully",
		Data: map[string]interface{}{
			"orders": orders,
			"total":  res.Total,
			"page":   page,
			"limit":  pageSize,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// ProcessPayment processes the payment for an order
func (c *Config) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Parse request
	var paymentRequest struct {
		OrderID       string  `json:"orderId"`
		PaymentMethod string  `json:"paymentMethod"`
		PaymentAmount float64 `json:"paymentAmount"`
		Currency      string  `json:"currency"`
		ReturnURL     string  `json:"returnUrl,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call checkout service via gRPC
	res, err := c.CheckoutClient.ProcessPayment(r.Context(), &checkoutpb.ProcessPaymentRequest{
		OrderId:       paymentRequest.OrderID,
		PaymentMethod: paymentRequest.PaymentMethod,
		Amount:        paymentRequest.PaymentAmount,
		Currency:      paymentRequest.Currency,
		ReturnUrl:     paymentRequest.ReturnURL,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Payment processed successfully",
		Data: map[string]interface{}{
			"transactionId": res.TransactionId,
			"status":        res.Status,
			"redirectUrl":   res.RedirectUrl,
			"success":       res.Success,
			"message":       res.Message,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// ValidateCheckout validates checkout information
func (c *Config) ValidateCheckout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Items       []OrderItemRequest `json:"items"`
		BillingInfo struct {
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			CompanyName string `json:"companyName,omitempty"`
			Address     string `json:"address"`
			Country     string `json:"country"`
			Region      string `json:"region"`
			City        string `json:"city"`
			ZipCode     string `json:"zipCode"`
			Email       string `json:"email"`
			Phone       string `json:"phone"`
		} `json:"billingInfo"`
		ShippingInfo struct {
			ShipToDifferentAddress bool   `json:"shipToDifferentAddress"`
			FirstName              string `json:"firstName,omitempty"`
			LastName               string `json:"lastName,omitempty"`
			CompanyName            string `json:"companyName,omitempty"`
			Address                string `json:"address,omitempty"`
			Country                string `json:"country,omitempty"`
			Region                 string `json:"region,omitempty"`
			City                   string `json:"city,omitempty"`
			ZipCode                string `json:"zipCode,omitempty"`
			ShippingMethod         string `json:"shippingMethod"`
		} `json:"shippingInfo"`
		PaymentMethod string `json:"paymentMethod"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Convert request items to proto message format
	items := make([]*checkoutpb.OrderItem, 0)
	for _, item := range requestPayload.Items {
		orderItem := &checkoutpb.OrderItem{
			ProductId: strconv.FormatInt(item.ProductID, 10),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			ImageUrl:  item.ImageURL,
		}
		items = append(items, orderItem)
	}

	// Create billing info proto message
	billingInfo := &checkoutpb.BillingInfo{
		FirstName:   requestPayload.BillingInfo.FirstName,
		LastName:    requestPayload.BillingInfo.LastName,
		CompanyName: requestPayload.BillingInfo.CompanyName,
		Address:     requestPayload.BillingInfo.Address,
		Country:     requestPayload.BillingInfo.Country,
		Region:      requestPayload.BillingInfo.Region,
		City:        requestPayload.BillingInfo.City,
		ZipCode:     requestPayload.BillingInfo.ZipCode,
		Email:       requestPayload.BillingInfo.Email,
		Phone:       requestPayload.BillingInfo.Phone,
	}

	// Create shipping info proto message
	shippingInfo := &checkoutpb.ShippingInfo{
		ShipToDifferentAddress: requestPayload.ShippingInfo.ShipToDifferentAddress,
		FirstName:              requestPayload.ShippingInfo.FirstName,
		LastName:               requestPayload.ShippingInfo.LastName,
		CompanyName:            requestPayload.ShippingInfo.CompanyName,
		Address:                requestPayload.ShippingInfo.Address,
		Country:                requestPayload.ShippingInfo.Country,
		Region:                 requestPayload.ShippingInfo.Region,
		City:                   requestPayload.ShippingInfo.City,
		ZipCode:                requestPayload.ShippingInfo.ZipCode,
		ShippingMethod:         requestPayload.ShippingInfo.ShippingMethod,
	}

	// Call checkout service via gRPC
	res, err := c.CheckoutClient.ValidateCheckout(r.Context(), &checkoutpb.ValidateCheckoutRequest{
		Items:         items,
		BillingInfo:   billingInfo,
		ShippingInfo:  shippingInfo,
		PaymentMethod: requestPayload.PaymentMethod,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Format validation errors
	errors := make([]map[string]string, 0)
	for _, validationError := range res.Errors {
		errors = append(errors, map[string]string{
			"field":   validationError.Field,
			"message": validationError.Message,
		})
	}

	responseData := util.JsonResponse{
		Error:   !res.Valid,
		Message: "Checkout validation completed",
		Data: map[string]interface{}{
			"valid":  res.Valid,
			"errors": errors,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// OrderItemRequest is used for order item request parsing
type OrderItemRequest struct {
	ProductID int64   `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int32   `json:"quantity"`
	ImageURL  string  `json:"imageUrl,omitempty"`
}

// Helper function to format order response for the frontend
func formatOrderResponse(order *checkoutpb.Order) map[string]interface{} {
	// Format order items
	items := make([]map[string]interface{}, 0)
	for _, item := range order.Items {
		// Convert ProductId string to number if possible
		productId, err := strconv.ParseInt(item.ProductId, 10, 64)
		if err != nil {
			productId = 0 // Default value if conversion fails
		}

		orderItem := map[string]interface{}{
			"id":        item.Id,
			"productId": productId,
			"name":      item.Name,
			"price":     item.Price,
			"quantity":  item.Quantity,
			"subtotal":  item.Subtotal,
			"imageUrl":  item.ImageUrl,
		}
		items = append(items, orderItem)
	}

	// Format billing info
	billingInfo := map[string]interface{}{
		"firstName":   order.BillingInfo.FirstName,
		"lastName":    order.BillingInfo.LastName,
		"companyName": order.BillingInfo.CompanyName,
		"address":     order.BillingInfo.Address,
		"country":     order.BillingInfo.Country,
		"region":      order.BillingInfo.Region,
		"city":        order.BillingInfo.City,
		"zipCode":     order.BillingInfo.ZipCode,
		"email":       order.BillingInfo.Email,
		"phone":       order.BillingInfo.Phone,
	}

	// Format shipping info
	shippingInfo := map[string]interface{}{
		"shipToDifferentAddress": order.ShippingInfo.ShipToDifferentAddress,
		"firstName":              order.ShippingInfo.FirstName,
		"lastName":               order.ShippingInfo.LastName,
		"companyName":            order.ShippingInfo.CompanyName,
		"address":                order.ShippingInfo.Address,
		"country":                order.ShippingInfo.Country,
		"region":                 order.ShippingInfo.Region,
		"city":                   order.ShippingInfo.City,
		"zipCode":                order.ShippingInfo.ZipCode,
		"shippingMethod":         order.ShippingInfo.ShippingMethod,
		"shippingCost":           order.ShippingInfo.ShippingCost,
	}

	// Format payment info
	paymentInfo := map[string]interface{}{
		"paymentMethod": order.PaymentInfo.PaymentMethod,
		"transactionId": order.PaymentInfo.TransactionId,
		"status":        order.PaymentInfo.Status,
		"amount":        order.PaymentInfo.Amount,
		"currency":      order.PaymentInfo.Currency,
		"paymentDate":   order.PaymentInfo.PaymentDate,
	}

	// Format order totals
	totals := map[string]interface{}{
		"subtotal": order.Totals.Subtotal,
		"shipping": order.Totals.Shipping,
		"discount": order.Totals.Discount,
		"tax":      order.Totals.Tax,
		"total":    order.Totals.Total,
	}

	// Create complete order data structure
	orderData := map[string]interface{}{
		"id":           order.Id,
		"userId":       order.UserId,
		"orderNumber":  order.OrderNumber,
		"status":       order.Status,
		"items":        items,
		"billingInfo":  billingInfo,
		"shippingInfo": shippingInfo,
		"paymentInfo":  paymentInfo,
		"totals":       totals,
		"notes":        order.Notes,
		"createdAt":    order.CreatedAt,
		"updatedAt":    order.UpdatedAt,
	}

	return orderData
}

// UpdateOrderStatus updates the status of an order
func (c *Config) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Get role from context - admin check done via middleware but double check
	role, _ := r.Context().Value("role").(string)
	if role != "admin" {
		util.ErrorJSON(w, errors.New("admin privileges required"), http.StatusForbidden)
		return
	}

	// Get order ID from URL
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		util.ErrorJSON(w, errors.New("order ID is required"), http.StatusBadRequest)
		return
	}

	// Parse request body
	var requestPayload struct {
		Status string `json:"status"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"paid":       true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
		"refunded":   true,
	}

	if !validStatuses[requestPayload.Status] {
		util.ErrorJSON(w, errors.New("invalid status"), http.StatusBadRequest)
		return
	}

	// Call checkout service via gRPC
	log.Printf("Updating order %s status to %s by admin %s", orderID, requestPayload.Status, userID)

	res, err := c.CheckoutClient.UpdateOrderStatus(r.Context(), &checkoutpb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  requestPayload.Status,
	})

	if err != nil {
		log.Printf("ERROR: Failed to update order status: %v", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   !res.Success,
		Message: res.Message,
		Data: map[string]interface{}{
			"status": requestPayload.Status,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}
