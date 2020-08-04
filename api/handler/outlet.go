package handler

import (
	"context"
	"net/http"
	"strconv"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
)

//CreateMerchantOutlet _
func (s *Service) CreateMerchantOutlet(c *gin.Context) {
	var outlet pbM.OutletRequest
	if err := c.ShouldBindJSON(&outlet); err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to bind Outlet", e))
		return
	}

	response, err := ClientMerchant.CreateMerchantOutlet(context.TODO(), &outlet)
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call CreateMerchantOutlet service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Created, "successful"))
}

//GetMerchantOutlets _
func (s *Service) GetMerchantOutlets(c *gin.Context) {
	merchantID, err := strconv.ParseUint(c.Param("merchant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginH("Failed to parse string to uint64", err))
		return
	}
	response, err := ClientMerchant.GetMerchantOutlets(context.TODO(), &pbM.GetOutletIdRequest{
		MerchantId: merchantID,
	})

	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call GetMerchantOutlets service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.GetOutlets(), "successful"))
}

//UpdateMerchantOutlet _
func (s *Service) UpdateMerchantOutlet(c *gin.Context) {
	var outlet pbM.UpdateOutletRequest
	if err := c.ShouldBindJSON(&outlet); err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to bind Outlet", e))
		return
	}

	response, err := ClientMerchant.UpdateMerchantOutlet(context.TODO(), &outlet)
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call UpdateMerchantOutlet service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Created, "successful"))
}

//DeleteMerchantOutlet _
func (s *Service) DeleteMerchantOutlet(c *gin.Context) {
	var (
		merchantID, outletID uint64
		err                  error
	)
	merchantID, err = strconv.ParseUint(c.Query("merchant_id"), 10, 64)
	outletID, err = strconv.ParseUint(c.Query("outlet_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginH("Failed to pass uint64", err))
		return
	}
	rsp, err := ClientMerchant.DeleteMerchantOutlet(context.TODO(), &pbM.GetOutletIdRequest{
		Id:         outletID,
		MerchantId: merchantID,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call DeleteMerchantOutlet", e))
		return
	}
	c.JSON(http.StatusOK, ginH(rsp.Deleted, "successful"))
}
