package controllers

import (
	"fmt"
	"goProject/database"
	"goProject/helpers"
	"goProject/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get all the data from the db and returns it to the user
func GetData(c *gin.Context) {
	database.ConnectToDb()
	defer database.CloseDb()
	rows, err := database.Db.Query(`SELECT * FROM users`)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while fetching data from db"})
		log.Fatal(err)
		return
	}
	defer rows.Close()
	tempPersons := []models.Person{}
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.CompanyId, &person.Name, &person.Age, &person.IsDeleted)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while scanning data from db"})
			log.Fatal(err)
			return
		}
		tempPersons = append(tempPersons, person)
	}
	//Edge case when no data is present in the db
	if len(tempPersons) == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"status": "204", "message": "No data found"})
		return
	}
	c.IndentedJSON(200, gin.H{"status": "200", "message": "Success!", "persons": tempPersons})
}

// Get all the data from the db where isDeleted = 0 (false) and return it to the user
func GetNonDeletedData(c *gin.Context) {
	database.ConnectToDb()
	defer database.CloseDb()
	rows, err := database.Db.Query(`SELECT * FROM users WHERE isDeleted = 0`)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while fetching data from db"})
		log.Fatal(err)
		return
	}
	defer rows.Close()
	tempPersons := []models.Person{}
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.CompanyId, &person.Name, &person.Age, &person.IsDeleted)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while scanning data from db"})
			log.Fatal(err)
			return
		}
		tempPersons = append(tempPersons, person)
	}
	//Edge case when no data is present in the db
	if len(tempPersons) == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"status": "204", "message": "No non-deleted data found"})
		return
	}
	c.IndentedJSON(200, gin.H{"status": "200", "message": "Success!", "persons": tempPersons})
}

// Get all the data from the db where isDeleted = 1(true) and return it to the user
func GetDeletedData(c *gin.Context) {
	database.ConnectToDb() // Initiating connection before processing the request
	defer database.CloseDb()
	rows, err := database.Db.Query(`SELECT * FROM users WHERE isDeleted = 1`)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while fetching data from db"})
		log.Fatal(err)
		return
	}
	defer rows.Close()
	tempPersons := []models.Person{}
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.CompanyId, &person.Name, &person.Age, &person.IsDeleted)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while scanning data from db"})
			log.Fatal(err)
			return
		}
		tempPersons = append(tempPersons, person)
	}
	//Edge case when no data is present in the db
	if len(tempPersons) == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"status": "204", "message": "No deleted data found"})
		return
	}
	c.IndentedJSON(200, gin.H{"message": "Success!", "persons": tempPersons})
}

// Get data by id from the db and return it to the user
func GetDataById(c *gin.Context) {
	database.ConnectToDb()
	defer database.CloseDb()
	id := c.Param("companyId") //Get the id from the url
	rows, err := database.Db.Query(`SELECT * FROM users WHERE companyId = ?`, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while fetching data from db"})
		log.Fatal(err)
		return
	}
	defer rows.Close()
	// tempPersons := []models.Person{}

	var person models.Person
	rowPresent := false
	for rows.Next() {
		rowPresent = true
		err := rows.Scan(&person.CompanyId, &person.Name, &person.Age, &person.IsDeleted)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while scanning data from db"})
			log.Fatal(err)
			return
		}
	}
	//Edge case when no data is present in the db
	if !rowPresent {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "400", "message": fmt.Sprintf("No data found with the given id: %s", id)})
		return
	}
	c.IndentedJSON(200, gin.H{"status": "200", "message": "Success!", "person": person})

}

