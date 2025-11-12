package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database Models
type Capability struct {
	ID           string   `gorm:"primaryKey" json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Enablers     string   `json:"enablers"` // JSON array stored as string
	Dependencies string   `json:"dependencies"` // JSON array stored as string
	CreatedAt    time.Time `json:"created_at"`
}

type Enabler struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Environment struct {
	ID                 string    `gorm:"primaryKey" json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Owner              string    `json:"owner"`
	Tags               string    `json:"tags"`
	Status             string    `json:"status"` // provisioning, running, stopped, error
	Capabilities       string    `json:"capabilities"` // JSON array
	EnablersConfig     string    `json:"enablers_config"` // JSON object
	ComputeConfig      string    `json:"compute_config"` // JSON object
	Storage            int       `json:"storage"`
	Network            string    `json:"network"`
	Priority           string    `json:"priority"`
	Duration           int       `json:"duration"`
	EstimatedCost      float64   `json:"estimated_cost"`
	ActualCost         float64   `json:"actual_cost"`
	Health             int       `json:"health"`
	Uptime             string    `json:"uptime"`
	FleetWiseConfig    string    `json:"fleetwise_config"` // JSON object for AWS FleetWise configuration
	UseRealAWSBackend  bool      `json:"use_real_aws_backend"` // Flag to use real AWS instead of simulation
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type StateTransition struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EnvironmentID string    `json:"environment_id"`
	FromState     string    `json:"from_state"`
	ToState       string    `json:"to_state"`
	Reason        string    `json:"reason"`
	Metadata      string    `json:"metadata"` // JSON object
	CreatedAt     time.Time `json:"created_at"`
}

type AuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EnvironmentID string    `json:"environment_id"`
	Action        string    `json:"action"`
	UserID        string    `json:"user_id"`
	Details       string    `json:"details"` // JSON object
	CreatedAt     time.Time `json:"created_at"`
}

type Upload struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EnvironmentID string    `json:"environment_id"`
	Filename      string    `json:"filename"`
	FileType      string    `json:"file_type"` // binary, config
	Version       string    `json:"version"`
	Size          int64     `json:"size"`
	Status        string    `json:"status"` // pending, processing, completed, failed
	CreatedAt     time.Time `json:"created_at"`
}

