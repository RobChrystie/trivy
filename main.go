package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"

	TrivyTypes "github.com/aquasecurity/trivy/pkg/types"
	"github.com/gin-gonic/gin"
)

type VulnDisplay struct {
	ID    string
	Title string
}

var router *gin.Engine
var tmplt *template.Template
var db *sql.DB

func main() {
	db, err := dbConnection()
	if err != nil {
		fmt.Printf("Error creating database: %s\n", err)
		return
	}
	defer db.Close()

	err = createScanTable(db)
	if err != nil {
		fmt.Printf("Error creating scan table: %s\n", err)
		return
	}

	plan, _ := os.ReadFile("./stuff/results.json")
	err = insert(db, string(plan))
	if err != nil {
		fmt.Printf("Error inserting scan: %s\n", err)
		return
	}

	scan, _ := selectScan(db, 1)

	fmt.Println("Starting server")
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")
	initializeRoutes()
	router.Run("localhost:8081")
}

func initializeRoutes() {
	router.GET("/home", handlePage)
	// EXAMPLE: get with ID - router.GET("/albums/:id", getAlbumByID)
	// EXAMPLE: post - router.POST("/albums", postAlbums)
}

func loadScan(path string) (disp []VulnDisplay) {
	var scan TrivyTypes.Report

	plan, _ := os.ReadFile(path)

	if err := json.Unmarshal(plan, &scan); err != nil {
		log.Fatal(err)
		return
	}

	for _, result := range scan.Results {
		for _, vuln := range result.Vulnerabilities {
			var vulnDisp VulnDisplay
			vulnDisp.ID = vuln.VulnerabilityID
			vulnDisp.Title = vuln.Title
			disp = append(disp, vulnDisp)
		}
	}

	return
}

func handlePage(c *gin.Context) {
	c.HTML(200, "template.html", loadScan("stuff/results.json"))
}
