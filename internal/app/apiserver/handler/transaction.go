package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
)

// Reserve money from main account into reserve account
func (h *Handler) reserveMoneyFromUser(c *gin.Context) {

	type Request struct {
		UserID    int     `json:"user_id"`
		ServiceID int     `json:"service_id"`
		OrderID   int     `json:"order_id"`
		Amount    float64 `json:"amount"`
	}

	// Request body decoding
	req := &Request{}
	if err := c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	// User existing check
	_, err := h.Services.GetUserByID(req.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "user does not exist")
		return
	}

	err = h.Services.ReserveMoneyFromAccount(req.UserID, req.ServiceID, req.OrderID, req.Amount)

	if err != nil {
		if strings.Contains(err.Error(), "the amount must be not negative") {
			newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("reserve money error: %s", err.Error()))
		} else {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("reserve money error: %s", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%.2f were successfully reserved from user", req.Amount),
	})
}

// Method for revenue recognition
func (h *Handler) completeReservationTransaction(c *gin.Context) {

	type Request struct {
		UserID    int     `json:"user_id"`
		ServiceID int     `json:"service_id"`
		OrderID   int     `json:"order_id"`
		Amount    float64 `json:"amount"`
	}

	// Request body decoding
	req := &Request{}
	if err := c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	// User existing check
	_, err := h.Services.GetUserByID(req.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "user does not exist")
		return
	}

	_, err = h.Services.CompleteReservationTransaction(req.UserID, req.ServiceID, req.OrderID, req.Amount)

	if err != nil {
		if strings.Contains(err.Error(), "the amount must be not negative") {
			newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("reservation completing error: %s", err.Error()))
		} else {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("reservation completing error: %s", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%.2f were successfully withdrawn from user", req.Amount),
	})
}

// Report generation
func (h *Handler) getReport(c *gin.Context) {

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid parametres: %s", err.Error()))
		return
	}
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid parametres: %s", err.Error()))
		return
	}

	reports, err := h.Services.GetReport(year, month)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("unable to format report: %s", err.Error()))
		return
	}

	// Save report to csv file
	csvContent, err := gocsv.MarshalString(&reports)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("unable to create csv: %s", err.Error()))
		return
	}
	csvFile, err := os.Create(fmt.Sprintf("./reports/report_%d_%d.csv", year, month))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("unable to create csv: %s", err.Error()))
		return
	}
	csvFile.WriteString(csvContent)
	csvFile.Close()

	c.Writer.Header().Set("Content-type", "application/csv")
	c.File(fmt.Sprintf("./reports/report_%d_%d.csv", year, month))
	// c.JSON(http.StatusOK, map[string]interface{}{
	// 	"report": reports,
	// 	"csv":    csvContent,
	// })
}
