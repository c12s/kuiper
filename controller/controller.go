package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/c12s/kuiper/errors"
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
	router.HandleFunc("/create", controller.CreateVersion).Methods("POST")
	router.HandleFunc("/list", controller.ListVersions).Methods("GET")
	router.HandleFunc("/diff/list", controller.ListVersionsDiff).Methods("GET")
	http.Handle("/", router)
	controller.logger.Info("Controller router endpoints initialized and handle run.")
}

func (controller Controller) CreateVersion(writer http.ResponseWriter, request *http.Request) {
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

	if !version.Type.IsValid() {
		controller.logger.Error("invalid value for configType",
			zap.String("configType", string(version.Type)),
		)
		writer.Write([]byte("invalid config type"))
		return
	}

	controller.logger.Info("version received", zap.Any("version", version))

	version, err = controller.service.CreateNewVersion(version)
	if err != nil {
		if err.Error() == errors.VersionAlreadyExist {
			writer.WriteHeader(http.StatusNotAcceptable)
			writer.Write([]byte(err.Error()))
			return
		} else if err.Error() == errors.VersionTagIsRequired {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
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
	writer.WriteHeader(http.StatusOK)
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
		return
	}

	controller.logger.Info("listRequest parsed", zap.Any("listRequest", listRequest))

	results, err := controller.service.ListVersions(listRequest)
	if err != nil {
		writer.Write([]byte("error in list version"))
		return
	}

	responseBytes, err := results.ToJSON()
	if err != nil {
		controller.logger.Error("error in marshalling results",
			zap.Any("results", results),
		)
		writer.Write([]byte("error in returning results"))
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseBytes)
}

func (controller Controller) ListVersionsDiff(writer http.ResponseWriter, request *http.Request) {

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
		return
	}

	controller.logger.Info("listRequest parsed", zap.Any("listRequest", listRequest))

	results, err := controller.service.ListVersionsDiff(listRequest)
	if err != nil {
		writer.Write([]byte("error in list versions diff"))
		return
	}

	responseBytes, err := results.ToJSON()
	if err != nil {
		controller.logger.Error("error in marshalling results",
			zap.Any("results", results),
		)
		writer.Write([]byte("error in returning results"))
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseBytes)
}
