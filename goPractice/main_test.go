package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goProject/controllers"
	"goProject/database"
	"goProject/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// Test function to get all data irrespective of deleted or non-deleted
func TestGetData(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.GET("/data", controllers.GetData)
	req, _ := http.NewRequest("GET", "/data", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var users map[string]interface{}

	// fmt.Println(w.Body)
	json.Unmarshal(w.Body.Bytes(), &users)

	// fmt.Println(users)
	assert.NotEmpty(t, users["persons"])
	assert.Equal(t, http.StatusOK, w.Code)

}

// Test function to get non deleted data
func TestGetNonDeletedData(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.GET("/data/nonDeleted", controllers.GetNonDeletedData)
	req, _ := http.NewRequest("GET", "/data/nonDeleted", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var users map[string]int

	json.Unmarshal(w.Body.Bytes(), &users)

	// fmt.Println(users)

	//Check if status key in the response body is 200
	// assert.Equal(t, http.StatusOK, users["status"])

	assert.Equal(t, http.StatusOK, w.Code) //Checking only the status since if no nonDeletedData present then empty persons key.

}

// Test function to get deleted data
func TestGetDeletedData(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.GET("/data/deleted", controllers.GetDeletedData)
	req, _ := http.NewRequest("GET", "/data/deleted", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var users map[string]int

	json.Unmarshal(w.Body.Bytes(), &users)

	assert.Equal(t, http.StatusOK, w.Code) //Checking only the status since if no deletedData present then empty persons key.

}

// Test function to get data by id
func TestGetDataById(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.GET("/data/:companyId", controllers.GetDataById)
	reqFound, _ := http.NewRequest("GET", "/data/1020", nil)
	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, reqFound)

	// var users map[string]interface{}

	// fmt.Println(w.Body)
	// json.Unmarshal(w.Body.Bytes(), &users)

	// fmt.Println(users)

	//Check if status key in the response body is 200
	assert.Equal(t, http.StatusOK, w1.Code)

	reqNotFound, _ := http.NewRequest("GET", "/data/10", nil)
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, reqNotFound)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

}

// Test function to soft delete data by id
func TestSoftDeleteDataById(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.GET("/data/softDeleteData/:companyId", controllers.SoftDeleteDataById)
	reqFound, _ := http.NewRequest("GET", "/data/softDeleteData/1020", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, reqFound)
	assert.Equal(t, http.StatusOK, w1.Code) //Check if status key in the response body is 200

	reqNotFound, _ := http.NewRequest("GET", "/data/softDeleteData/10", nil)
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, reqNotFound)
	assert.Equal(t, http.StatusBadRequest, w2.Code) // Status = 400 if data not found

}

//Test function to create a new user

func TestPostData(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.POST("/data/postPerson", controllers.PostData)

	testInput := models.Person{
		CompanyId: 1048,
		Name:      "test 3",
		Age:       21,
		IsDeleted: 0,
	}
	jsonValue, _ := json.Marshal(testInput)
	fmt.Println(string(jsonValue))
	req, _ := http.NewRequest("POST", "/data/postPerson", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	database.DeleteRow(testInput.CompanyId)

}

//Test function to update a user

func TestUpdateDataById(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.PATCH("/data/updatePerson/:companyId", controllers.UpdateDataById)

	testInput := models.Person{
		Name: "test 4",
		Age:  22,
	}
	jsonValue, _ := json.Marshal(testInput)
	fmt.Println(string(jsonValue))
	reqFound, _ := http.NewRequest("PATCH", "/data/updatePerson/1025", bytes.NewBuffer(jsonValue))

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, reqFound)
	assert.Equal(t, http.StatusOK, w1.Code)

	reqNotFound, _ := http.NewRequest("PATCH", "/data/updatePerson/10", bytes.NewBuffer(jsonValue))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, reqNotFound)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

}

// Test function to delete data by id
func TestDeleteDataById(t *testing.T) {
	database.ConnectToDb()
	defer database.CloseDb()

	router := SetUpRouter()
	router.DELETE("/data/deletePerson/:companyId", controllers.DeleteDataById)
	reqFound, _ := http.NewRequest("DELETE", "/data/deletePerson/1022", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, reqFound)
	assert.Equal(t, http.StatusOK, w1.Code) //Check if status key in the response body is 200

	reqNotFound, _ := http.NewRequest("DELETE", "/data/deletePerson/10", nil)
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, reqNotFound)
	assert.Equal(t, http.StatusBadRequest, w2.Code) // Status = 400 if data not found

}
