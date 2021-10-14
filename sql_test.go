package belajargolangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES('Sul', 'Kifli')"

	_, err := db.ExecContext(ctx, script)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new Customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"

	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		rows.Scan(&id , &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("Id : ", id)
		fmt.Println("Name : ", name)
	}

}

func TestQuerySqlComplex(t *testing.T){
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating , birth_date, married, created_at FROM customer"

	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var  createdAt time.Time
		var married bool

		rows.Scan(&id , &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("====================")
		fmt.Println("Id : ", id)
		fmt.Println("Name : ", name)
		if email.Valid {
			fmt.Println("Email : ", email.String)
		}
		fmt.Println("Balance : ", balance)
		fmt.Println("Rating : ", rating)
		if birthDate.Valid{
			fmt.Println("Birth Date : ", birthDate.Time)
		}
		fmt.Println("Married : ", married)
		fmt.Println("CreatedAt : ", createdAt)
	}
}
func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #" 
	password := "admin"

	script := "SELECT username FROM user WHERE username = '"+username+
	"' AND password = '"+password+"' LIMIT 1 "

	fmt.Println(script)

	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses login", username)
	}else{
		fmt.Println("Gagal Login")
	}

}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin" 
	password := "admin"

	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1 "

	fmt.Println(script)
	
	rows, err := db.QueryContext(ctx, script, username,password)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses login", username)
	}else{
		fmt.Println("Gagal Login")
	}

}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "sul" 
	password := "Kifli"
	script := "INSERT INTO user(username, password) VALUES(? , ?)"

	_, err := db.ExecContext(ctx, script, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new User")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "kifli@gmail.com" 
	comment := "belajar itu indah"
	script := "INSERT INTO comments(email, comment) VALUES(? , ?)"

	result, err := db.ExecContext(ctx, script, email, comment)

	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new Comment with id " , insertId )
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	script := "INSERT INTO comments(email, comment) VALUES(? , ?)"

	statement, err := db.PrepareContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "Kifli" + strconv.Itoa(i)+ "@gamil.com" 
		comment := "ini kommen ke" + strconv.Itoa(i)+ "untuk kursus 1" 

		result, err := statement.ExecContext(ctx, email, comment)

		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id ", id)
	}

}

func TestTransaction(t *testing.T) {
	db :=GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	// do transaction

	

	script := "INSERT INTO comments(email, comment) VALUES(? , ?)"

	for i := 0; i < 10; i++ {
		email := "Kifli" + strconv.Itoa(i)+ "@gamil.com" 
		comment := "ini kommen ke" + strconv.Itoa(i)+ "untuk kursus 1" 

		result, err := tx.ExecContext(ctx, script, email, comment)

		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id ", id)
	}

	// err = tx.Commit()
	err = tx.Rollback()

	if err != nil {
		panic(err)
	}
}