package server

import (
	"errors"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ch configHandler) CreateConfig(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.CreateConfig")
	defer span.End()

	contentType := c.Request.Header.Get("Content-Type")
	// idemKey := c.Request.Header.Get("x-idempotency-key")

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		span.RecordError(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		span.RecordError(err)
		http.Error(c.Writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigBody(c.Request.Body)
	if err != nil || rt.Version == "" || rt.Entries == nil {
		span.RecordError(err)
		http.Error(c.Writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// if ch.idempotencyService.FindRequestId(ctx, requestId) == true {
	// 	http.Error(c.Writer, "Request has been already sent", http.StatusBadRequest)
	// 	return
	// }

	cid, err := ch.configService.CreateConfig(ctx, rt)

	// reqId := ""
	// if err == nil {
	// 	reqId = ts.idempotencyService.SaveRequestId(ctx)
	// }

	c.JSON(http.StatusOK, gin.H{"id": cid})
}