type Reservation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EnvironmentID string    `json:"environment_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Priority      string    `json:"priority"`
	Status        string    `json:"status"` // scheduled, active, completed, cancelled
	CreatedAt     time.Time `json:"created_at"`
}

// Request/Response DTOs
type CreateEnvironmentRequest struct {
	Name              string                 `json:"name" binding:"required"`
	Description       string                 `json:"description"`
	Owner             string                 `json:"owner"`
	Tags              string                 `json:"tags"`
	Capabilities      []string               `json:"capabilities"`
	EnablersConfig    map[string]interface{} `json:"enablers"`
	Compute           ComputeConfig          `json:"compute"`
	Storage           int                    `json:"storage"`
	Network           string                 `json:"network"`
	Priority          string                 `json:"priority"`
	Duration          int                    `json:"duration"`
	FleetWiseConfig   *FleetWiseConfig       `json:"fleetwise_config,omitempty"`
	UseRealAWSBackend bool                   `json:"use_real_aws_backend"`
}

type ComputeConfig struct {
	CPU       int `json:"cpu"`
	Memory    int `json:"memory"`
	Instances int `json:"instances"`
}

type ValidationResponse struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type CostEstimationResponse struct {
	DailyCost       float64            `json:"daily_cost"`
	MonthlyCost     float64            `json:"monthly_cost"`
	Breakdown       map[string]float64 `json:"breakdown"`
	OptimizationTip string             `json:"optimization_tip,omitempty"`
}

// Global DB instance
var db *gorm.DB

func main() {
	// Initialize database
	initDB()

	// Seed initial data
	seedData()

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API Routes
	v1 := router.Group("/api/v1")
	{
		// Capabilities and Enablers
		v1.GET("/capabilities", getCapabilities)
		v1.GET("/enablers", getEnablers)

		// Environments
		v1.GET("/environments", listEnvironments)
		v1.POST("/environments", createEnvironment)
		v1.GET("/environments/:id", getEnvironment)
		v1.PUT("/environments/:id", updateEnvironment)
		v1.DELETE("/environments/:id", deleteEnvironment)

		// Environment Operations
		v1.POST("/environments/:id/provision", provisionEnvironment)
		v1.POST("/environments/:id/start", startEnvironment)
		v1.POST("/environments/:id/stop", stopEnvironment)
		v1.POST("/environments/:id/upload", uploadArtifact)
		v1.GET("/environments/:id/status", getEnvironmentStatus)
		v1.GET("/environments/:id/metrics", getEnvironmentMetrics)
		v1.GET("/environments/:id/logs", getEnvironmentLogs)

		// Validation and Cost
		v1.POST("/validate", validateSpec)
		v1.POST("/cost/estimate", estimateCost)

		// Templates
		v1.GET("/templates", getTemplates)

		// Audit and History
		v1.GET("/audit", getAuditLogs)
		v1.GET("/environments/:id/history", getEnvironmentHistory)

		// AWS FleetWise Operations
		v1.POST("/fleetwise/vehicles", createFleetWiseVehicle)
		v1.POST("/fleetwise/vehicles/batch", batchCreateFleetWiseVehicles)
		v1.GET("/fleetwise/vehicles/:name", getFleetWiseVehicle)
		v1.PUT("/fleetwise/vehicles/:name", updateFleetWiseVehicle)
		v1.DELETE("/fleetwise/vehicles/:name", deleteFleetWiseVehicle)
		v1.GET("/fleetwise/vehicles/:name/status", getFleetWiseVehicleStatus)
		v1.GET("/fleetwise/vehicles", listFleetWiseVehicles)

		v1.POST("/fleetwise/campaigns", createFleetWiseCampaign)
		v1.GET("/fleetwise/campaigns/:name", getFleetWiseCampaign)
		v1.PUT("/fleetwise/campaigns/:name", updateFleetWiseCampaign)
		v1.DELETE("/fleetwise/campaigns/:name", deleteFleetWiseCampaign)
		v1.GET("/fleetwise/campaigns", listFleetWiseCampaigns)

		v1.POST("/fleetwise/fleets", createFleetWiseFleet)
		v1.POST("/fleetwise/fleets/:id/vehicles", associateVehicleToFleet)
	}

	// Start server
	log.Println("Starting SES Platform API on :8080")
	router.Run(":8080")
}

func initDB() {
	dsn := "host=localhost user=ses_user password=ses_password dbname=ses_platform port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schemas
	db.AutoMigrate(
		&Capability{},
		&Enabler{},
		&Environment{},
		&StateTransition{},
		&AuditLog{},
		&Upload{},
		&Reservation{},
	)
}

func seedData() {
	// Seed Capabilities
	capabilities := []Capability{
		{ID: "C01", Name: "Spec Authoring & Validation", Description: "Web editor and services to author SES", Enablers: `["E02","E17","E11","E20"]`, Dependencies: `[]`},
		{ID: "C02", Name: "Parsing & Internal Modeling", Description: "Convert validated SES into normalized models", Enablers: `["E03","E04","E17","E01"]`, Dependencies: `["C01"]`},
		{ID: "C03", Name: "Planning Engine", Description: "Generates execution plan", Enablers: `["E03","E04","E08","E16"]`, Dependencies: `["C01","C02"]`},
		{ID: "C04", Name: "Provisioning Automation", Description: "Executes plan to create infrastructure", Enablers: `["E05","E04","E01","E18","E20"]`, Dependencies: `["C03","C09"]`},
		{ID: "C06", Name: "Monitoring & Metrics", Description: "Collects infrastructure metrics", Enablers: `["E06","E07","E11","E19"]`, Dependencies: `["C04"]`},
		{ID: "C08", Name: "Cost Management", Description: "Real-time cost tracking", Enablers: `["E08","E16","E06","E11"]`, Dependencies: `["C06"]`},
		{ID: "C09", Name: "Security & Compliance", Description: "Credential vault, RBAC", Enablers: `["E09","E10","E19","E01"]`, Dependencies: `[]`},
		{ID: "C12", Name: "Simulation Execution", Description: "Runs scenarios", Enablers: `["E12","E13","E04","E06","E20"]`, Dependencies: `["C04"]`},
		{ID: "C19", Name: "AWS Automotive Integration", Description: "Real AWS IoT FleetWise integration for vehicle simulation", Enablers: `["E05","E21","E01","E18","E20"]`, Dependencies: `["C04","C09","C14"]`},
	}
	for _, cap := range capabilities {
		db.FirstOrCreate(&cap, Capability{ID: cap.ID})
	}

	// Seed Enablers
	enablers := []Enabler{
		{ID: "E01", Name: "Core Platform Infra", Description: "DB, storage, vault"},
		{ID: "E02", Name: "Schema & Validation", Description: "JSON/YAML validators"},
		{ID: "E03", Name: "Graph & Planning", Description: "DAG construction"},
		{ID: "E04", Name: "Execution Framework", Description: "Workflow engine"},
		{ID: "E06", Name: "Metrics Stack", Description: "Prometheus, time-series"},
		{ID: "E08", Name: "Cost Engine", Description: "Pricing cache"},
		{ID: "E09", Name: "Security/RBAC", Description: "Auth, MFA"},
		{ID: "E11", Name: "UI Components", Description: "Web editor, dashboards"},
	}
	for _, enb := range enablers {
		db.FirstOrCreate(&enb, Enabler{ID: enb.ID})
	}
}

// API Handlers
func getCapabilities(c *gin.Context) {
	var capabilities []Capability
	db.Find(&capabilities)
	c.JSON(http.StatusOK, capabilities)
}

func getEnablers(c *gin.Context) {
	var enablers []Enabler
	db.Find(&enablers)
	c.JSON(http.StatusOK, enablers)
}

func listEnvironments(c *gin.Context) {
	var environments []Environment
	db.Order("created_at desc").Find(&environments)
	c.JSON(http.StatusOK, environments)
}

func createEnvironment(c *gin.Context) {
	var req CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate estimated cost
	estimatedCost := calculateCost(req.Compute, req.Storage, len(req.Capabilities))

	// Serialize data
	capabilitiesJSON, _ := json.Marshal(req.Capabilities)
	enablersJSON, _ := json.Marshal(req.EnablersConfig)
	computeJSON, _ := json.Marshal(req.Compute)

	// Serialize FleetWise config if provided
	fleetwiseConfigJSON := ""
	if req.FleetWiseConfig != nil {
		fwJSON, _ := json.Marshal(req.FleetWiseConfig)
		fleetwiseConfigJSON = string(fwJSON)
	}

	env := Environment{
		ID:                fmt.Sprintf("env-%d", time.Now().Unix()),
		Name:              req.Name,
		Description:       req.Description,
		Owner:             req.Owner,
		Tags:              req.Tags,
		Status:            "pending",
		Capabilities:      string(capabilitiesJSON),
		EnablersConfig:    string(enablersJSON),
		ComputeConfig:     string(computeJSON),
		Storage:           req.Storage,
		Network:           req.Network,
		Priority:          req.Priority,
		Duration:          req.Duration,
		EstimatedCost:     estimatedCost,
		ActualCost:        0,
		Health:            100,
		Uptime:            "0h",
		FleetWiseConfig:   fleetwiseConfigJSON,
		UseRealAWSBackend: req.UseRealAWSBackend,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := db.Create(&env).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create audit log
	auditLog := AuditLog{
		EnvironmentID: env.ID,
		Action:        "created",
		UserID:        req.Owner,
		Details:       `{"message":"Environment created"}`,
		CreatedAt:     time.Now(),
	}
	db.Create(&auditLog)

	// Simulate provisioning in background
	go simulateProvisioning(env.ID)

	c.JSON(http.StatusCreated, env)
}

func getEnvironment(c *gin.Context) {
	id := c.Param("id")
	var env Environment
	if err := db.First(&env, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}
	c.JSON(http.StatusOK, env)
}

func updateEnvironment(c *gin.Context) {
	id := c.Param("id")
	var env Environment
	if err := db.First(&env, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()
	db.Model(&env).Updates(updates)

	c.JSON(http.StatusOK, env)
}

func deleteEnvironment(c *gin.Context) {
	id := c.Param("id")
	var env Environment
	if err := db.First(&env, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Create state transition
	transition := StateTransition{
		EnvironmentID: id,
		FromState:     env.Status,
		ToState:       "deleted",
		Reason:        "User requested deletion",
		CreatedAt:     time.Now(),
	}
	db.Create(&transition)

	db.Delete(&env)
	c.JSON(http.StatusOK, gin.H{"message": "Environment deleted"})
}

func provisionEnvironment(c *gin.Context) {
	id := c.Param("id")
	var env Environment
	if err := db.First(&env, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	if env.Status != "pending" && env.Status != "stopped" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Environment cannot be provisioned in current state"})
		return
	}

	// Update status
	db.Model(&env).Updates(map[string]interface{}{
		"status":     "provisioning",
		"updated_at": time.Now(),
	})

	// Create state transition
	transition := StateTransition{
		EnvironmentID: id,
		FromState:     env.Status,
		ToState:       "provisioning",
		Reason:        "Provisioning initiated",
		CreatedAt:     time.Now(),
	}
	db.Create(&transition)

	// Use real AWS or simulate provisioning based on flag
	if env.UseRealAWSBackend && env.FleetWiseConfig != "" {
		go provisionAWSFleetWise(id, env.FleetWiseConfig)
	} else {
		go simulateProvisioning(id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provisioning started"})
}

func startEnvironment(c *gin.Context) {
	id := c.Param("id")
	updateEnvironmentStatus(id, "running", c)
}

func stopEnvironment(c *gin.Context) {
	id := c.Param("id")
	updateEnvironmentStatus(id, "stopped", c)
}

func updateEnvironmentStatus(id, newStatus string, c *gin.Context) {
	var env Environment
	if err := db.First(&env, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	oldStatus := env.Status
	db.Model(&env).Updates(map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	})

	// State transition
	transition := StateTransition{
		EnvironmentID: id,
		FromState:     oldStatus,
		ToState:       newStatus,
		Reason:        fmt.Sprintf("Status changed to %s", newStatus),
		CreatedAt:     time.Now(),
	}
	db.Create(&transition)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Environment %s", newStatus)})
}

func uploadArtifact(c *gin.Context) {
	id := c.Param("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	upload := Upload{
		EnvironmentID: id,
		Filename:      header.Filename,
		FileType:      c.PostForm("file_type"),
		Version:       c.PostForm("version"),
		Size:          header.Size,
		Status:        "completed",
		CreatedAt:     time.Now(),
	}
	db.Create(&upload)

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload successful",
		"upload":  upload,
	})
}

func getEnvironmentStatus(c *gin.Context) {
	id := c.Param("id")
	var env Environment
	if err := db.First(&env, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      env.ID,
		"status":  env.Status,
		"health":  env.Health,
		"uptime":  env.Uptime,
		"cost":    env.ActualCost,
	})
}

func getEnvironmentMetrics(c *gin.Context) {
	id := c.Param("id")

	// Simulate metrics data
	metrics := gin.H{
		"environment_id": id,
		"timestamp":      time.Now(),
		"cpu_usage":      rand.Float64() * 100,
		"memory_usage":   rand.Float64() * 100,
		"disk_usage":     rand.Float64() * 100,
		"network_in":     rand.Float64() * 1000,
		"network_out":    rand.Float64() * 1000,
	}

	c.JSON(http.StatusOK, metrics)
}

func getEnvironmentLogs(c *gin.Context) {
	id := c.Param("id")

	// Simulate log entries
	logs := []gin.H{
		{
			"timestamp": time.Now().Add(-5 * time.Minute),
			"level":     "INFO",
			"message":   "Environment health check passed",
		},
		{
			"timestamp": time.Now().Add(-10 * time.Minute),
			"level":     "INFO",
			"message":   "Metrics collection started",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"environment_id": id,
		"logs":           logs,
	})
}

func validateSpec(c *gin.Context) {
	var req CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, ValidationResponse{
			Valid:  false,
			Errors: []string{err.Error()},
		})
		return
	}

	errors := []string{}
	warnings := []string{}

	// Validate capabilities dependencies
	for _, capID := range req.Capabilities {
		var cap Capability
		if err := db.First(&cap, "id = ?", capID).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Capability %s not found", capID))
			continue
		}

		var deps []string
		json.Unmarshal([]byte(cap.Dependencies), &deps)
		for _, depID := range deps {
			found := false
			for _, selectedCap := range req.Capabilities {
				if selectedCap == depID {
					found = true
					break
				}
			}
			if !found {
				errors = append(errors, fmt.Sprintf("Capability %s requires %s", capID, depID))
			}
		}
	}

	if req.Compute.CPU < 1 {
		errors = append(errors, "CPU must be at least 1")
	}
	if req.Compute.Memory < 1 {
		errors = append(errors, "Memory must be at least 1 GB")
	}
	if req.Storage < 1 {
		warnings = append(warnings, "Storage should be at least 1 GB")
	}

	c.JSON(http.StatusOK, ValidationResponse{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	})
}

func estimateCost(c *gin.Context) {
	var req CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dailyCost := calculateCost(req.Compute, req.Storage, len(req.Capabilities))
	
	breakdown := map[string]float64{
		"compute":      float64(req.Compute.CPU*req.Compute.Memory*req.Compute.Instances) * 0.05,
		"storage":      float64(req.Storage) * 0.1,
		"capabilities": float64(len(req.Capabilities)) * 5.0,
	}

	tip := ""
	if req.Compute.Instances > 5 {
		tip = "Consider using auto-scaling to optimize costs during low usage"
	}

	c.JSON(http.StatusOK, CostEstimationResponse{
		DailyCost:       dailyCost,
		MonthlyCost:     dailyCost * 30,
		Breakdown:       breakdown,
		OptimizationTip: tip,
	})
}

func getTemplates(c *gin.Context) {
	templates := []gin.H{
		{
			"id":           "tmpl-001",
			"name":         "Standard Load Test",
			"description":  "Pre-configured for load testing",
			"capabilities": []string{"C01", "C02", "C03", "C04", "C06", "C12"},
			"popularity":   145,
			"cost":         98.50,
		},
		{
			"id":           "tmpl-002",
			"name":         "Full Production Mirror",
			"description":  "Complete production environment",
			"capabilities": []string{"C01", "C02", "C03", "C04", "C06", "C08", "C09", "C12"},
			"popularity":   89,
			"cost":         287.20,
		},
	}

	c.JSON(http.StatusOK, templates)
}

func getAuditLogs(c *gin.Context) {
	envID := c.Query("environment_id")
	var logs []AuditLog

	query := db.Order("created_at desc").Limit(100)
	if envID != "" {
		query = query.Where("environment_id = ?", envID)
	}
	query.Find(&logs)

	c.JSON(http.StatusOK, logs)
}

func getEnvironmentHistory(c *gin.Context) {
	id := c.Param("id")
	var transitions []StateTransition
	db.Where("environment_id = ?", id).Order("created_at desc").Find(&transitions)
	c.JSON(http.StatusOK, transitions)
}

// Helper Functions
func calculateCost(compute ComputeConfig, storage, capabilityCount int) float64 {
	computeCost := float64(compute.CPU*compute.Memory*compute.Instances) * 0.05
	storageCost := float64(storage) * 0.1
	capabilityCost := float64(capabilityCount) * 5.0
	return computeCost + storageCost + capabilityCost
}

func simulateProvisioning(envID string) {
	// Simulate provisioning stages
	stages := []string{"validating", "allocating", "configuring", "starting"}
	
	for i, stage := range stages {
		time.Sleep(2 * time.Second)
		
		// Update environment status
		var env Environment
		db.First(&env, "id = ?", envID)
		
		progress := int((float64(i+1) / float64(len(stages))) * 100)
		status := "provisioning"
		if progress == 100 {
			status = "running"
		}
		
		db.Model(&env).Updates(map[string]interface{}{
			"status":     status,
			"health":     90 + rand.Intn(10),
			"updated_at": time.Now(),
		})
		
		// Create state transition
		transition := StateTransition{
			EnvironmentID: envID,
			FromState:     env.Status,
			ToState:       status,
			Reason:        fmt.Sprintf("Stage: %s (%d%%)", stage, progress),
			Metadata:      fmt.Sprintf(`{"stage":"%s","progress":%d}`, stage, progress),
			CreatedAt:     time.Now(),
		}
		db.Create(&transition)
	}
	
	// Update uptime
	go updateUptime(envID)
}

func updateUptime(envID string) {
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		var env Environment
		if err := db.First(&env, "id = ?", envID).Error; err != nil {
			return
		}

		if env.Status != "running" {
			return
		}

		duration := time.Since(startTime)
		uptime := fmt.Sprintf("%dd %dh", int(duration.Hours()/24), int(duration.Hours())%24)

		db.Model(&env).Updates(map[string]interface{}{
			"uptime":      uptime,
			"actual_cost": env.EstimatedCost * (duration.Hours() / 24),
		})
	}
}

// provisionAWSFleetWise provisions an environment using real AWS IoT FleetWise
func provisionAWSFleetWise(envID string, configJSON string) {
	log.Printf("Starting real AWS FleetWise provisioning for environment: %s", envID)

	// Parse FleetWise configuration
	var config FleetWiseConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		log.Printf("Error parsing FleetWise config: %v", err)
		updateEnvironmentStatusWithError(envID, "error", fmt.Sprintf("Failed to parse FleetWise config: %v", err))
		return
	}

	// Create FleetWise client
	client, err := NewAWSFleetWiseClient(config.Region)
	if err != nil {
		log.Printf("Error creating FleetWise client: %v", err)
		updateEnvironmentStatusWithError(envID, "error", fmt.Sprintf("Failed to create AWS client: %v", err))
		return
	}

	// Update status to validating
	updateEnvironmentStatusWithTransition(envID, "provisioning", "Validating FleetWise configuration")
	time.Sleep(1 * time.Second)

	// Provision the environment
	err = client.ProvisionFleetWiseEnvironment(envID, config)
	if err != nil {
		log.Printf("Error provisioning FleetWise environment: %v", err)
		updateEnvironmentStatusWithError(envID, "error", fmt.Sprintf("Provisioning failed: %v", err))
		return
	}

	// Update status to running
	updateEnvironmentStatusWithTransition(envID, "running", "AWS FleetWise environment provisioned successfully")

	// Update health to 95-100 (real AWS environment)
	var env Environment
	db.First(&env, "id = ?", envID)
	db.Model(&env).Updates(map[string]interface{}{
		"health":     95 + rand.Intn(5),
		"updated_at": time.Now(),
	})

	log.Printf("Successfully provisioned AWS FleetWise environment: %s", envID)

	// Start uptime tracking
	go updateUptime(envID)
}

// Helper function to update environment status with transition
func updateEnvironmentStatusWithTransition(envID, status, reason string) {
	var env Environment
	db.First(&env, "id = ?", envID)

	oldStatus := env.Status
	db.Model(&env).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	})

	transition := StateTransition{
		EnvironmentID: envID,
		FromState:     oldStatus,
		ToState:       status,
		Reason:        reason,
		CreatedAt:     time.Now(),
	}
	db.Create(&transition)
}

// Helper function to update environment with error status
func updateEnvironmentStatusWithError(envID, status, errorMsg string) {
	var env Environment
	db.First(&env, "id = ?", envID)

	oldStatus := env.Status
	db.Model(&env).Updates(map[string]interface{}{
		"status":     status,
		"health":     0,
		"updated_at": time.Now(),
	})

	transition := StateTransition{
		EnvironmentID: envID,
		FromState:     oldStatus,
		ToState:       status,
		Reason:        errorMsg,
		Metadata:      fmt.Sprintf(`{"error": "%s"}`, errorMsg),
		CreatedAt:     time.Now(),
	}
	db.Create(&transition)

	// Create audit log for error
	auditLog := AuditLog{
		EnvironmentID: envID,
		Action:        "provisioning_failed",
		UserID:        "system",
		Details:       fmt.Sprintf(`{"error": "%s"}`, errorMsg),
		CreatedAt:     time.Now(),
	}
	db.Create(&auditLog)
}