// Soft delete data by id from the db and return it to the user
func SoftDeleteDataById(c *gin.Context) {
	database.ConnectToDb()
	defer database.CloseDb()
	id := c.Param("companyId") //Get the id from the url
	// Convert the id to int
	idInt, errr := strconv.Atoi(id)
	if errr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while converting id to int"})
		log.Fatal(errr)
		return
	}
	// Logic to soft delete the data
	stmt, err := database.Db.Prepare(`UPDATE users SET isDeleted = 1 WHERE companyId = ?`)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error preparing the query"})
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(idInt) // Execute the query
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while deleting data from db"})
		log.Fatal(err)
		return
	}
	//If no data is present with the given id then return a message
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	if rowsAffected == 0 {
		c.IndentedJSON(400, gin.H{"status": "400", "message": fmt.Sprintf("No data found with the given id: %s", id)})
		return
	}
	c.IndentedJSON(200, gin.H{"message": "Success!", "person": fmt.Sprintf("Data with id: %s deleted successfully", id)})
}

// Checks if the id is present in the db or not, if not then it inserts the data into the db
func PostData(c *gin.Context) {
	database.ConnectToDb()
	var newPerson models.Person
	if err := c.ShouldBindJSON(&newPerson); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Handle error to check if the same id is not present in db
	//Process -> Get Data(id's) from Db -> Iterate through the rows -> If matching id found -> throw error
	message, flag := helpers.GetDataAndCheckId(newPerson.CompanyId)
	if flag {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "Failed", "error": message + "! , Please try again with a different id"})
		return
	}
	defer database.CloseDb()
	_, err := database.Db.Exec(`INSERT INTO users (companyId, name, age, isDeleted) VALUES (?, ?, ?, ?)`, newPerson.CompanyId, newPerson.Name, newPerson.Age, newPerson.IsDeleted)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while inserting data into db"})
		log.Fatal(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Success!", "person": newPerson})
}

// Update data by id from the db if present and return it to the user
func UpdateDataById(c *gin.Context) {
	database.ConnectToDb()
	defer database.CloseDb()
	id := c.Param("companyId") //Get the id from the url
	// Convert the id to int
	idInt, errr := strconv.Atoi(id)
	if errr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while converting id to int"})
		log.Fatal(errr)
		return
	}
	// Logic to update the data
	stmt, err := database.Db.Prepare(`UPDATE users SET name = ?, age = ? WHERE companyId = ?`)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error preparing the query"})
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	var person models.Person
	err = c.BindJSON(&person)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while binding the request body"})
		log.Fatal(err)
		return
	}
	res, err := stmt.Exec(person.Name, person.Age, idInt) // Execute the query
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while updating data from db"})
		log.Fatal(err)
		return
	}
	//If no data is present with the given id then return a message
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	if rowsAffected == 0 {
		c.IndentedJSON(400, gin.H{"status": "400", "message": fmt.Sprintf("No data found with the given id: %s", id)})
		return
	}
	c.IndentedJSON(200, gin.H{"message": "Success!", "person": fmt.Sprintf("Data with id: %s updated successfully", id)})
}

// Delete data by id from the db and return it to the user
func DeleteDataById(c *gin.Context) {
	database.ConnectToDb()
	defer database.CloseDb()
	id := c.Param("companyId") //Get the id from the url
	// Convert the id to int
	idInt, errr := strconv.Atoi(id)
	if errr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while converting id to int"})
		log.Fatal(errr)
		return
	}
	// Logic to hard delete the data
	stmt, err := database.Db.Prepare(`DELETE FROM users WHERE companyId = ?`)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error preparing the query"})
		log.Fatal(err)
		return
	}

	defer stmt.Close()
	res, err := stmt.Exec(idInt) // Execute the query
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error while deleting data from db"})
		log.Fatal(err)
		return
	}

	//If no data is present with the given id then return a message
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	if rowsAffected == 0 {
		c.IndentedJSON(400, gin.H{"status": "400", "message": fmt.Sprintf("No data found with the given id: %s", id)})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Success!", "person": fmt.Sprintf("Data with id: %s deleted successfully", id)})

}

// PUT request to update existing data or insert new data

//1. First query to UPDATE the data using company id.
/*   If no rows affected, that means no data is present with the given id.

2. Second query to INSERT the data if no data is present with the given id.

*/
