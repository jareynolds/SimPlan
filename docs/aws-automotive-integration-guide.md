# AWS Automotive Integration Guide

This guide explains how to use the AWS IoT FleetWise integration capability (C19) in SimPlan to configure and manage real automotive vehicle simulators using AWS services.

---

## Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Configuration](#configuration)
5. [API Usage](#api-usage)
6. [Examples](#examples)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

---

## Overview

The AWS Automotive Integration capability (C19) replaces simulated backend functions with real AWS IoT FleetWise APIs. This enables you to:

- **Create and manage real vehicles** in AWS IoT FleetWise
- **Configure data collection campaigns** for vehicle telemetry
- **Organize vehicles into fleets** for efficient management
- **Collect and store vehicle data** in S3, Timestream, or MQTT
- **Monitor vehicle health and status** in real-time
- **Scale to thousands of vehicles** with batch operations

---

## Prerequisites

### AWS Account Setup

1. **AWS Account** with IoT FleetWise enabled
2. **IAM Credentials** with the following permissions:
   - `iotfleetwise:*`
   - `iot:*`
   - `s3:*` (for data storage)
   - `timestream:*` (optional, for time-series data)
   - `iam:PassRole`

3. **AWS Resources** pre-configured:
   - Signal Catalog ARN
   - Model Manifest ARN
   - Decoder Manifest ARN
   - (Optional) S3 bucket for data storage
   - (Optional) Timestream database and table

### Environment Setup

```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="us-east-1"
```

### Backend Dependencies

The Go backend requires AWS SDK v2:

```bash
cd ses-platform/backend
go mod download
```

---

## Quick Start

### 1. Create an Environment with AWS FleetWise

```bash
curl -X POST http://localhost:8080/api/v1/environments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production Vehicle Fleet",
    "description": "Fleet of connected vehicles for data collection",
    "owner": "fleet-team",
    "capabilities": ["C19"],
    "use_real_aws_backend": true,
    "fleetwise_config": {
      "region": "us-east-1",
      "signal_catalog_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:signal-catalog/my-catalog",
      "model_manifest_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:model-manifest/my-model",
      "decoder_manifest_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:decoder-manifest/my-decoder",
      "fleet_id": "production-fleet",
      "vehicle_names": ["vehicle-001", "vehicle-002", "vehicle-003"],
      "data_destination_s3": "arn:aws:s3:::my-fleetwise-data",
      "enable_compression": true,
      "enable_spooling": true,
      "enable_diagnostics": true
    },
    "compute": {
      "cpu": 4,
      "memory": 16,
      "instances": 3
    },
    "storage": 100,
    "priority": "high",
    "duration": 168
  }'
```

### 2. Provision the Environment

```bash
curl -X POST http://localhost:8080/api/v1/environments/{env-id}/provision
```

This will:
1. Create a fleet in AWS IoT FleetWise
2. Create all specified vehicles
3. Associate vehicles with the fleet
4. Set up a data collection campaign
5. Start collecting vehicle data

---

## Configuration

### FleetWise Configuration Object

```json
{
  "region": "us-east-1",
  "signal_catalog_arn": "arn:aws:iotfleetwise:...:signal-catalog/name",
  "model_manifest_arn": "arn:aws:iotfleetwise:...:model-manifest/name",
  "decoder_manifest_arn": "arn:aws:iotfleetwise:...:decoder-manifest/name",
  "fleet_id": "my-fleet",
  "campaign_arn": "arn:aws:iotfleetwise:...:campaign/name",
  "vehicle_names": ["vehicle-001", "vehicle-002"],
  "data_destination_s3": "arn:aws:s3:::bucket-name",
  "data_destination_mqtt": "arn:aws:iot:...:topic/name",
  "enable_compression": true,
  "enable_spooling": true,
  "enable_diagnostics": true
}
```

### Configuration Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `region` | string | Yes | AWS region (e.g., "us-east-1") |
| `signal_catalog_arn` | string | Yes | ARN of signal catalog |
| `model_manifest_arn` | string | Yes | ARN of model manifest |
| `decoder_manifest_arn` | string | Yes | ARN of decoder manifest |
| `fleet_id` | string | No | Fleet ID to create/use |
| `campaign_arn` | string | No | Existing campaign ARN |
| `vehicle_names` | array | Yes | List of vehicle names to create |
| `data_destination_s3` | string | No | S3 bucket ARN for data storage |
| `data_destination_mqtt` | string | No | MQTT topic ARN for streaming |
| `enable_compression` | boolean | No | Enable SNAPPY compression (default: false) |
| `enable_spooling` | boolean | No | Enable offline data spooling (default: false) |
| `enable_diagnostics` | boolean | No | Enable DTC collection (default: false) |

---

## API Usage

### Vehicle Management

#### Create a Single Vehicle

```bash
POST /api/v1/fleetwise/vehicles?region=us-east-1
Content-Type: application/json

{
  "name": "vehicle-001",
  "model_manifest_arn": "arn:aws:iotfleetwise:...:model-manifest/my-model",
  "decoder_manifest_arn": "arn:aws:iotfleetwise:...:decoder-manifest/my-decoder",
  "attributes": {
    "Make": "Tesla",
    "Model": "Model 3",
    "Year": "2024",
    "VIN": "5YJ3E1EA8KF123456"
  },
  "create_iot_thing": true
}
```

#### Create Multiple Vehicles (Batch)

```bash
POST /api/v1/fleetwise/vehicles/batch
Content-Type: application/json

{
  "region": "us-east-1",
  "vehicles": [
    {
      "name": "vehicle-001",
      "model_manifest_arn": "arn:aws:iotfleetwise:...",
      "decoder_manifest_arn": "arn:aws:iotfleetwise:...",
      "attributes": {...},
      "create_iot_thing": true
    },
    {
      "name": "vehicle-002",
      "model_manifest_arn": "arn:aws:iotfleetwise:...",
      "decoder_manifest_arn": "arn:aws:iotfleetwise:...",
      "attributes": {...},
      "create_iot_thing": true
    }
  ]
}
```

#### Get Vehicle Information

```bash
GET /api/v1/fleetwise/vehicles/vehicle-001?region=us-east-1
```

#### Update Vehicle

```bash
PUT /api/v1/fleetwise/vehicles/vehicle-001?region=us-east-1
Content-Type: application/json

{
  "attributes": {
    "Odometer": "25000",
    "FirmwareVersion": "2024.3.5"
  }
}
```

#### Delete Vehicle

```bash
DELETE /api/v1/fleetwise/vehicles/vehicle-001?region=us-east-1
```

#### Get Vehicle Status

```bash
GET /api/v1/fleetwise/vehicles/vehicle-001/status?region=us-east-1
```

#### List All Vehicles

```bash
GET /api/v1/fleetwise/vehicles?region=us-east-1&model_manifest_arn=arn:...
```

### Campaign Management

#### Create a Campaign

```bash
POST /api/v1/fleetwise/campaigns?region=us-east-1
Content-Type: application/json

{
  "name": "production-data-collection",
  "description": "Collect vehicle telemetry data",
  "signal_catalog_arn": "arn:aws:iotfleetwise:...:signal-catalog/my-catalog",
  "target_arn": "arn:aws:iotfleetwise:...:fleet/production-fleet",
  "collection_scheme": {
    "type": "time-based",
    "period_ms": 10000
  },
  "signals_to_collect": [
    {
      "name": "Vehicle.Speed",
      "max_sample_count": 1000,
      "minimum_sampling_interval_ms": 100
    },
    {
      "name": "Vehicle.Powertrain.EngineSpeed",
      "max_sample_count": 1000,
      "minimum_sampling_interval_ms": 100
    }
  ],
  "data_destinations": [
    {
      "type": "s3",
      "s3_bucket_arn": "arn:aws:s3:::my-bucket",
      "s3_prefix": "fleetwise-data/",
      "s3_data_format": "JSON",
      "s3_storage_compression": "GZIP"
    }
  ],
  "compression": "SNAPPY",
  "diagnostics_mode": "SEND_ACTIVE_DTCS",
  "spooling_mode": "TO_DISK",
  "post_trigger_duration_ms": 5000
}
```

#### Create Event-Based Campaign

```bash
POST /api/v1/fleetwise/campaigns?region=us-east-1
Content-Type: application/json

{
  "name": "high-speed-event-collection",
  "description": "Collect data when vehicle exceeds 100 km/h",
  "signal_catalog_arn": "arn:aws:iotfleetwise:...",
  "target_arn": "arn:aws:iotfleetwise:...:fleet/my-fleet",
  "collection_scheme": {
    "type": "condition-based",
    "expression": "$variable.`Vehicle.Speed` > 100",
    "minimum_trigger_interval_ms": 5000,
    "trigger_mode": "ALWAYS"
  },
  "signals_to_collect": [...],
  "data_destinations": [...],
  "compression": "SNAPPY",
  "post_trigger_duration_ms": 10000
}
```

#### Get Campaign

```bash
GET /api/v1/fleetwise/campaigns/my-campaign?region=us-east-1
```

#### Update Campaign (Approve/Suspend/Resume)

```bash
PUT /api/v1/fleetwise/campaigns/my-campaign?region=us-east-1
Content-Type: application/json

{
  "action": "APPROVE"
}
```

Actions: `APPROVE`, `SUSPEND`, `RESUME`, `UPDATE`

#### Delete Campaign

```bash
DELETE /api/v1/fleetwise/campaigns/my-campaign?region=us-east-1
```

#### List Campaigns

```bash
GET /api/v1/fleetwise/campaigns?region=us-east-1
```

### Fleet Management

#### Create a Fleet

```bash
POST /api/v1/fleetwise/fleets
Content-Type: application/json

{
  "fleet_id": "production-fleet",
  "description": "Production vehicle fleet",
  "signal_catalog_arn": "arn:aws:iotfleetwise:...:signal-catalog/my-catalog",
  "region": "us-east-1"
}
```

#### Associate Vehicle to Fleet

```bash
POST /api/v1/fleetwise/fleets/production-fleet/vehicles?region=us-east-1
Content-Type: application/json

{
  "vehicle_name": "vehicle-001"
}
```

---

## Examples

### Example 1: Basic Vehicle Fleet Setup

```bash
# 1. Create fleet
curl -X POST http://localhost:8080/api/v1/fleetwise/fleets \
  -H "Content-Type: application/json" \
  -d '{
    "fleet_id": "test-fleet",
    "description": "Test fleet",
    "signal_catalog_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:signal-catalog/my-catalog",
    "region": "us-east-1"
  }'

# 2. Create vehicles
curl -X POST http://localhost:8080/api/v1/fleetwise/vehicles/batch \
  -H "Content-Type: application/json" \
  -d '{
    "region": "us-east-1",
    "vehicles": [
      {
        "name": "test-vehicle-001",
        "model_manifest_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:model-manifest/my-model",
        "decoder_manifest_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:decoder-manifest/my-decoder",
        "create_iot_thing": true
      }
    ]
  }'

# 3. Associate vehicle to fleet
curl -X POST http://localhost:8080/api/v1/fleetwise/fleets/test-fleet/vehicles?region=us-east-1 \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_name": "test-vehicle-001"
  }'

# 4. Create campaign
curl -X POST http://localhost:8080/api/v1/fleetwise/campaigns?region=us-east-1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-campaign",
    "signal_catalog_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:signal-catalog/my-catalog",
    "target_arn": "arn:aws:iotfleetwise:us-east-1:123456789012:fleet/test-fleet",
    "collection_scheme": {
      "type": "time-based",
      "period_ms": 60000
    },
    "signals_to_collect": [
      {
        "name": "Vehicle.Speed",
        "max_sample_count": 100,
        "minimum_sampling_interval_ms": 1000
      }
    ],
    "compression": "SNAPPY"
  }'
```

### Example 2: Event-Based Data Collection

```bash
# Create campaign that triggers on harsh braking
curl -X POST http://localhost:8080/api/v1/fleetwise/campaigns?region=us-east-1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "harsh-braking-detection",
    "signal_catalog_arn": "arn:aws:iotfleetwise:...",
    "target_arn": "arn:aws:iotfleetwise:...:fleet/my-fleet",
    "collection_scheme": {
      "type": "condition-based",
      "expression": "$variable.`Vehicle.Speed` - last_change($variable.`Vehicle.Speed`, 1000) < -20",
      "minimum_trigger_interval_ms": 2000,
      "trigger_mode": "ALWAYS"
    },
    "signals_to_collect": [
      {
        "name": "Vehicle.Speed",
        "max_sample_count": 500,
        "minimum_sampling_interval_ms": 50
      },
      {
        "name": "Vehicle.Chassis.Brake.PedalPosition",
        "max_sample_count": 500,
        "minimum_sampling_interval_ms": 50
      }
    ],
    "data_destinations": [
      {
        "type": "s3",
        "s3_bucket_arn": "arn:aws:s3:::safety-events",
        "s3_prefix": "harsh-braking/",
        "s3_data_format": "JSON",
        "s3_storage_compression": "GZIP"
      }
    ],
    "post_trigger_duration_ms": 5000
  }'
```

---

## Best Practices

### 1. Resource Naming

- Use consistent naming conventions for vehicles, fleets, and campaigns
- Include environment identifiers (dev, staging, prod)
- Use descriptive names that indicate purpose

Example: `prod-fleet-west-coast`, `vehicle-model3-001`, `campaign-speed-monitoring-prod`

### 2. Signal Selection

- Only collect signals you need to reduce costs
- Set appropriate sampling intervals based on signal importance
- Use condition-based collection for events

### 3. Data Destinations

- Use S3 for long-term storage and batch analytics
- Use Timestream for real-time queries and dashboards
- Use MQTT for real-time integrations and alerts

### 4. Cost Optimization

- Enable compression to reduce data transfer costs
- Use appropriate sampling rates
- Set up S3 lifecycle policies for old data
- Use condition-based campaigns instead of time-based when possible

### 5. Security

- Use IAM roles with least privilege
- Enable encryption at rest and in transit
- Rotate credentials regularly
- Use VPC endpoints for private connectivity

### 6. Monitoring

- Monitor campaign status regularly
- Check vehicle connectivity status
- Set up CloudWatch alarms for failures
- Review data volume and costs weekly

---

## Troubleshooting

### Common Issues

#### 1. Vehicle Creation Fails

**Error**: `Failed to create vehicle: AccessDeniedException`

**Solution**: Ensure IAM credentials have `iotfleetwise:CreateVehicle` and `iot:CreateThing` permissions.

#### 2. Campaign Not Collecting Data

**Error**: Campaign status shows `SUSPENDED`

**Solution**:
1. Approve the campaign: `PUT /api/v1/fleetwise/campaigns/{name}` with action `APPROVE`
2. Check vehicle connectivity
3. Verify signal catalog contains all signals in campaign

#### 3. No Data in S3 Bucket

**Checklist**:
- Campaign is `APPROVED` and `RUNNING`
- Vehicles are connected and reporting status
- S3 bucket permissions allow FleetWise to write
- Data destination is correctly configured
- Wait at least one collection period

#### 4. High Costs

**Actions**:
- Reduce sampling frequency
- Decrease number of signals collected
- Use condition-based collection
- Enable compression
- Review and delete old data

### Debug Commands

```bash
# Check vehicle status
curl http://localhost:8080/api/v1/fleetwise/vehicles/vehicle-001/status?region=us-east-1

# Check campaign details
curl http://localhost:8080/api/v1/fleetwise/campaigns/my-campaign?region=us-east-1

# List all vehicles
curl http://localhost:8080/api/v1/fleetwise/vehicles?region=us-east-1

# Check environment status
curl http://localhost:8080/api/v1/environments/{env-id}/status
```

---

## Additional Resources

- [AWS IoT FleetWise Developer Guide](https://docs.aws.amazon.com/iot-fleetwise/latest/developerguide/)
- [AWS IoT FleetWise API Reference](https://docs.aws.amazon.com/iot-fleetwise/latest/APIReference/)
- [Vehicle Signal Specification (VSS)](https://covesa.github.io/vehicle_signal_specification/)
- [AWS IoT FleetWise Pricing](https://aws.amazon.com/iot-fleetwise/pricing/)
- [SimPlan AWS Automotive Configuration Options](/docs/aws-automotive-configuration-options.md)

---

## Support

For issues or questions:
- GitHub Issues: [SimPlan Repository](https://github.com/your-org/simplan)
- Documentation: `/docs/aws-automotive-configuration-options.md`
- AWS Support: [AWS Support Center](https://console.aws.amazon.com/support/)
