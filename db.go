package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "password"
	hostname = "127.0.0.1:3306"
	dbname   = "trivy"
)

func dsn(dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func dbConnection() (*sql.DB, error) {
	fmt.Println(dsn(""))

	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		fmt.Printf("Error %s when creating DB\n", err)
		return nil, err
	}

	no, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Errpr %s when fetching rows\n", err)
		return nil, err
	}
	fmt.Printf("Rows affected: %d\n", no)

	db.Close()

	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelFunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Printf("Error %s pinging DB", err)
	}

	fmt.Printf("Connected to DB %s successfully\n", dbname)

	return db, nil
}

func createScanTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS scan(scan_id int primary key auto_increment, scan_json blob, 
        created_at datetime default CURRENT_TIMESTAMP)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when creating product table\n", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when getting rows affected\n", err)
		return err
	}
	fmt.Printf("Rows affected when creating table: %d\n", rows)

	return nil
}

func insert(db *sql.DB, scan string) error {
	query := "INSERT INTO scan(scan_json) VALUES (?)"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when preparing SQL statement\n", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, scan)
	if err != nil {
		fmt.Printf("Error %s when inserting row into scan table\n", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when finding rows affected\n", err)
		return err
	}
	fmt.Printf("%d scan created\n", rows)

	prdID, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error %s when getting last inserted scan\n", err)
		return err
	}

	fmt.Printf("Scan with ID %d created\n", prdID)

	return nil
}

func selectScan(db *sql.DB, scan_id int) (scan, error) {
	fmt.Printf("Getting scan\n")
	query := `select scan_json from scan where scan_id = ?`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("Error %s when preparing SQL statement\n", err)
		return "", err
	}
	defer stmt.Close()
	var scan string
	row := stmt.QueryRowContext(ctx, scan)
	if err := row.Scan(&scan); err != nil {
		return "", err
	}
	return scan, nil
}
