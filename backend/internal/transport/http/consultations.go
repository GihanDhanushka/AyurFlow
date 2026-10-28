package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"ayurflow/backend/internal/consultations"

	"github.com/gin-gonic/gin"
)

type ConsultationDraftService interface {
	List(context.Context, string, string) ([]consultations.Summary, error)
	Create(context.Context, string, string, string) (consultations.Consultation, error)
	GetByID(context.Context, string, string) (consultations.Consultation, error)
	UpdateByID(context.Context, string, string, int64, string) (consultations.Consultation, error)
}

func registerConsultationDrafts(group *gin.RouterGroup, service ConsultationDraftService) {
	group.GET("/patients/:id/consultations", func(c *gin.Context) {
		items, err := service.List(c.Request.Context(), ActorID(c), c.Param("id"))
		if err != nil {
			consultationError(c, err, true)
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})
	group.POST("/patients/:id/consultations", func(c *gin.Context) {
		var input *struct {
			DraftText string `json:"draftText"`
		}
		if !consultationJSON(c, &input) || input == nil {
			badConsultationRequest(c)
			return
		}
		record, err := service.Create(c.Request.Context(), ActorID(c), c.Param("id"), input.DraftText)
		if err != nil {
			consultationError(c, err, true)
			return
		}
		c.JSON(http.StatusCreated, record)
	})
	group.GET("/consultations/:id", func(c *gin.Context) {
		record, err := service.GetByID(c.Request.Context(), ActorID(c), c.Param("id"))
		if err != nil {
			consultationError(c, err, true)
			return
		}
		c.JSON(http.StatusOK, record)
	})
	group.PATCH("/consultations/:id", func(c *gin.Context) {
		var input *struct {
			Version   int64  `json:"version"`
			DraftText string `json:"draftText"`
		}
		if !consultationJSON(c, &input) || input == nil {
			badConsultationRequest(c)
			return
		}
		record, err := service.UpdateByID(c.Request.Context(), ActorID(c), c.Param("id"), input.Version, input.DraftText)
		if err != nil {
			consultationError(c, err, true)
			return
		}
		c.JSON(http.StatusOK, record)
	})
}

func consultationJSON(c *gin.Context, target any) bool {
	mediaType, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return false
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 72*1024)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(target) != nil {
		return false
	}
	return decoder.Decode(new(any)) == io.EOF
}

func badConsultationRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, failure("bad_request", "Invalid consultation request."))
}

func consultationError(c *gin.Context, err error, hideRecord bool) {
	switch {
	case errors.Is(err, consultations.ErrNotFound), hideRecord && errors.Is(err, consultations.ErrForbidden):
		c.JSON(http.StatusNotFound, failure("not_found", "The requested resource was not found."))
	case errors.Is(err, consultations.ErrForbidden):
		c.JSON(http.StatusForbidden, failure("forbidden", "This action is not permitted."))
	case errors.Is(err, consultations.ErrInvalidDraft):
		c.JSON(http.StatusUnprocessableEntity, failure("validation_failed", "Check the consultation fields."))
	case errors.Is(err, consultations.ErrStaleVersion), errors.Is(err, consultations.ErrNotDraft):
		c.JSON(http.StatusConflict, failure("conflict", "This action conflicts with the current record state."))
	default:
		c.JSON(http.StatusServiceUnavailable, failure("unavailable", "The request could not be completed."))
	}
}
