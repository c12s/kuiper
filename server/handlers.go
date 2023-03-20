package server

import (
	"errors"
	"kuiper/service"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ch configHandler) SaveConfig(c *gin.Context) {
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
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error:": "Only application/json is accepted"})
		return
	}

	rt, err := decodeConfigBody(c.Request.Body)
	if err != nil || rt.Entries == nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Invalid JSON"})
		return
	}

	// if ch.idempotencyService.FindRequestId(ctx, requestId) == true {
	// 	http.Error(c.Writer, "Request has been already sent", http.StatusBadRequest)
	// 	return
	// }

	cid, err := ch.configService.CreateConfig(ctx, rt)
	if err == service.NoVersionError {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "No version supplied"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Error when saving config"})
		return
	}

	// reqId := ""
	// if err == nil {
	// 	reqId = ts.idempotencyService.SaveRequestId(ctx)
	// }

	c.JSON(http.StatusOK, gin.H{"id": cid})
}

func (ch configHandler) GetConfig(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.CreateConfig")
	defer span.End()

	id := c.Param("id")
	ver := c.Param("ver")

	cfg, err := ch.configService.GetConfig(ctx, id, ver)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error:": "No value under key"})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

func (ch configHandler) CreateNewVersion(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.CreateNewVersion")
	defer span.End()

	id := c.Param("id")

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
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error:": "Only application/json is accepted"})
		return
	}

	rt, err := decodeConfigBody(c.Request.Body)
	if err != nil || rt.Entries == nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Invalid JSON"})
		return
	}

	err = ch.configService.CreateNewVersion(ctx, rt, id)
	if err == service.NoVersionError {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "No version supplied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "Invalid JSON"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
	return
}
