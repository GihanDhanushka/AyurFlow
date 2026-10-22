package billing

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler translates HTTP transport concerns (parsing, status codes, JSON
// response shapes) into service calls. It does not implement authorization
// or business logic directly.
//
// Routes registered here must NOT be wired into the router until AUTH-03
// (role and record checks, CSRF middleware) is implemented by Member A.
// See RegisterRoutes for the registration point.
type Handler struct {
	svc *Service
}

// NewHandler creates a Handler backed by the given Service.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all billing routes under the /api/v1 group.
//
// AUTH-03 GATE: Do not call this function from the router until:
//   1. AUTH-03 role/record middleware is in place.
//   2. CSRF protection is applied to all mutations.
//   3. The caller passes a validated Session (actorID, roles) to each handler.
//
// The function exists now so that the route shapes are reviewable and so
// that integration can happen in a single place when AUTH-03 lands.
func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	g := api.Group("/billing")

	// Invoices
	g.POST("/invoices", h.createInvoice)
	g.GET("/invoices/:invoiceId", h.getInvoice)
	g.POST("/invoices/:invoiceId/items", h.addItem)
	g.DELETE("/invoices/:invoiceId/items/:itemId", h.removeItem)
	g.POST("/invoices/:invoiceId/issue", h.issueInvoice)
	g.POST("/invoices/:invoiceId/void", h.voidInvoice)

	// Receipts and balance
	g.POST("/invoices/:invoiceId/receipts", h.postReceipt)
	g.GET("/invoices/:invoiceId/receipts", h.listReceipts)
	g.GET("/invoices/:invoiceId/balance", h.getBalance)

	// Report
	g.GET("/reports/outstanding", h.outstandingReport)
}

// ---------------------------------------------------------------------------
// Invoice handlers
// ---------------------------------------------------------------------------

func (h *Handler) createInvoice(c *gin.Context) {
	// AUTH-03 TODO: extract actorID from validated session middleware.
	actorID := "TODO-auth03"

	var req CreateInvoiceRequest
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiFailure("bad_request", "Invalid request body."))
		return
	}

	inv, err := h.svc.CreateInvoice(c.Request.Context(), actorID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func (h *Handler) getInvoice(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	inv, err := h.svc.GetInvoice(c.Request.Context(), actorID, invoiceID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *Handler) addItem(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	var req AddItemRequest
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiFailure("bad_request", "Invalid request body."))
		return
	}

	item, err := h.svc.AddItem(c.Request.Context(), actorID, invoiceID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) removeItem(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")
	itemID := c.Param("itemId")

	if err := h.svc.RemoveItem(c.Request.Context(), actorID, invoiceID, itemID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) issueInvoice(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	inv, err := h.svc.IssueInvoice(c.Request.Context(), actorID, invoiceID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *Handler) voidInvoice(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	if err := h.svc.VoidInvoice(c.Request.Context(), actorID, invoiceID); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Receipt handlers
// ---------------------------------------------------------------------------

func (h *Handler) postReceipt(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	var req PostReceiptRequest
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiFailure("bad_request", "Invalid request body."))
		return
	}

	rec, err := h.svc.PostReceipt(c.Request.Context(), actorID, invoiceID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (h *Handler) listReceipts(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	recs, err := h.svc.ListReceipts(c.Request.Context(), actorID, invoiceID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	if recs == nil {
		recs = []Receipt{}
	}
	c.JSON(http.StatusOK, gin.H{"items": recs})
}

func (h *Handler) getBalance(c *gin.Context) {
	actorID := "TODO-auth03"
	invoiceID := c.Param("invoiceId")

	bal, err := h.svc.GetBalance(c.Request.Context(), actorID, invoiceID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, bal)
}

// ---------------------------------------------------------------------------
// Report handler
// ---------------------------------------------------------------------------

func (h *Handler) outstandingReport(c *gin.Context) {
	actorID := "TODO-auth03"

	report, err := h.svc.OutstandingReport(c.Request.Context(), actorID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, report)
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

// handleError maps billing sentinel errors to HTTP status codes using the
// shared error envelope format already established in the transport layer.
func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, apiFailure("not_found", "The requested record was not found."))
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, apiFailure("forbidden", "You do not have permission to perform this action."))
	case errors.Is(err, ErrConflict):
		msg := "This action conflicts with the current record state."
		if detail := extractDetail(err); detail != "" {
			msg = detail
		}
		c.JSON(http.StatusConflict, apiFailure("conflict", msg))
	case errors.Is(err, ErrValidation):
		msg := "Check the submitted fields."
		if detail := extractDetail(err); detail != "" {
			msg = detail
		}
		c.JSON(http.StatusUnprocessableEntity, apiFailure("validation_failed", msg))
	default:
		c.JSON(http.StatusServiceUnavailable, apiFailure("unavailable", "The request could not be completed."))
	}
}

func apiFailure(code, message string) gin.H {
	return gin.H{"error": gin.H{"code": code, "message": message}}
}

// extractDetail strips the sentinel prefix from a wrapped error message.
func extractDetail(err error) string {
	msg := err.Error()
	for _, prefix := range []string{"billing: state conflict: ", "billing: validation failed: "} {
		if strings.HasPrefix(msg, prefix) {
			return strings.TrimPrefix(msg, prefix)
		}
	}
	return ""
}
