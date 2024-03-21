package handler

import (
	"net/http"

	"openapi/pkg/api/objapi"
)

var _ objapi.ServerInterface = (*Handler)(nil)

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct {
}

func (h Handler) GetObjects(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h Handler) PostObjects(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h Handler) DeleteObjectsObjectId(w http.ResponseWriter, r *http.Request, objectId string) {
	// TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h Handler) GetObjectsObjectId(w http.ResponseWriter, r *http.Request, objectId string) {
	// TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h Handler) GetObjectsObjectIdDistance(
	w http.ResponseWriter,
	r *http.Request,
	objectId string,
	params objapi.GetObjectsObjectIdDistanceParams,
) {
	// TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}
