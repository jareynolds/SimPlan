package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotfleetwise"
	"github.com/aws/aws-sdk-go-v2/service/iotfleetwise/types"
)

// AWSFleetWiseClient wraps AWS IoT FleetWise operations
type AWSFleetWiseClient struct {
	client *iotfleetwise.Client
	ctx    context.Context
}

// FleetWiseConfig holds configuration for AWS IoT FleetWise integration
type FleetWiseConfig struct {
	Region              string   `json:"region"`
	SignalCatalogARN    string   `json:"signal_catalog_arn"`
	ModelManifestARN    string   `json:"model_manifest_arn"`
	DecoderManifestARN  string   `json:"decoder_manifest_arn"`
	FleetID             string   `json:"fleet_id"`
	CampaignARN         string   `json:"campaign_arn"`
	VehicleNames        []string `json:"vehicle_names"`
	DataDestinationS3   string   `json:"data_destination_s3"`
	DataDestinationMQTT string   `json:"data_destination_mqtt"`
	EnableCompression   bool     `json:"enable_compression"`
	EnableSpooling      bool     `json:"enable_spooling"`
	EnableDiagnostics   bool     `json:"enable_diagnostics"`
}

// VehicleConfig represents vehicle-specific configuration
type VehicleConfig struct {
	Name               string            `json:"name"`
	ModelManifestARN   string            `json:"model_manifest_arn"`
	DecoderManifestARN string            `json:"decoder_manifest_arn"`
	Attributes         map[string]string `json:"attributes"`
	CreateIoTThing     bool              `json:"create_iot_thing"`
}

// CampaignConfig represents campaign configuration
type CampaignConfig struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	SignalCatalogARN      string                 `json:"signal_catalog_arn"`
	TargetARN             string                 `json:"target_arn"`
	CollectionScheme      CollectionScheme       `json:"collection_scheme"`
	SignalsToCollect      []SignalToCollect      `json:"signals_to_collect"`
	DataDestinations      []DataDestination      `json:"data_destinations"`
	Compression           string                 `json:"compression"`
	DiagnosticsMode       string                 `json:"diagnostics_mode"`
	SpoolingMode          string                 `json:"spooling_mode"`
	DataExtraDimensions   []string               `json:"data_extra_dimensions"`
	PostTriggerDurationMs int64                  `json:"post_trigger_duration_ms"`
}

// CollectionScheme defines how data is collected
type CollectionScheme struct {
	Type                     string `json:"type"` // "time-based" or "condition-based"
	PeriodMs                 int64  `json:"period_ms,omitempty"`
	Expression               string `json:"expression,omitempty"`
	MinimumTriggerIntervalMs int64  `json:"minimum_trigger_interval_ms,omitempty"`
	TriggerMode              string `json:"trigger_mode,omitempty"` // "ALWAYS" or "RISING_EDGE"
}

// SignalToCollect defines which signals to collect
type SignalToCollect struct {
	Name                      string `json:"name"`
	MaxSampleCount            int64  `json:"max_sample_count"`
	MinimumSamplingIntervalMs int64  `json:"minimum_sampling_interval_ms"`
}

// DataDestination defines where data should be sent
type DataDestination struct {
	Type                     string `json:"type"` // "s3", "timestream", or "mqtt"
	S3BucketARN              string `json:"s3_bucket_arn,omitempty"`
	S3Prefix                 string `json:"s3_prefix,omitempty"`
	S3DataFormat             string `json:"s3_data_format,omitempty"` // "JSON" or "PARQUET"
	S3StorageCompression     string `json:"s3_storage_compression,omitempty"` // "NONE" or "GZIP"
	TimestreamTableARN       string `json:"timestream_table_arn,omitempty"`
	TimestreamExecutionRole  string `json:"timestream_execution_role,omitempty"`
	MQTTTopicARN             string `json:"mqtt_topic_arn,omitempty"`
	MQTTExecutionRole        string `json:"mqtt_execution_role,omitempty"`
}

