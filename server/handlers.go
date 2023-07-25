package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"kuiper/service"
	"kuiper/store"
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
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Only application/json is accepted"})
		return
	}

	newCfg, err := decodeConfigBody(c.Request.Body)
	if err != nil || newCfg.Entries == nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON"})
		return
	}

	cid, err := ch.configService.CreateConfig(ctx, newCfg)
	if err == service.NoVersionError {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No version supplied"})
		return
	} else if err == store.KeyAlreadyExistsError {
		c.JSON(http.StatusConflict, gin.H{"error": "Version already exists for the service"})
		return
	} else if err == service.NoServiceNameError {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No service name supplied"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when saving config"})
		return
	}

	cfgJson, _ := json.Marshal(newCfg.Entries)
	ch.nats.Publish(natsKey(newCfg.Service), cfgJson)
	c.JSON(http.StatusOK, gin.H{"id": cid})
}

func (ch configHandler) GetConfig(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.GetConfig")
	defer span.End()

	id := c.Param("id")
	ver := c.Param("ver")

	cfg, err := ch.configService.GetConfig(ctx, id, ver)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No value under key"})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

func (ch configHandler) GetConfigsByService(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.GetConfigsByService")
	defer span.End()

	id := c.Param("id")
	if c.Query("latest") == "true" {
		cfg, err := ch.configService.GetLatestConfigByService(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No value under key"})
			return
		}
		c.JSON(http.StatusOK, cfg)
		return
	}
	cfgs, err := ch.configService.GetConfigsByService(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No value under key"})
		return
	}
	c.JSON(http.StatusOK, cfgs)

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
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Only application/json is accepted"})
		return
	}

	rt, err := decodeConfigBody(c.Request.Body)
	if err != nil || rt.Entries == nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON"})
		return
	}

	err = ch.configService.CreateNewVersion(ctx, rt, id)
	switch err {
	case service.NoVersionError:
		c.JSON(http.StatusBadRequest, gin.H{"error": "No version supplied"})
		return
	case store.KeyAlreadyExistsError:
		c.JSON(http.StatusConflict, gin.H{"error": "Version already exists"})
		return
	case store.ErrorNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration with given ID doesn't exist"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
	return
}

func (ch configHandler) DeleteConfig(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.DeleteConfig")
	defer span.End()

	id := c.Param("id")
	ver := c.Param("ver")

	cfg, err := ch.configService.DeleteConfig(ctx, id, ver)
	if err == store.ErrorNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration with the given ID and version doesn't exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failure when connecting to database"})
	}

	c.JSON(http.StatusOK, cfg)
}

func (ch configHandler) DeleteConfigsWithPrefix(c *gin.Context) {
	ctx, span := ch.tracer.Start(c.Request.Context(), "configServer.DeleteConfigsWithPrefix")
	defer span.End()

	id := c.Param("id")

	cfg, err := ch.configService.DeleteConfigsWithPrefix(ctx, id)
	if err == store.ErrorNotFound {
		errorMsg := fmt.Sprintf("No configurations found for %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failure when connecting to database"})
		return
	}

	c.JSON(http.StatusOK, cfg)
}
