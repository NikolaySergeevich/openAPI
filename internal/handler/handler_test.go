package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openapi/internal/memstore"
	"openapi/pkg/api/objapi"
	"testing"

	"github.com/google/uuid"
	// "github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_PostObjects(t *testing.T) {
	t.Parallel()
	//для начала нужно инициализировать handler и memstore
	s := memstore.New()
	h := NewHandler(s)

	//нужно создать объект, который будем передавать в метод PostObjects 
	//и нужно его замаршалить
	id := uuid.New().String()
	obj := objapi.Object{
		Id: id,
		Name: "obj 1",
		Lat: 10.10,
		Lon: 20.20,
	} 
	encodedObj, err := json.Marshal(obj)
	require.NoError(t, err)//это вместо проверки ошибки классическим способом if err...

	//теперь нужно создать запрос:
	req, err := http.NewRequest(http.MethodPost, "/objects", bytes.NewReader(encodedObj))
	require.NoError(t, err)

	//теперь нужно создать recoder.
	rr := httptest.NewRecorder()

	h.PostObjects(rr, req)//вызываем метод, который проверяем. Тепррь нужно проверить что он там сделал.
	//результат запроса можем получить из rr
	//проверим код ответа, который возвращает запрос:
	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)

	//теперь нужно проверить как метод обработал объект, который мы создали выше и передали в ручку
	//По логике ручка кладёт переданый объект в хранилище store. Мы вытащим из хранилища этот объект
	//и сравним его с тем, что мы тут создали выше:
	item, ok := s.FindByObjectID(id)
	assert.True(t, ok)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, obj.Name, item.Name)
	assert.Equal(t, obj.Lat, item.Lat)
	assert.Equal(t, obj.Lon, item.Lon)
}

func TestGetObjects(t *testing.T) {
	mStore := memstore.New()
	handler := NewHandler(mStore)

	mStore.Add(memstore.Item{ID: "1", Name: "Object 1", Lat: 10.0, Lon: 20.0})
	mStore.Add(memstore.Item{ID: "2", Name: "Object 2", Lat: 30.0, Lon: 40.0})

	req, err := http.NewRequest("GET", "/objects", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handler.GetObjects(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var respObjects []objapi.Object
	err = json.NewDecoder(rr.Body).Decode(&respObjects)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(respObjects))
}


func TestPostObjects(t *testing.T) {
	mStore := memstore.New()
	handler := NewHandler(mStore)

	obj := objapi.Object{Id: "3", Name: "Object 3", Lat: 50.0, Lon: 60.0}
	body, _ := json.Marshal(obj)
	req, err := http.NewRequest("POST", "/objects", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.PostObjects(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	item, exists := mStore.FindByObjectID(obj.Id)
	assert.True(t, exists)
	assert.Equal(t, obj.Name, item.Name)
	assert.Equal(t, obj.Id, item.ID)
	assert.Equal(t, obj.Lat, item.Lat)
	assert.Equal(t, obj.Lon, item.Lon)
}

func TestGetObjectsObjectId(t *testing.T) {
	mStore := memstore.New()
	handler := NewHandler(mStore)

	item := memstore.Item{ID: "2", Name: "Object 2", Lat: 30.0, Lon: 40.0}
	mStore.Add(item)

	req, err := http.NewRequest("GET", "/objects/"+item.ID, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetObjectsObjectId(rr, req, item.ID)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseObject objapi.Object
	err = json.NewDecoder(rr.Body).Decode(&responseObject)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, responseObject.Id)
	assert.Equal(t, item.Name, responseObject.Name)
	assert.Equal(t, item.Lat, responseObject.Lat)
	assert.Equal(t, item.Lon, responseObject.Lon)

	nonExistentID := "999"
	req, err = http.NewRequest("GET", "/objects/"+nonExistentID, nil)
	assert.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.GetObjectsObjectId(rr, req, nonExistentID)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetObjectsObjectIdDistance(t *testing.T) {
	mStore := memstore.New()
	handler := NewHandler(mStore)

	testID := "3"
	expectedDistance := 650
	mStore.Add(memstore.Item{ID: testID, Name: "Object 3", Lat: 50.0, Lon: 60.0})

	req, err := http.NewRequest("GET", "/objects/"+testID+"/distance?lat=55.0&lon=65.0", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetObjectsObjectIdDistance(rr, req, testID, objapi.GetObjectsObjectIdDistanceParams{Lat: 55.0, Lon: 65.0})

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseDistance struct{ Distance float64 }
	err = json.NewDecoder(rr.Body).Decode(&responseDistance)
	assert.NoError(t, err)
	assert.NotZero(t, responseDistance.Distance)
	assert.Equal(t, expectedDistance, int(responseDistance.Distance))
}

