package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ZhigarinKirill/billing-api/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// New user registraion
func (h *Handler) createUser(c *gin.Context) {
	user := &model.User{}
	if err := c.BindJSON(user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	id, err := h.Services.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("creating user error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// Exising users view
func (h *Handler) getAllUsers(c *gin.Context) {
	users, err := h.Services.GetAllUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("getting users error: %s", err.Error()))
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "no users",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUserBillAmount(c *gin.Context) {
	// User ID from URL getting
	userID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	// User existing check
	_, err = h.Services.GetUserByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "user does not exist")
		return
	}

	var userBillAmount string

	res, err := h.Services.GetBillAmountByID(userID)
	if err != nil {
		if res == "" {
			newErrorResponse(c, http.StatusNotFound, "user does not have an account")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("getting user balance error: %s", err.Error()))
		}
		return
	}
	userBillAmount = res

	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": userBillAmount,
	})

}

// Deposit or withdrawing money
func (h *Handler) changeUserBillAmount(c *gin.Context) {

	// User ID from URL getting
	userID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	type Request struct {
		Amount float64 `json:"amount"`
	}

	// Request body decoding
	req := &Request{}
	if err = c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	// User existing check
	_, err = h.Services.GetUserByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "user does not exist")
		return
	}
	log.Info().Msg(fmt.Sprintf("request url on change account %s", c.Request.RequestURI))

	// Deposit and withdrawing check
	if strings.Contains(c.Request.RequestURI, "deposit") {
		log.Info().Msg("deposit")
		err = h.Services.DepositMoneyIntoAccount(userID, req.Amount, true)
	} else if strings.Contains(c.Request.RequestURI, "withdraw") {
		log.Info().Msg("withdraw")
		err = h.Services.WithdrawMoneyFromAccount(userID, req.Amount, true)
	}

	if err != nil {
		if strings.Contains(err.Error(), "the amount must not be negative") {
			newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("update user account error: %s", err.Error()))
		} else {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("update user account error: %s", err.Error()))
		}
		return
	}

	if strings.Contains(c.Request.RequestURI, "deposit") {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("%.2f were successfully deposited into user account", req.Amount),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%.2f were successfully withdrawn from user account", req.Amount),
	})
}

// Transfer between two users
func (h *Handler) transferMoneyBetweenUsers(c *gin.Context) {
	// User ID from URL getting
	fromUserID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	type Request struct {
		ToUserID int     `json:"to_user_id"`
		Amount   float64 `json:"amount"`
	}

	// Request body decoding
	req := &Request{}
	if err := c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}

	// User existing check
	if _, err = h.Services.GetUserByID(fromUserID); err != nil {
		newErrorResponse(c, http.StatusNotFound, "sender user does not exist")
		return
	}

	// User existing check
	if _, err = h.Services.GetUserByID(req.ToUserID); err != nil {
		newErrorResponse(c, http.StatusNotFound, "target user does not exist")
		return
	}

	// Same user check
	if fromUserID == req.ToUserID {
		newErrorResponse(c, http.StatusBadRequest, "users must be different")
		return
	}

	log.Info().Msg(fmt.Sprintf("from user %d to user %d amount %.2f", fromUserID, req.ToUserID, req.Amount))
	if err = h.Services.TransferMoneyBetweenUsers(fromUserID, req.ToUserID, req.Amount); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("transfer error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%.2f were successfully transferred", req.Amount),
	})
}

// User ID from URL getting
func (h *Handler) getUserID(c *gin.Context) (int, error) {
	strUserID := c.Param("user_id")
	intUserID, err := strconv.Atoi(strUserID)
	if err != nil {
		return 0, err
	}
	return int(intUserID), nil
}
