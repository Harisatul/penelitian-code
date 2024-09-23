package order

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"log/slog"
	"net/http"
	"order-svc/feature/category"
	"order-svc/feature/shared"
	"order-svc/pkg"
)

func HttpRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/orders", createOrderHandler)
	mux.HandleFunc("POST /api/payments/notify", paymentNotificationHandler)
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState3       = shared.LogEventStateFetchDB
		lfState3Status = "state_3_check_user_has_order_status"

		lvState4       = shared.LogEventStateCreatePayment
		lfState4Status = "state_4_create_payment_status"

		lvState6       = shared.LogEventStateInsertDB
		lfState6Status = "state_6_insert_order_status"

		lvState7       = shared.LogEventStateKafkaPublish
		lfState7Status = "state_7_publish_message_status"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("CreateOrder"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "invalid request", err, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(req),
	)

	/*------------------------------------
	| Step 3 : Check if user has order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	isExist, err := isUserHasOrder(ctx, req.Email)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	if isExist {
		shared.WriteErrorResponse(w, http.StatusConflict, errOrderExist)
		return
	}

	lock, err := cache.SetNX(ctx, req.Email, true, orderCancellationDuration).Result()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	if !lock {
		shared.WriteErrorResponse(w, http.StatusConflict, errOrderExist)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	/*------------------------------------
	| Step 4 : Create payment
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState4))

	orderId, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState4Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	price := category.Categories[req.CategoryID].Price
	vaCode, err := createVirtualAccountPayment(ctx, orderId, price)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState4Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState4Status))

	/*------------------------------------
	| Step 6 : Insert order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState6))

	order := orderEntity{
		ID:         orderId,
		CategoryID: req.CategoryID,
		Email:      req.Email,
	}

	ticketId, version, err := insertOrder(ctx, order)
	if err != nil {
		if err == errTicketNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}

		lf = append(lf, pkg.LogStatusFailed(lfState6Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState6Status))

	/*------------------------------------
	| Step 7 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState7))

	message := createOrderMessage{
		ID: orderId,
		Ticket: createOrderTicketPayload{
			ID:         ticketId,
			CategoryID: req.CategoryID,
			Version:    version,
		},
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState7Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	err = pkg.PublishMessage(kp, CreateOrderTopic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState7Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState7Status))

	shared.WriteSuccessResponse(w, http.StatusCreated,
		createOrderResponse{
			ID:     orderId,
			Total:  price,
			VaCode: vaCode,
		},
	)

	pkg.LogInfoWithContext(ctx, "success create order", lf)
}

func paymentNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateValidateRequest
		lfState2Status = "state_2_validate_payment_status"

		lvState3       = shared.LogEventStateUpdateDB
		lfState3Status = "state_2_update_order_status"

		lvState4       = shared.LogEventStateKafkaPublish
		lfState4Status = "state_4_publish_message_status"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("PaymentNotification"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var req paymentNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "invalid request", err, lf)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(req),
	)

	/*------------------------------------
	| Step 2 : Validate payment
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	if req.StatusCode != "200" || req.TransactionStatus != "settlement" {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx, "invalid payment status", nil, lf)
		shared.WriteErrorResponse(w, http.StatusBadRequest, errors.New("invalid payment status"))
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Update order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	orderId := cast.ToUint64(req.OrderID)
	err := updateOrderStatusToSuccess(ctx, orderId)
	if err != nil {
		if err == errOrderNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}

		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	/*------------------------------------
	| Step 3 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState4))

	message := completeOrderMessage{
		ID: orderId,
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState4Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	err = pkg.PublishMessage(kp, CompleteOrderTopic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState4Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState4Status))

	shared.WriteSuccessResponse(w, http.StatusOK, nil)
	pkg.LogInfoWithContext(ctx, "success update order status", lf)
}