// NewAWSFleetWiseClient creates a new FleetWise client
func NewAWSFleetWiseClient(region string) (*AWSFleetWiseClient, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := iotfleetwise.NewFromConfig(cfg)

	return &AWSFleetWiseClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// CreateVehicle creates a new vehicle in AWS IoT FleetWise
func (c *AWSFleetWiseClient) CreateVehicle(vehicleConfig VehicleConfig) (*types.CreateVehicleOutput, error) {
	log.Printf("Creating vehicle: %s", vehicleConfig.Name)

	// Convert attributes map to AWS SDK format
	attributes := make(map[string]string)
	for k, v := range vehicleConfig.Attributes {
		attributes[k] = v
	}

	// Determine association behavior
	var associationBehavior types.VehicleAssociationBehavior
	if vehicleConfig.CreateIoTThing {
		associationBehavior = types.VehicleAssociationBehaviorCreateIotThing
	} else {
		associationBehavior = types.VehicleAssociationBehaviorValidateIotThingExists
	}

	input := &iotfleetwise.CreateVehicleInput{
		VehicleName:          aws.String(vehicleConfig.Name),
		ModelManifestArn:     aws.String(vehicleConfig.ModelManifestARN),
		DecoderManifestArn:   aws.String(vehicleConfig.DecoderManifestARN),
		Attributes:           attributes,
		AssociationBehavior:  associationBehavior,
	}

	result, err := c.client.CreateVehicle(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %v", err)
	}

	log.Printf("Successfully created vehicle: %s (ARN: %s)", vehicleConfig.Name, *result.Arn)
	return result, nil
}

// BatchCreateVehicles creates multiple vehicles in a batch
func (c *AWSFleetWiseClient) BatchCreateVehicles(vehicles []VehicleConfig) ([]string, []error) {
	var createdARNs []string
	var errors []error

	// AWS supports up to 10 vehicles per batch, so we'll batch them
	batchSize := 10
	for i := 0; i < len(vehicles); i += batchSize {
		end := i + batchSize
		if end > len(vehicles) {
			end = len(vehicles)
		}

		batch := vehicles[i:end]
		var batchInput []types.CreateVehicleRequestItem

		for _, v := range batch {
			attributes := make(map[string]string)
			for k, val := range v.Attributes {
				attributes[k] = val
			}

			var associationBehavior types.VehicleAssociationBehavior
			if v.CreateIoTThing {
				associationBehavior = types.VehicleAssociationBehaviorCreateIotThing
			} else {
				associationBehavior = types.VehicleAssociationBehaviorValidateIotThingExists
			}

			batchInput = append(batchInput, types.CreateVehicleRequestItem{
				VehicleName:          aws.String(v.Name),
				ModelManifestArn:     aws.String(v.ModelManifestARN),
				DecoderManifestArn:   aws.String(v.DecoderManifestARN),
				Attributes:           attributes,
				AssociationBehavior:  associationBehavior,
			})
		}

		input := &iotfleetwise.BatchCreateVehicleInput{
			Vehicles: batchInput,
		}

		result, err := c.client.BatchCreateVehicle(c.ctx, input)
		if err != nil {
			errors = append(errors, fmt.Errorf("batch create failed: %v", err))
			continue
		}

		// Process successful vehicles
		for _, v := range result.Vehicles {
			if v.Arn != nil {
				createdARNs = append(createdARNs, *v.Arn)
				log.Printf("Created vehicle: %s (ARN: %s)", *v.VehicleName, *v.Arn)
			}
		}

		// Process errors
		for _, e := range result.Errors {
			errors = append(errors, fmt.Errorf("vehicle %s failed: %s", *e.VehicleName, *e.Message))
		}
	}

	return createdARNs, errors
}

// GetVehicle retrieves vehicle information
func (c *AWSFleetWiseClient) GetVehicle(vehicleName string) (*iotfleetwise.GetVehicleOutput, error) {
	input := &iotfleetwise.GetVehicleInput{
		VehicleName: aws.String(vehicleName),
	}

	result, err := c.client.GetVehicle(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %v", err)
	}

	return result, nil
}

// UpdateVehicle updates vehicle configuration
func (c *AWSFleetWiseClient) UpdateVehicle(vehicleName string, updates VehicleConfig) error {
	log.Printf("Updating vehicle: %s", vehicleName)

	attributes := make(map[string]string)
	for k, v := range updates.Attributes {
		attributes[k] = v
	}

	input := &iotfleetwise.UpdateVehicleInput{
		VehicleName:        aws.String(vehicleName),
		ModelManifestArn:   aws.String(updates.ModelManifestARN),
		DecoderManifestArn: aws.String(updates.DecoderManifestARN),
		Attributes:         attributes,
		AttributeUpdateMode: types.UpdateModeOverwrite,
	}

	_, err := c.client.UpdateVehicle(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %v", err)
	}

	log.Printf("Successfully updated vehicle: %s", vehicleName)
	return nil
}

// DeleteVehicle deletes a vehicle
func (c *AWSFleetWiseClient) DeleteVehicle(vehicleName string) error {
	log.Printf("Deleting vehicle: %s", vehicleName)

	input := &iotfleetwise.DeleteVehicleInput{
		VehicleName: aws.String(vehicleName),
	}

	_, err := c.client.DeleteVehicle(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %v", err)
	}

	log.Printf("Successfully deleted vehicle: %s", vehicleName)
	return nil
}

// CreateCampaign creates a data collection campaign
func (c *AWSFleetWiseClient) CreateCampaign(campaignConfig CampaignConfig) (*types.CreateCampaignOutput, error) {
	log.Printf("Creating campaign: %s", campaignConfig.Name)

	// Build collection scheme
	var collectionScheme types.CollectionScheme
	if campaignConfig.CollectionScheme.Type == "time-based" {
		collectionScheme = &types.CollectionSchemeMemberTimeBasedCollectionScheme{
			Value: types.TimeBasedCollectionScheme{
				PeriodMs: aws.Int64(campaignConfig.CollectionScheme.PeriodMs),
			},
		}
	} else if campaignConfig.CollectionScheme.Type == "condition-based" {
		var triggerMode types.TriggerMode
		if campaignConfig.CollectionScheme.TriggerMode == "RISING_EDGE" {
			triggerMode = types.TriggerModeRisingEdge
		} else {
			triggerMode = types.TriggerModeAlways
		}

		collectionScheme = &types.CollectionSchemeMemberConditionBasedCollectionScheme{
			Value: types.ConditionBasedCollectionScheme{
				Expression:                 aws.String(campaignConfig.CollectionScheme.Expression),
				MinimumTriggerIntervalMs:   aws.Int64(campaignConfig.CollectionScheme.MinimumTriggerIntervalMs),
				TriggerMode:                triggerMode,
				ConditionLanguageVersion:   aws.Int32(1),
			},
		}
	}

	// Build signals to collect
	var signalsToCollect []types.SignalInformation
	for _, signal := range campaignConfig.SignalsToCollect {
		signalsToCollect = append(signalsToCollect, types.SignalInformation{
			Name:                       aws.String(signal.Name),
			MaxSampleCount:             aws.Int64(signal.MaxSampleCount),
			MinimumSamplingIntervalMs:  aws.Int64(signal.MinimumSamplingIntervalMs),
		})
	}

	// Build data destination configs
	var dataDestinationConfigs []types.DataDestinationConfig
	for _, dest := range campaignConfig.DataDestinations {
		switch dest.Type {
		case "s3":
			var dataFormat types.DataFormat
			if dest.S3DataFormat == "PARQUET" {
				dataFormat = types.DataFormatParquet
			} else {
				dataFormat = types.DataFormatJson
			}

			var compression types.StorageCompressionFormat
			if dest.S3StorageCompression == "GZIP" {
				compression = types.StorageCompressionFormatGzip
			} else {
				compression = types.StorageCompressionFormatNone
			}

			dataDestinationConfigs = append(dataDestinationConfigs, &types.DataDestinationConfigMemberS3Config{
				Value: types.S3Config{
					BucketArn:                 aws.String(dest.S3BucketARN),
					Prefix:                    aws.String(dest.S3Prefix),
					DataFormat:                dataFormat,
					StorageCompressionFormat:  compression,
				},
			})

		case "timestream":
			dataDestinationConfigs = append(dataDestinationConfigs, &types.DataDestinationConfigMemberTimestreamConfig{
				Value: types.TimestreamConfig{
					TimestreamTableArn: aws.String(dest.TimestreamTableARN),
					ExecutionRoleArn:   aws.String(dest.TimestreamExecutionRole),
				},
			})

		case "mqtt":
			dataDestinationConfigs = append(dataDestinationConfigs, &types.DataDestinationConfigMemberMqttTopicConfig{
				Value: types.MqttTopicConfig{
					MqttTopicArn:     aws.String(dest.MQTTTopicARN),
					ExecutionRoleArn: aws.String(dest.MQTTExecutionRole),
				},
			})
		}
	}

	// Set compression
	var compression types.Compression
	if campaignConfig.Compression == "OFF" {
		compression = types.CompressionOff
	} else {
		compression = types.CompressionSnappy
	}

	// Set diagnostics mode
	var diagnosticsMode types.DiagnosticsMode
	if campaignConfig.DiagnosticsMode == "SEND_ACTIVE_DTCS" {
		diagnosticsMode = types.DiagnosticsModeSendActiveDtcs
	} else {
		diagnosticsMode = types.DiagnosticsModeOff
	}

	// Set spooling mode
	var spoolingMode types.SpoolingMode
	if campaignConfig.SpoolingMode == "TO_DISK" {
		spoolingMode = types.SpoolingModeToDisk
	} else {
		spoolingMode = types.SpoolingModeOff
	}

	input := &iotfleetwise.CreateCampaignInput{
		Name:                           aws.String(campaignConfig.Name),
		Description:                    aws.String(campaignConfig.Description),
		SignalCatalogArn:               aws.String(campaignConfig.SignalCatalogARN),
		TargetArn:                      aws.String(campaignConfig.TargetARN),
		CollectionScheme:               collectionScheme,
		SignalsToCollect:               signalsToCollect,
		DataDestinationConfigs:         dataDestinationConfigs,
		Compression:                    compression,
		DiagnosticsMode:                diagnosticsMode,
		SpoolingMode:                   spoolingMode,
		DataExtraDimensions:            campaignConfig.DataExtraDimensions,
		PostTriggerCollectionDuration:  aws.Int64(campaignConfig.PostTriggerDurationMs),
	}

	result, err := c.client.CreateCampaign(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %v", err)
	}

	log.Printf("Successfully created campaign: %s (ARN: %s)", campaignConfig.Name, *result.Arn)
	return result, nil
}

// GetCampaign retrieves campaign information
func (c *AWSFleetWiseClient) GetCampaign(campaignName string) (*iotfleetwise.GetCampaignOutput, error) {
	input := &iotfleetwise.GetCampaignInput{
		Name: aws.String(campaignName),
	}

	result, err := c.client.GetCampaign(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %v", err)
	}

	return result, nil
}

// UpdateCampaign updates campaign configuration
func (c *AWSFleetWiseClient) UpdateCampaign(campaignName string, action string) error {
	log.Printf("Updating campaign %s with action: %s", campaignName, action)

	var updateAction types.UpdateCampaignAction
	switch action {
	case "APPROVE":
		updateAction = types.UpdateCampaignActionApprove
	case "SUSPEND":
		updateAction = types.UpdateCampaignActionSuspend
	case "RESUME":
		updateAction = types.UpdateCampaignActionResume
	case "UPDATE":
		updateAction = types.UpdateCampaignActionUpdate
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	input := &iotfleetwise.UpdateCampaignInput{
		Name:   aws.String(campaignName),
		Action: updateAction,
	}

	_, err := c.client.UpdateCampaign(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %v", err)
	}

	log.Printf("Successfully updated campaign: %s", campaignName)
	return nil
}

// DeleteCampaign deletes a campaign
func (c *AWSFleetWiseClient) DeleteCampaign(campaignName string) error {
	log.Printf("Deleting campaign: %s", campaignName)

	input := &iotfleetwise.DeleteCampaignInput{
		Name: aws.String(campaignName),
	}

	_, err := c.client.DeleteCampaign(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %v", err)
	}

	log.Printf("Successfully deleted campaign: %s", campaignName)
	return nil
}

// CreateFleet creates a vehicle fleet
func (c *AWSFleetWiseClient) CreateFleet(fleetID, description, signalCatalogARN string) (*types.CreateFleetOutput, error) {
	log.Printf("Creating fleet: %s", fleetID)

	input := &iotfleetwise.CreateFleetInput{
		FleetId:          aws.String(fleetID),
		Description:      aws.String(description),
		SignalCatalogArn: aws.String(signalCatalogARN),
	}

	result, err := c.client.CreateFleet(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create fleet: %v", err)
	}

	log.Printf("Successfully created fleet: %s (ARN: %s)", fleetID, *result.Arn)
	return result, nil
}

// AssociateVehicleToFleet associates a vehicle with a fleet
func (c *AWSFleetWiseClient) AssociateVehicleToFleet(vehicleName, fleetID string) error {
	log.Printf("Associating vehicle %s to fleet %s", vehicleName, fleetID)

	input := &iotfleetwise.AssociateVehicleFleetInput{
		VehicleName: aws.String(vehicleName),
		FleetId:     aws.String(fleetID),
	}

	_, err := c.client.AssociateVehicleFleet(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to associate vehicle to fleet: %v", err)
	}

	log.Printf("Successfully associated vehicle %s to fleet %s", vehicleName, fleetID)
	return nil
}

// ListVehicles lists all vehicles
func (c *AWSFleetWiseClient) ListVehicles(modelManifestARN string, maxResults int32) ([]types.VehicleSummary, error) {
	input := &iotfleetwise.ListVehiclesInput{
		MaxResults: aws.Int32(maxResults),
	}

	if modelManifestARN != "" {
		input.ModelManifestArn = aws.String(modelManifestARN)
	}

	result, err := c.client.ListVehicles(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles: %v", err)
	}

	return result.VehicleSummaries, nil
}

// ListCampaigns lists all campaigns
func (c *AWSFleetWiseClient) ListCampaigns(maxResults int32) ([]types.CampaignSummary, error) {
	input := &iotfleetwise.ListCampaignsInput{
		MaxResults: aws.Int32(maxResults),
	}

	result, err := c.client.ListCampaigns(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %v", err)
	}

	return result.CampaignSummaries, nil
}

// GetVehicleStatus gets the status of a vehicle
func (c *AWSFleetWiseClient) GetVehicleStatus(vehicleName string) (*iotfleetwise.GetVehicleStatusOutput, error) {
	input := &iotfleetwise.GetVehicleStatusInput{
		VehicleName: aws.String(vehicleName),
	}

	result, err := c.client.GetVehicleStatus(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle status: %v", err)
	}

	return result, nil
}

// ProvisionFleetWiseEnvironment provisions a complete FleetWise environment
func (c *AWSFleetWiseClient) ProvisionFleetWiseEnvironment(envID string, config FleetWiseConfig) error {
	log.Printf("Provisioning FleetWise environment: %s", envID)

	// Step 1: Create fleet if specified
	if config.FleetID != "" {
		_, err := c.CreateFleet(config.FleetID, fmt.Sprintf("Fleet for environment %s", envID), config.SignalCatalogARN)
		if err != nil {
			return fmt.Errorf("failed to create fleet: %v", err)
		}
	}

	// Step 2: Create vehicles
	var vehicles []VehicleConfig
	for _, name := range config.VehicleNames {
		vehicles = append(vehicles, VehicleConfig{
			Name:               name,
			ModelManifestARN:   config.ModelManifestARN,
			DecoderManifestARN: config.DecoderManifestARN,
			Attributes: map[string]string{
				"EnvironmentID": envID,
				"CreatedAt":     time.Now().Format(time.RFC3339),
			},
			CreateIoTThing: true,
		})
	}

	createdARNs, errors := c.BatchCreateVehicles(vehicles)
	if len(errors) > 0 {
		log.Printf("Errors creating vehicles: %v", errors)
	}
	log.Printf("Created %d vehicles", len(createdARNs))

	// Step 3: Associate vehicles to fleet
	if config.FleetID != "" {
		for _, name := range config.VehicleNames {
			err := c.AssociateVehicleToFleet(name, config.FleetID)
			if err != nil {
				log.Printf("Warning: failed to associate vehicle %s to fleet: %v", name, err)
			}
		}
	}

	// Step 4: Create campaign if configured
	if config.CampaignARN != "" {
		// Extract campaign name from ARN or use a generated name
		campaignName := fmt.Sprintf("campaign-%s", envID)

		// Build data destinations
		var destinations []DataDestination
		if config.DataDestinationS3 != "" {
			destinations = append(destinations, DataDestination{
				Type:                 "s3",
				S3BucketARN:          config.DataDestinationS3,
				S3Prefix:             fmt.Sprintf("fleetwise/%s/", envID),
				S3DataFormat:         "JSON",
				S3StorageCompression: "GZIP",
			})
		}
		if config.DataDestinationMQTT != "" {
			destinations = append(destinations, DataDestination{
				Type:              "mqtt",
				MQTTTopicARN:      config.DataDestinationMQTT,
				MQTTExecutionRole: "", // Should be provided in config
			})
		}

		compression := "OFF"
		if config.EnableCompression {
			compression = "SNAPPY"
		}

		diagnosticsMode := "OFF"
		if config.EnableDiagnostics {
			diagnosticsMode = "SEND_ACTIVE_DTCS"
		}

		spoolingMode := "OFF"
		if config.EnableSpooling {
			spoolingMode = "TO_DISK"
		}

		targetARN := config.FleetID
		if targetARN == "" && len(createdARNs) > 0 {
			targetARN = createdARNs[0] // Use first vehicle if no fleet
		}

		campaignConfig := CampaignConfig{
			Name:             campaignName,
			Description:      fmt.Sprintf("Data collection campaign for environment %s", envID),
			SignalCatalogARN: config.SignalCatalogARN,
			TargetARN:        targetARN,
			CollectionScheme: CollectionScheme{
				Type:     "time-based",
				PeriodMs: 10000, // 10 seconds
			},
			SignalsToCollect: []SignalToCollect{
				{
					Name:                      "Vehicle.Speed",
					MaxSampleCount:            1000,
					MinimumSamplingIntervalMs: 100,
				},
			},
			DataDestinations:      destinations,
			Compression:           compression,
			DiagnosticsMode:       diagnosticsMode,
			SpoolingMode:          spoolingMode,
			PostTriggerDurationMs: 0,
		}

		_, err := c.CreateCampaign(campaignConfig)
		if err != nil {
			log.Printf("Warning: failed to create campaign: %v", err)
		}
	}

	log.Printf("Successfully provisioned FleetWise environment: %s", envID)
	return nil
}

// DeProvisionFleetWiseEnvironment cleans up FleetWise resources
func (c *AWSFleetWiseClient) DeProvisionFleetWiseEnvironment(envID string, config FleetWiseConfig) error {
	log.Printf("De-provisioning FleetWise environment: %s", envID)

	// Delete campaign
	if config.CampaignARN != "" {
		campaignName := fmt.Sprintf("campaign-%s", envID)
		err := c.DeleteCampaign(campaignName)
		if err != nil {
			log.Printf("Warning: failed to delete campaign: %v", err)
		}
	}

	// Delete vehicles
	for _, name := range config.VehicleNames {
		err := c.DeleteVehicle(name)
		if err != nil {
			log.Printf("Warning: failed to delete vehicle %s: %v", name, err)
		}
	}

	// Note: Fleets, signal catalogs, model manifests, and decoder manifests
	// are typically not deleted as they may be reused across environments

	log.Printf("Successfully de-provisioned FleetWise environment: %s", envID)
	return nil
}

// MarshalFleetWiseConfig converts FleetWiseConfig to JSON string
func MarshalFleetWiseConfig(config FleetWiseConfig) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalFleetWiseConfig converts JSON string to FleetWiseConfig
func UnmarshalFleetWiseConfig(data string) (FleetWiseConfig, error) {
	var config FleetWiseConfig
	err := json.Unmarshal([]byte(data), &config)
	return config, err
}
