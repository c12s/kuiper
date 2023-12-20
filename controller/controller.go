package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/c12s/kuiper/model"
	"github.com/c12s/kuiper/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Controller struct {
	service service.ConfigService
	logger  *zap.Logger
}

func New(service service.ConfigService, logger *zap.Logger) *Controller {
	return &Controller{
		service,
		logger,
	}
}

func (controller Controller) Init(router *mux.Router) {
	router.HandleFunc("/create/{type}", controller.CreateVersion).Methods("POST")
	router.HandleFunc("/list", controller.ListVersions).Methods("GET")
	http.Handle("/", router)
	controller.logger.Info("Controller router endpoints initialized and handle run.")
}

func (controller Controller) CreateVersion(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	typeVar, ok := vars["type"]
	if !ok {
		controller.logger.Error("configType variable isn't presented in url")
		writer.Write([]byte("error in parsing type"))
	}

	configType := model.ConfigType(typeVar)

	if !(&configType).IsValid() {
		controller.logger.Error("invalid value for configType",
			zap.String("configType", string(configType)),
		)
		writer.Write([]byte("invalid config type"))
		return
	}

	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		controller.logger.Error("error in reading bytes from Body reader",
			zap.Error(err),
		)
		writer.Write([]byte("unable to read bytes from body"))
		return
	}

	var version model.Version
	err = json.Unmarshal(bytes, &version)
	if err != nil {
		controller.logger.Error("error in unmarshalling bytes to model.Version struct",
			zap.Error(err),
		)
		return
	}

	controller.logger.Info("version received", zap.Any("version", version))

	version, err = controller.service.CreateNewVersion(version)
	if err != nil {
		writer.Write([]byte("error in creating version"))
		return
	}

	var response any = version
	responseBytes, err := json.Marshal(response)
	if err != nil {
		controller.logger.Error("error in marshalling newly created version object",
			zap.Any("version", version),
		)
		writer.Write([]byte("error in returning newly created version"))
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responseBytes)
}

func (controller Controller) ListVersions(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()

	entityTypeRaw := query.Get("type")
	idRaw := query.Get("id")
	namespaceRaw := query.Get("namespace")
	appNameRaw := query.Get("appName")
	fromVersionRaw := query.Get("fromVersion")
	withFromRaw := query.Get("withFrom")
	toVersionRaw := query.Get("toVersion")
	withToRaw := query.Get("withTo")
	sortTypeRaw := query.Get("sortType")

	listRequest, err := model.ParseListRequest(
		entityTypeRaw,
		idRaw,
		namespaceRaw,
		appNameRaw,
		fromVersionRaw,
		withFromRaw,
		toVersionRaw,
		withToRaw,
		sortTypeRaw,
	)

	if err != nil {
		writer.Write([]byte(err.Error()))
	}

	controller.logger.Info("listRequest parsed", zap.Any("listRequest", listRequest))

	results, err := controller.service.ListVersions(listRequest)
	if err != nil {
		writer.Write([]byte("error in creating version"))
		return
	}

	responseBytes, err := json.Marshal(results)
	if err != nil {
		controller.logger.Error("error in marshalling results",
			zap.Any("results", results),
		)
		writer.Write([]byte("error in returning results"))
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responseBytes)
}
