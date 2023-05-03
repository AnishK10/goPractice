package helpers

import (
	"goProject/database"
)

func GetDataAndCheckId(idFromPost int) (string, bool) {
	//Returns false if no id is found
	rows, err := database.Db.Query("SELECT companyId FROM users")
	if err != nil {
		return err.Error(), false
	}
	for rows.Next() {
		//create a pointer to the id
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err.Error(), true // Here false == nil
		}

		if idFromPost == id {
			return "Id already exists", true
		}
	}
	defer rows.Close()
	return "", false
}
