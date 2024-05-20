package routes

// Path: routes/users.go
import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jalong4/stock-service-go/config"
	"github.com/jalong4/stock-service-go/models"
)

type Response struct {
	Summary  interface{} `json:"summary"`
	Holdings interface{} `json:"holdings"`
}

// GetAllHoldingsHandler handles requests to get all holdings
func GetAllHoldingsHandler(c *gin.Context) {
    holdings, err := config.GetAllHoldings()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve holdings"})
        return
    }

    // Calculating the total cost and organizing the summary
    totalCost := 0.0
    for _, holding := range holdings {
        totalCost += holding.TotalCost
    }

	// Round totalCost to two decimal places
	totalCost = math.Round(totalCost*100) / 100

    summary := map[string]interface{}{
        "found": len(holdings),
        "totalCost": totalCost,
    }

    response := Response {
        Summary: summary,
        Holdings: holdings,
    }

    c.JSON(http.StatusOK, response)
}

func AddHoldingsHandler(c *gin.Context) {
    var input map[string]interface{}

    // Bind JSON to a map to check for unknown fields
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
        return
    }

    // Check if ID is provided and warn it will be ignored
    if _, exists := input["id"]; exists {
        delete(input, "id")
        log.Println("ID field will be ignored")
    }

    // Define allowed fields
    allowedFields := map[string]bool{
        "ticker":    true,
        "quantity":  true,
        "totalCost": true,
        "account":   true,
    }

    // Check for unknown fields
    for key := range input {
        if !allowedFields[key] {
            c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unknown field: %s", key)})
            return
        }
    }

    // Re-marshal the map into JSON and then unmarshal into the Holding struct
    jsonBytes, err := json.Marshal(input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process input data"})
        return
    }

    var holding models.Holding
    if err := json.Unmarshal(jsonBytes, &holding); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
        return
    }

    // Add the holding to the database
    newId, err := config.AddHolding(holding)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add holding"})
        return
    }

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Successfully added holdings for ticker %s with ID %s", holding.Ticker, newId)})
}

// GetHoldingsByTickerHandler handles requests to get holdings by ticker
func GetHoldingsByTickerHandler(c *gin.Context) {
    ticker := c.Param("ticker")
    holdings, err := config.GetHoldingsByTicker(ticker)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve holdings"})
        return
    }

    if len(holdings) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No holdings found for the given ticker"})
        return
    }

    c.JSON(http.StatusOK, holdings)
}

// GetHoldingsByAccountHandler handles requests to get holdings by account regex pattern
func GetHoldingsByAccountHandler(c *gin.Context) {
    accountPattern := c.Param("account")
    holdings, err := config.GetHoldingsByAccount(accountPattern)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve holdings"})
        return
    }

    if len(holdings) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No holdings found for the given account pattern"})
        return
    }

    totalCost := 0.0
    for _, holding := range holdings {
        totalCost += holding.TotalCost
    }
    totalCost = math.Round(totalCost*100) / 100  // Round total cost to two decimal places

    summary := map[string]interface{}{
        "found": len(holdings),
        "totalCost": totalCost,
    }
    response := Response {
        Summary: summary,
        Holdings: holdings,
    }

    c.JSON(http.StatusOK, response)
}

func GetHoldingByIDHandler(c *gin.Context) {
	id := c.Param("_id")
	holding, err := config.GetHoldingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve holding"})
		return
	}

	c.JSON(http.StatusOK, holding)
}

func DeleteHoldingHandler(c *gin.Context) {
	id := c.Param("_id")

	holding, err := config.GetHoldingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find holding with ID: " + id})
		return
	}

	// Delete the holding from the database
	_, err = config.DeleteHoldingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete holding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Holdings for %s deleted successfully!", holding.Ticker)})
}

func UpdateHoldingHandler(c *gin.Context) {
	id := c.Param("_id")

	log.Printf("Updating Holding ID: %s", id)

	holding, err := config.GetHoldingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find holding with ID: " + id})
		return
	}

	log.Printf("Updating Holding: %s", holding.Ticker)

	var updatedHolding models.Holding
	if err = c.ShouldBindJSON(&updatedHolding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Update the holding in the database
	_, err = config.UpdateHoldingByID(id, updatedHolding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update holding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Holdings for %s updated successfully!", holding.Ticker)})

}