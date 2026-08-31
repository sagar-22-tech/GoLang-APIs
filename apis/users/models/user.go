package models

type User struct {
	ID   int `json:"id"`
	Name struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	} `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name   string `json:"name"`
		Post   string `json:"post"`
		Salary string `json:"salary"`
	} `json:"company"`
}
