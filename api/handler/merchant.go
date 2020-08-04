package handler

import (
	"context"
	"net/http"
	"strconv"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
)

//Service holdes futher configurations
type Service struct{}

//ClientMerchant _
var ClientMerchant pbM.MerchantService

func ginH(msg, in interface{}) gin.H {
	switch in.(type) {
	case error:
		return gin.H{"Message": msg, "Error": in.(error).Error()}
	default:
		return gin.H{"Message": msg, "Success": in}
	}
}

//CreateMerchant restful api
func (s *Service) CreateMerchant(c *gin.Context) {
	merchant := pbM.CreateRequest{}

	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, ginH("failed to bind request", err))
		return
	}
	response, err := ClientMerchant.CreateMerchant(context.TODO(), &merchant)
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call merchant service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Created, "successful"))
}

//ListMerchants _
func (s *Service) ListMerchants(c *gin.Context) {
	response, err := ClientMerchant.GetMerchants(context.TODO(), &pbM.Request{})
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to fetch merchant list", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.GetMerchants(), "successful"))
}

//GetMerchantByID _
func (s *Service) GetMerchantByID(c *gin.Context) {
	merchantIDstr := c.Param("merchant_id")
	merchantID, err := strconv.ParseUint(merchantIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginH("Failed to convert str to uint", err))
		return
	}
	response, err := ClientMerchant.GetMerchantByID(context.TODO(), &pbM.GetIdRequest{
		Id: merchantID,
	})

	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to fetch merchant", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.GetMerchant(), "successful"))
}

//DeleteMerchantByID _
func (s *Service) DeleteMerchantByID(c *gin.Context) {
	merchantIDstr := c.Param("merchant_id")
	merchantID, err := strconv.ParseUint(merchantIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginH("Failed to convert str to uint", err))
		return
	}
	response, err := ClientMerchant.DeleteMerchantByID(context.TODO(), &pbM.GetIdRequest{
		Id: merchantID,
	})

	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusBadRequest, ginH("Failed to delete merchant", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Deleted, "successful"))
}

//UpdateMerchant update merchant
func (s *Service) UpdateMerchant(c *gin.Context) {
	var merchant pbM.UpdateRequest

	if err := c.ShouldBindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, ginH("failed to bind request", err))
		return
	}

	response, err := ClientMerchant.UpdateMerchant(context.TODO(), &merchant)
	if err != nil {
		e := errors.Parse(err.Error())
		c.JSON(http.StatusInternalServerError, ginH("Failed to call merchant service", e))
		return
	}
	c.JSON(http.StatusOK, ginH(response.Updated, "successful"))
}
