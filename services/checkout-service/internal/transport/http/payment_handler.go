package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"checkout-service/internal/domain"
	"checkout-service/internal/event"
)

// PaymentHandler handles HTTP requests for payment operations
type PaymentHandler struct {
	orderService domain.OrderService
	orderRepo    domain.OrderRepository
	eventEmitter *event.Emitter
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(orderService domain.OrderService, orderRepo domain.OrderRepository, eventEmitter *event.Emitter) *PaymentHandler {
	return &PaymentHandler{
		orderService: orderService,
		orderRepo:    orderRepo,
		eventEmitter: eventEmitter,
	}
}

// HandleMomoCallback handles the callback from MoMo payment gateway
// URL: POST /api/v1/payments/momo/callback
func (h *PaymentHandler) HandleMomoCallback(w http.ResponseWriter, r *http.Request) {
	// Parse the callback parameters
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract relevant parameters
	orderId := r.FormValue("orderId")
	resultCode := r.FormValue("resultCode")
	message := r.FormValue("message")
	transId := r.FormValue("transId")

	log.Printf("MoMo callback received: orderId=%s, resultCode=%s, message=%s, transId=%s",
		orderId, resultCode, message, transId)

	// Validate the required parameters
	if orderId == "" {
		http.Error(w, "Missing orderId parameter", http.StatusBadRequest)
		return
	}

	// Check if the payment was successful
	// MoMo returns resultCode=0 for successful transactions
	if resultCode != "0" {
		log.Printf("MoMo payment failed: %s", message)

		// Return success to MoMo to acknowledge receipt of callback
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Payment verification failed",
		})
		return
	}

	// Get the order from database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get order using the repository
	order, err := h.orderRepo.GetOrderByID(ctx, orderId)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Update payment info in the order
	paymentInfo := domain.PaymentInfo{
		PaymentMethod: "momo",
		TransactionID: transId,
		Status:        "completed",        // Use consistent status name
		Amount:        order.Totals.Total, // Use the total from the order
		Currency:      "VND",
		PaymentDate:   time.Now().Format(time.RFC3339),
	}

	// Update payment info in database
	err = h.orderRepo.UpdatePaymentInfo(ctx, orderId, paymentInfo)
	if err != nil {
		log.Printf("Failed to update payment info: %v", err)
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	// Update order status to processing
	err = h.orderRepo.UpdateOrderStatus(ctx, orderId, "processing")
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		// Continue anyway, payment was successful
	}

	// Emit payment success event
	if h.eventEmitter != nil {
		// Prepare event data
		eventData := event.PaymentEventData{
			OrderID:       order.ID,
			OrderNumber:   order.OrderNumber,
			UserID:        order.UserID,
			UserEmail:     order.BillingInfo.Email,
			PaymentMethod: "momo",
			Amount:        order.Totals.Total,
			Currency:      "VND",
			TransactionID: transId,
			Status:        "completed", // Use consistent status name
			PaymentDate:   time.Now().Format(time.RFC3339),
		}

		// Emit the event
		go func() {
			if err := h.eventEmitter.EmitPaymentSucceeded(eventData); err != nil {
				log.Printf("Failed to emit payment_succeeded event: %v", err)
			} else {
				log.Printf("Successfully emitted payment_succeeded event for order %s", orderId)
			}
		}()
	}

	// Return success to MoMo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Payment verified successfully",
	})
}
