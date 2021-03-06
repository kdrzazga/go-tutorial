package main

import (
    "database/sql"
	"fmt"
	"log"
	"context"  
	"time"
	"strings"
)
	
const (  
    username = "root"
    password = "S3cret"
    hostname = "127.0.0.1:3306"
    dbname   = "ecommerce"
	table = "product"
	columns = "(product_id, product_name, product_price)"
	TIMEOUT = 5 * time.Second
)

type Timestamp *time.Time

func dsn(dbName string) string {  
    return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func openDbManager() *sql.DB{
	fmt.Printf("Opening db manager...\n")
    db, err := sql.Open("mysql", dsn(""))
    if err != nil {
		var msg = "Error %s when opening DB\n"
        log.Printf(msg, err)
		fmt.Printf(msg)
        return nil
    }

	fmt.Printf("Connection: " + dsn(dbname) + "\n")

	return db
}

func resetDb() {  //create if not exists, drop and recreate otherwise
    db, err := sql.Open("mysql", dsn(""))
    if err != nil {
        log.Printf("Error %s when opening DB\n", err)
        return
    }
    defer db.Close()

    ctx, cancelfunc := context.WithTimeout(context.Background(), TIMEOUT)
    defer cancelfunc()
	
	db.ExecContext(ctx, "DROP DATABASE " + dbname)
    res, err := db.ExecContext(ctx, "CREATE DATABASE " + dbname)
    if err != nil {
        log.Printf("Error %s when creating DB\n", err)
        return
    }
    no, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when fetching rows", err)
        return
    }
    log.Printf("rows affected %d\n", no)

    db.Close()
    db, err = sql.Open("mysql", dsn(dbname))
    if err != nil {
        log.Printf("Error %s when opening DB", err)
        return
    }
    defer db.Close()

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(time.Minute * 5)

    ctx, cancelfunc = context.WithTimeout(context.Background(), TIMEOUT)
    defer cancelfunc()
    err = db.PingContext(ctx)
    if err != nil {
        log.Printf("Errors %s pinging DB", err)
        return
    }
    log.Printf("Connected to DB %s successfully\n", dbname)
}

func connectToDb() *sql.DB{
	dbConn, errCon := dbConnection()

    if errCon != nil {
        log.Printf("Error %s when getting db connection\n", errCon)
        return nil
    }	
	fmt.Printf("Successfully connected to database\n")
	
	return dbConn
}

func dbConnection() (*sql.DB, error) {  
    db, err := sql.Open("mysql", dsn(""))
    if err != nil {
        log.Printf("Error %s when opening DB\n", err)
        return nil, err
    }

    ctx, cancelfunc := context.WithTimeout(context.Background(), TIMEOUT)
    defer cancelfunc()
	fmt.Printf("Dropping old database (if such exists) ...\n");	
	db.ExecContext(ctx, "DROP DATABASE IF EXISTS "  + dbname)
	
	fmt.Printf("Creating brand new database...");	
    res, err := db.ExecContext(ctx, "CREATE DATABASE " + dbname)
    if err != nil {
        log.Printf("Error %s when creating DB\n", err)
        return nil, err
    }
    no, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when fetching rows\n", err)
        return nil, err
    }
    log.Printf("rows affected %d\n", no)
    
    db, err = sql.Open("mysql", dsn(dbname))
    if err != nil {
        log.Printf("Error %s when opening DB\n", err)
        return nil, err
    }
	
    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(time.Minute * 5)

    ctx, cancelfunc = context.WithTimeout(context.Background(), TIMEOUT)
    defer cancelfunc()
    err = db.PingContext(ctx)
    if err != nil {
        log.Printf("Errors %s pinging DB", err)
        return nil, err
    }
    log.Printf("Connected to DB %s successfully\n", dbname)
    return db, nil
}

func createProductTable(db *sql.DB) error {  	
    query := `CREATE TABLE IF NOT EXISTS ` + table + `(product_id int primary key auto_increment, product_name text, 
        product_price int, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)`
		
    ctx, cancelfunc := context.WithTimeout(context.Background(), TIMEOUT)
	
    defer cancelfunc()
	
	log.Printf("Creating table %s...\n", dbname)
    res, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Error %s when creating %s table", err, table)
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when getting rows affected", err)
        return err
    }
    log.Printf("Rows affected when creating table: %d", rows)
    return nil
}

func fillTableProduct(db *sql.DB){		
	for i := 0; i < 51; i++{
		executeQuery(`INSERT INTO ` + table + columns + generateProductQuery(), db)
	}
}

func readTableContent(db *sql.DB){
	query := `SELECT * FROM ` + table

	var(
	    product_id int
	    product_name string
	    product_price int
	    created_at Timestamp
	    updated_at Timestamp
	)
	
	ctx, cancelfunc := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancelfunc()
	  
	rows, err := db.Query(query)
    if err != nil {
        log.Fatal("Error %s when reading from %s table", err, table)
        log.Fatal(ctx)
        return
    }

    log.Println()
    log.Println(strings.Repeat("-", 75))

    for rows.Next() {
        rows.Scan(&product_id, &product_name, &product_price, &created_at, &updated_at) //all fields need to be entered to map record from db

        log.Println("Product ID = ", product_id, " name = ", product_name, " created at ", created_at, " updated at ", updated_at)
    }

}

func executeQuery(query string, db *sql.DB){
	log.Printf(query + "\n")
	
	ctx, cancelfunc := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancelfunc()

	time.Sleep(78 * time.Millisecond)	
	res, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Error %s when filling table %s ", err, table)
        return
    }
	
	rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when getting rows affected", err)
        return
    }

    log.Printf("Entered number of rows: %d", rows)
}