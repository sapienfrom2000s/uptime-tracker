package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ctx context.Context
	db  *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	status := "up"
	if err = db.PingContext(ctx); err != nil {
		status = "down"
		log.Fatal(err)
	}
	fmt.Printf("DB is %q\n", status)

	// Checking if db tracker table is present
	var trackers string
	err = db.QueryRowContext(ctx, "SELECT name as trackers FROM sqlite_master WHERE type='table' AND name='Trackers';").
		Scan(&trackers)

	// table was not created initially
	if err == sql.ErrNoRows {
		fmt.Println("Trackers table isn't created yet")
		err = createTrackersTable()
		if err != nil {
			fmt.Println("Unable to create trackers table")
		}
	} else {
		fmt.Println("Trackers table already exists")
	}

	// Checking if db tracker table is present
	var uptimeStatuses string
	err = db.QueryRowContext(ctx, "SELECT name as uptimeStatuses FROM sqlite_master WHERE type='table' AND name='UptimeStatuses';").
		Scan(&uptimeStatuses)

	// table was not created initially
	if err == sql.ErrNoRows {
		fmt.Println("uptimeStatuses table isn't created yet")
		err = createUptimeStatuses()
		if err != nil {
			fmt.Println("Unable to create uptimeStatuses table")
		}
	} else {
		fmt.Println("uptimeStatuses table already exists")
	}

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/v1/trackers/status", func(c *gin.Context) {
		rows, err := db.Query("SELECT up from uptimeStatuses LIMIT 60")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var statuses []bool
		for rows.Next() {
			var up bool
			rows.Scan(&up)
			statuses = append(statuses, up)
		}
		c.JSON(200, statuses)
	})
	r.Run()
}

func createTrackersTable() error {
	insertTrackers := `CREATE TABLE Trackers (
    ID INTEGER,
    url TEXT
	);`

	fmt.Println("Creating Trackers Table...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, insertTrackers)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func createUptimeStatuses() error {
	insertUptimeStatuses := `CREATE TABLE UptimeStatuses(
		ID INTEGER,
		tracker_id INTEGER,
		up boolean
		)`

	fmt.Println("Creating UptimeStauses table...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, insertUptimeStatuses)
	if err != nil {
		return fmt.Errorf("Error: %q", err)
	}
	return nil
}
