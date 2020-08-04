package handler

import (
	"context"
	"net/http"
	"strconv"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
)

//CreateMerchantTerminal _
func (s *Service) CreateMerchantTerminal(c *gin.Context) {
	var terminal pbM.TerminalRequest
	if err := c.ShouldBindJSON(&terminal); err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to bind Terminal", e))
		return
	}

	response, err := ClientMerchant.CreateMerchantTerminal(context.TODO(), &terminal)
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call CreateMerchantTerminal service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Created, "successful"))
}

//GetMerchantTerminals _
func (s *Service) GetMerchantTerminals(c *gin.Context) {
	merchantID, err := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginH("Failed to parse string to uint64", err))
		return
	}
	response, err := ClientMerchant.GetMerchantTerminals(context.TODO(), &pbM.GetTerminalIdRequest{
		MerchantId: merchantID,
	})

	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call GetMerchantTerminals service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.GetTerminals(), "successful"))
}

//UpdateMerchantTerminal _
func (s *Service) UpdateMerchantTerminal(c *gin.Context) {
	var Terminal pbM.UpdateTerminalRequest
	if err := c.ShouldBindJSON(&Terminal); err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to bind Terminal", e))
		return
	}

	response, err := ClientMerchant.UpdateMerchantTerminal(context.TODO(), &Terminal)
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call UpdateMerchantTerminal service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Created, "successful"))
}

//DeleteMerchantTerminal _
func (s *Service) DeleteMerchantTerminal(c *gin.Context) {
	var (
		merchantID, TerminalID uint64
		err                    error
	)
	merchantID, err = strconv.ParseUint(c.Query("merchant_id"), 10, 64)
	TerminalID, err = strconv.ParseUint(c.Query("Terminal_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginH("Failed to pass uint64", err))
		return
	}
	rsp, err := ClientMerchant.DeleteMerchantTerminal(context.TODO(), &pbM.GetTerminalIdRequest{
		Id:         TerminalID,
		MerchantId: merchantID,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call DeleteMerchantTerminal", e))
		return
	}
	c.JSON(http.StatusOK, ginH(rsp.Deleted, "successful"))
}
