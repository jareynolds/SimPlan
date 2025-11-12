package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FleetWise Vehicle API Handlers

func createFleetWiseVehicle(c *gin.Context) {
	var vehicleConfig VehicleConfig
	if err := c.ShouldBindJSON(&vehicleConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	region := c.Query("region")
	if region == "" {
		region = "us-east-1" // Default region
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.CreateVehicle(vehicleConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Vehicle created successfully",
		"arn":     result.Arn,
		"vehicle": result.VehicleName,
	})
}

func batchCreateFleetWiseVehicles(c *gin.Context) {
	var req struct {
		Vehicles []VehicleConfig `json:"vehicles"`
		Region   string          `json:"region"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Region == "" {
		req.Region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(req.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdARNs, errors := client.BatchCreateVehicles(req.Vehicles)

	response := gin.H{
		"created_count": len(createdARNs),
		"created_arns":  createdARNs,
	}

	if len(errors) > 0 {
		errorMessages := make([]string, len(errors))
		for i, e := range errors {
			errorMessages[i] = e.Error()
		}
		response["errors"] = errorMessages
		c.JSON(http.StatusPartialContent, response)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func getFleetWiseVehicle(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.GetVehicle(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func updateFleetWiseVehicle(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	var updates VehicleConfig
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = client.UpdateVehicle(name, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle updated successfully"})
}

func deleteFleetWiseVehicle(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = client.DeleteVehicle(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted successfully"})
}

func getFleetWiseVehicleStatus(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.GetVehicleStatus(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func listFleetWiseVehicles(c *gin.Context) {
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	modelManifestARN := c.Query("model_manifest_arn")
	maxResults := int32(50) // Default

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	vehicles, err := client.ListVehicles(modelManifestARN, maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicles": vehicles,
		"count":    len(vehicles),
	})
}

// FleetWise Campaign API Handlers

func createFleetWiseCampaign(c *gin.Context) {
	var campaignConfig CampaignConfig
	if err := c.ShouldBindJSON(&campaignConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.CreateCampaign(campaignConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Campaign created successfully",
		"arn":      result.Arn,
		"campaign": result.Name,
	})
}

func getFleetWiseCampaign(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.GetCampaign(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func updateFleetWiseCampaign(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	var req struct {
		Action string `json:"action"` // APPROVE, SUSPEND, RESUME, UPDATE
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = client.UpdateCampaign(name, req.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign updated successfully"})
}

func deleteFleetWiseCampaign(c *gin.Context) {
	name := c.Param("name")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = client.DeleteCampaign(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}

func listFleetWiseCampaigns(c *gin.Context) {
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	maxResults := int32(50) // Default

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	campaigns, err := client.ListCampaigns(maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"campaigns": campaigns,
		"count":     len(campaigns),
	})
}

// FleetWise Fleet API Handlers

func createFleetWiseFleet(c *gin.Context) {
	var req struct {
		FleetID          string `json:"fleet_id" binding:"required"`
		Description      string `json:"description"`
		SignalCatalogARN string `json:"signal_catalog_arn" binding:"required"`
		Region           string `json:"region"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Region == "" {
		req.Region = "us-east-1"
	}

	client, err := NewAWSFleetWiseClient(req.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := client.CreateFleet(req.FleetID, req.Description, req.SignalCatalogARN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Fleet created successfully",
		"arn":     result.Arn,
		"fleet":   result.Id,
	})
}

func associateVehicleToFleet(c *gin.Context) {
	fleetID := c.Param("id")
	region := c.Query("region")
	if region == "" {
		region = "us-east-1"
	}

	var req struct {
		VehicleName string `json:"vehicle_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := NewAWSFleetWiseClient(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = client.AssociateVehicleToFleet(req.VehicleName, fleetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vehicle associated to fleet successfully",
		"vehicle": req.VehicleName,
		"fleet":   fleetID,
	})
}
