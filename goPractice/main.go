package main

import (
	"goProject/controllers"
	"goProject/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitialiseDb()
	router := gin.Default()

	router.GET("/data", controllers.GetData)                                      //Get all the data
	router.GET("/data/nonDeleted", controllers.GetNonDeletedData)                 //Get all the non-deleted data
	router.GET("/data/deleted", controllers.GetDeletedData)                       //Get all the deleted data
	router.GET("/data/:companyId", controllers.GetDataById)                       //Get data by companyId
	router.GET("/data/softDeleteData/:companyId", controllers.SoftDeleteDataById) //Soft delete data by id
	router.POST("/data/postPerson", controllers.PostData)                         //Add a new user
	// router.PUT("/putPerson", controllers.PutData)                          //Update data by PUT method (Add a new user if id does not exist)
	router.PATCH("/data/updatePerson/:companyId", controllers.UpdateDataById)  //Update data by id
	router.DELETE("/data/deletePerson/:companyId", controllers.DeleteDataById) //Delete data by id
	router.Run(":8080")

	// fmt.Println("Hello, World!")

}
