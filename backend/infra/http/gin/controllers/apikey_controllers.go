package controllers

import (
	"net/http"

	apikeyapp "ai_hub.com/app/core/app/apikey/apikeyapp"
	"github.com/gin-gonic/gin"
)

type APIKeyController struct {
	createHandler     *apikeyapp.CreateAPIKeyHandler
	updateHandler     *apikeyapp.UpdateAPIKeyHandler
	deleteHandler     *apikeyapp.DeleteAPIKeyHandler
	getByIDHandler    *apikeyapp.GetKeyByIDHandler
	getByOwnerHandler *apikeyapp.GetKeysByOwnerHandler
}

func NewAPIKeyController(
	createH *apikeyapp.CreateAPIKeyHandler,
	updateH *apikeyapp.UpdateAPIKeyHandler,
	deleteH *apikeyapp.DeleteAPIKeyHandler,
	getByIDH *apikeyapp.GetKeyByIDHandler,
	getByOwnerH *apikeyapp.GetKeysByOwnerHandler,
) *APIKeyController {
	return &APIKeyController{
		createHandler:     createH,
		updateHandler:     updateH,
		deleteHandler:     deleteH,
		getByIDHandler:    getByIDH,
		getByOwnerHandler: getByOwnerH,
	}
}

// POST /api/keys
func (ctrl *APIKeyController) Create(c *gin.Context) {
	var req struct {
		KeyName   string   `json:"keyName"`
		KeyValue  string   `json:"keyValue"`
		Provider  string   `json:"provider"`
		ModelName string   `json:"modelName"`
		UsageEnv  *string  `json:"usageEnv,omitempty"`
		Status    *string  `json:"status,omitempty"`
		Balance   *float64 `json:"balance,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	ownerID := c.GetString("userID")

	res, err := ctrl.createHandler.Create(c, apikeyapp.CreateAPIKeyCommand{
		OwnerID:   ownerID,
		KeyName:   req.KeyName,
		KeyValue:  req.KeyValue,
		Provider:  req.Provider,
		ModelName: req.ModelName,
		UsageEnv:  req.UsageEnv,
		Status:    req.Status,
		Balance:   req.Balance,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PATCH /api/api-keys/:id
func (ctrl *APIKeyController) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		KeyName   *string  `json:"keyName,omitempty"`
		KeyValue  *string  `json:"keyValue,omitempty"`
		Status    *string  `json:"status,omitempty"`
		Provider  *string  `json:"provider,omitempty"`
		ModelName *string  `json:"modelName,omitempty"`
		UsageEnv  *string  `json:"usageEnv,omitempty"`
		Balance   *float64 `json:"balance,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	ownerID := c.GetString("userID")
	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	updated, err := ctrl.updateHandler.Update(c, apikeyapp.UpdateAPIKeyCommand{
		ID:        id,
		OwnerID:   ownerID,
		KeyName:   req.KeyName,
		KeyValue:  req.KeyValue,
		Status:    req.Status,
		Provider:  req.Provider,
		ModelName: req.ModelName,
		UsageEnv:  req.UsageEnv,
		Balance:   req.Balance,
	})
	if err != nil {
		c.Error(err)
		return
	}

	p := updated.ToPrimitives()
	c.JSON(http.StatusOK, gin.H{
		"_id":      p.ID,
		"keyName":  p.KeyName,
		"provider": p.Provider,
		"model":    p.ModelName,
		"usageEnv": p.UsageEnv,
		"status":   p.Status,
		"balance":  p.Balance,
	})
}

// DELETE /api/api-keys/:id
func (ctrl *APIKeyController) Delete(c *gin.Context) {
	id := c.Param("id")
	ownerID := c.GetString("userID")

	if err := ctrl.deleteHandler.Delete(c, apikeyapp.DeleteAPIKeyCommand{
		ID:      id,
		OwnerID: ownerID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /api/keys/:id
func (ctrl *APIKeyController) GetByID(c *gin.Context) {
	id := c.Param("id")
	ownerID := c.GetString("userID")

	res, err := ctrl.getByIDHandler.GetKeyByID(c, apikeyapp.GetKeyByIDQuery{
		ID:      id,
		OwnerID: ownerID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *APIKeyController) GetMyKeys(c *gin.Context) {
	ownerID := c.GetString("userID")

	res, err := ctrl.getByOwnerHandler.GetKeysByOwner(c, apikeyapp.GetKeysByOwnerQuery{
		OwnerID: ownerID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
