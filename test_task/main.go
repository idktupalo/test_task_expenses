package main

func main() {
	db := ConnToDB()
	inputRequest(db)
	defer db.Close()
}
