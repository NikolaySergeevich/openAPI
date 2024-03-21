package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"openapi/internal/geo"
	"openapi/internal/memstore"
	"openapi/pkg/api/objapi"
)

var _ objapi.ServerInterface = (*Handler)(nil)

func NewHandler(store store) *Handler {
	return &Handler{store: store}
}

type Handler struct {
	store store
}

func (h *Handler) GetObjects(w http.ResponseWriter, r *http.Request) {
	items := h.store.FindAll()
	respObjects := make([]objapi.Object, 0, len(items))
	for _, item := range items {
		respObjects = append(
			respObjects, objapi.Object{
				Id:   item.ID,
				Lat:  item.Lat,
				Lon:  item.Lon,
				Name: item.Name,
			},
		)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respObjects); err != nil {
		slog.Error("json.NewEncoder Encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostObjects(w http.ResponseWriter, r *http.Request) {
	var req objapi.Object
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("json.NewDecoder Decode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.store.Add(
		memstore.Item{
			ID:   req.Id,
			Name: req.Name,
			Lat:  req.Lat,
			Lon:  req.Lon,
		},
	)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteObjectsObjectId(w http.ResponseWriter, r *http.Request, objectId string) {
	h.store.DeleteByObjectID(objectId)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetObjectsObjectId(w http.ResponseWriter, r *http.Request, objectId string) {
	item, ok := h.store.FindByObjectID(objectId)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(
		objapi.Object{
			Id:   item.ID,
			Lat:  item.Lat,
			Lon:  item.Lon,
			Name: item.Name,
		},
	); err != nil {
		slog.Error("json.NewEncoder Encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetObjectsObjectIdDistance(
	w http.ResponseWriter,
	r *http.Request,
	objectId string,
	params objapi.GetObjectsObjectIdDistanceParams,
) {
	item, ok := h.store.FindByObjectID(objectId)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	type distObject struct {
		Distance float64
	}
	dist := geo.ComputeDistance(item.Lat, item.Lon, params.Lat, params.Lon)
	if err := json.NewEncoder(w).Encode(distObject{Distance: dist}); err != nil {
		slog.Error("json.NewEncoder Encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
