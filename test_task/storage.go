package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "24204141"
	dbname   = "expenses_db"
)

var idSlice = []int{}
var nameSlice = []string{}
var categorySlice = []string{}
var dateSlice = []string{}
var monthSlice = []string{}
var daySlice = []string{}

func ConnToDB() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckErr(err)

	return db
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func inputRequest(db *sql.DB) {
	var check_str, id int
	var flag bool = true
	for flag {
		printRequests()
		fmt.Scan(&check_str)
		switch check_str {
		case 1:
			user_id, user_data, date, category, cost := scanInsertValues()
			pushValuesToDB(db, user_id, user_data, date, category, cost)
		case 2:
			fmt.Println("Press number 1 if you wanna delete row,number 2 if all data")
			var del_check int
			fmt.Scan(&del_check)
			if del_check == 1 {
				fmt.Println("Input id[type int]")
				fmt.Scan(&id)
				if !CheckErrId(idSlice, id) {
					fmt.Println("[ Bad request ---> <<id not exist>!> ]")
					os.Exit(1)
				}
				deleteRowDB(db, id)
			}
			if del_check == 2 {
				deleteDataDB(db)
			}
		case 3:
			inputSelectRequest(db)
		case 4:
			flag = false
		default:
			fmt.Println("Wrong request number !")
		}
	}

}

func printRequests() {
	fmt.Println("_______________________________________")
	fmt.Println("|         |         |         |       |")
	fmt.Println("|insert(1)|delete(2)|select(3)|exit(4)|")
	fmt.Println("|         |         |         |       |")
	fmt.Println("_______________________________________")
	fmt.Print("INPUT REQUEST NUMBER --->  ")
}

func printSelect() {
	fmt.Println("____________________________________________________________")
	fmt.Println("|all_data(1)|per_day(2)|per_month(3)|per_year(4)|category(5)|")
	fmt.Println("____________________________________________________________")
	fmt.Print("INPUT REQUEST NUMBER --->  ")
}

func inputSelectRequest(db *sql.DB) {
	printSelect()
	var req int
	fmt.Scan(&req)
	switch req {
	case 1:
		getAllItemsDB(db)
	case 2:
		distributor(db, selectUser(db), 2)
	case 3:
		distributor(db, selectUser(db), 3)
	case 4:
		distributor(db, selectUser(db), 4)
	case 5:
		distributor(db, selectUser(db), 5)
	default:
		fmt.Println("Wrong request number !")
	}
}

func scanInsertValues() (int, string, string, string, int) {
	var user_id, cost int
	var user_data, date, category string
	fmt.Println("Input user_id[type int],user_data[type string],date[type string],category[type string],cost[type int]")
	fmt.Scan(&user_id, &user_data, &date, &category, &cost)
	if !CheckErrId(idSlice, user_id) {
		fmt.Println("[ Bad request --->  << id exist! >> ]")
		os.Exit(1)
	}
	return user_id, user_data, date, category, cost
}

func CheckErrId(slice []int, id int) bool {
	for _, val := range slice {
		if val == id {
			return false
		}
	}
	return true
}

func pushValuesToDB(db *sql.DB, user_id int, user_data string, date string, category string, cost int) {
	insertData := `insert into "expenses_info"("User_ID","User_data","Date","Category","Cost") values($1,$2,$3,$4,$5)`
	_, err := db.Exec(insertData, user_id, user_data, date, category, cost)
	CheckErr(err)
	CheckErrId(idSlice, user_id)
	idSlice = append(idSlice, user_id)
}

func deleteRowDB(db *sql.DB, id int) {
	deleteValue := `delete from "expenses_info" where "User_ID"=$1`
	_, err := db.Exec(deleteValue, id)
	CheckErr(err)
	for idx, val := range idSlice {
		if val == id {
			idSlice = append(idSlice[0:idx], idSlice[idx+1:]...)
		}
	}
}

func deleteDataDB(db *sql.DB) {
	deleteData := `truncate "expenses_info"`
	_, err := db.Exec(deleteData)
	CheckErr(err)
	idSlice = []int{}
}

func getAllItemsDB(db *sql.DB) {
	var user_id, cost int
	var user_data, date, category string

	rows, err := db.Query(`select * from "expenses_info"`)
	defer rows.Close()

	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&user_id, &user_data, &date, &category, &cost)
		CheckErr(err)
		fmt.Println(user_id, user_data, date, category, cost)
	}
}

func distributor(db *sql.DB, user string, check_request int) {

	switch check_request {
	case 2:
		//per_day
		selectPerDay(db, user)
	case 3:
		//per_month
		selectPerMonth(db, user)
	case 4:
		//per_year
		getDate(db, user)
		selectPerYear(db, user)
	case 5:
		//per_category
		selectPerCategory(db, user)
	default:
		fmt.Println("Bad request number!")
	}

}

func selectUser(db *sql.DB) string {
	rows, err := db.Query(`select distinct "User_data" from "expenses_info"`)
	defer rows.Close()
	CheckErr(err)

	var user_data, choice string
	fmt.Println("Input user surname >>> ")
	fmt.Print("Exsisting surnames: ")

	for rows.Next() {
		err := rows.Scan(&user_data)
		CheckErr(err)
		fmt.Print("[", user_data, "]", " ")
		nameSlice = append(nameSlice, user_data)
	}

	fmt.Println(";")
	fmt.Print(">")

	fmt.Scan(&choice)
	if !checkName(nameSlice, choice) {
		fmt.Println("Error: surname does not exist!")
		os.Exit(1)
	}
	return choice
}

