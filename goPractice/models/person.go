package models

type Person struct {
	CompanyId int    `json:"companyId"`
	Name      string `json:"name" binding:"required"`
	Age       int    `json:"age" binding:"required"`
	IsDeleted int8   `json:"isDeleted"`
}