func checkName(slice []string, name string) bool {
	for _, val := range slice {
		if val == name {
			return true
		}
	}
	return false
}

func selectCategory(db *sql.DB, user string) string {
	req_row := `select "Category" from "expenses_info" where "User_data"='` + user + "'"
	rows, err := db.Query(req_row)
	defer rows.Close()
	CheckErr(err)

	var category, choice string
	fmt.Println("Input category >>> ")
	fmt.Print("Exsisting categories: ")

	for rows.Next() {
		err := rows.Scan(&category)
		CheckErr(err)
		fmt.Print("[", category, "]", " ")
		categorySlice = append(categorySlice, category)
	}

	fmt.Println(";")
	fmt.Print(">")

	fmt.Scan(&choice)
	if !checkCategory(categorySlice, choice) {
		fmt.Println("Error: category does not exist!")
		os.Exit(1)
	}
	return choice
}

func checkCategory(slice []string, name string) bool {
	for _, val := range slice {
		if val == name {
			return true
		}
	}
	return false
}

func getDate(db *sql.DB, user string) {
	req_row := `select "Date" from "expenses_info" where "User_data"='` + user + "'"
	rows, err := db.Query(req_row)
	defer rows.Close()
	CheckErr(err)
	var date string
	for rows.Next() {
		err := rows.Scan(&date)
		CheckErr(err)
		dateSlice = append(dateSlice, date)
	}
	fmt.Println("Exsisting dates :", dateSlice)
}

func getDayFromDate() {

	for _, date_val := range dateSlice {
		for idx, val := range date_val {
			if string(val) == "." {
				date_val = date_val[:idx]
				break
			}
		}
		daySlice = append(daySlice, date_val)
	}

}

func getDayFromInput(inputDay string) string {
	var day string
	for idx, val := range inputDay {
		if string(val) == "." {
			day = inputDay[:idx]
			break
		}
	}
	return day
}

func checkDay(daySlice []string, inputDay string) bool {
	for _, val := range daySlice {
		if val == inputDay {
			return true
		}
	}
	return false
}

func getMonthFromDate() {
	for _, date_val := range dateSlice {
		for idx, val := range date_val {
			if string(val) == "." {
				date_val = date_val[idx+1 : idx+3]
				monthSlice = append(monthSlice, date_val)
				break
			}
		}
	}
}

func getMonthFromInput(inputMonth string) string {
	var month string
	for idx, val := range inputMonth {
		if string(val) == "." {
			month = inputMonth[idx+1 : idx+3]
			break
		}
	}
	return month
}

func checkMonth(monthSlice []string, inputMonth string) bool {
	for _, val := range monthSlice {
		if val == inputMonth {
			return true
		}
	}
	return false
}

func selectPerYear(db *sql.DB, user string) {
	var expensesPerYear int
	req_row := `select "Cost" from "expenses_info" where "User_data"='` + user + "'"
	rows, err := db.Query(req_row)
	defer rows.Close()

	CheckErr(err)
	var cost int
	for rows.Next() {
		err := rows.Scan(&cost)
		CheckErr(err)
		expensesPerYear += cost
		fmt.Println(cost)
	}
	fmt.Println("Expenses per year : ", expensesPerYear)
}

func selectPerMonth(db *sql.DB, user string) {
	getDate(db, user)
	getMonthFromDate()
	var inputMonth string
	fmt.Println("Input the date with the desired month [format day.month.year] >")
	fmt.Scan(&inputMonth)
	month := getMonthFromInput(inputMonth)
	if !checkMonth(monthSlice, month) {
		fmt.Println("Error: month does not exist!")
		os.Exit(1)
	}
	req_row := `select "Cost" from "expenses_info" where "Date"='` + inputMonth + "'"
	rows, err := db.Query(req_row)
	defer rows.Close()
	CheckErr(err)
	var cost, expensesPerMonth int
	for rows.Next() {
		err := rows.Scan(&cost)
		CheckErr(err)
		expensesPerMonth += cost
	}
	fmt.Println("Expenses per month : ", expensesPerMonth)
}

func selectPerDay(db *sql.DB, user string) {
	getDate(db, user)
	getDayFromDate()
	var inputDay string
	fmt.Println("Input day [format day.month.year] >")
	fmt.Scan(&inputDay)
	day := getDayFromInput(inputDay)
	if !checkDay(daySlice, day) {
		fmt.Println("Error: day does not exist!")
		os.Exit(1)
	}
	req_row := `select "Cost" from "expenses_info" where "Date"='` + inputDay + "'"
	rows, err := db.Query(req_row)
	defer rows.Close()
	CheckErr(err)
	var cost, expensesPerDay int
	for rows.Next() {
		err := rows.Scan(&cost)
		CheckErr(err)
		expensesPerDay += cost
	}
	fmt.Println("Expenses per day : ", expensesPerDay)
}

func selectPerCategory(db *sql.DB, user string) {
	category := selectCategory(db, user)
	var expensesPerCategory int
	req_row := `select "Cost" from "expenses_info" where "Category"='` + category + "'"
	rows, err := db.Query(req_row)
	defer rows.Close()

	CheckErr(err)
	var cost int
	for rows.Next() {
		err := rows.Scan(&cost)
		CheckErr(err)
		expensesPerCategory += cost

	}
	fmt.Println("Expenses per category : ", expensesPerCategory)
}
