# AWS IoT FleetWise Configuration Options

This document provides a comprehensive list of all configurable options available in AWS IoT FleetWise for automotive vehicle simulation and data collection.

---

## Table of Contents
1. [Signal Catalog Configuration](#signal-catalog-configuration)
2. [Vehicle Configuration](#vehicle-configuration)
3. [Model Manifest Configuration](#model-manifest-configuration)
4. [Decoder Manifest Configuration](#decoder-manifest-configuration)
5. [Campaign Configuration](#campaign-configuration)
6. [Fleet Configuration](#fleet-configuration)
7. [Data Destinations](#data-destinations)
8. [Collection Schemes](#collection-schemes)

---

## 1. Signal Catalog Configuration

### CreateSignalCatalog API Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String | Yes | Unique name for the signal catalog |
| `description` | String | No | Description of the signal catalog |
| `nodes` | Array | No | Signal definitions (attributes, branches, sensors, actuators) |
| `tags` | Array | No | Metadata tags for organization and cost allocation |

### Signal Node Types

#### Attribute Signals (Static Information)
```json
{
  "fullyQualifiedName": "Vehicle.VIN",
  "attribute": {
    "fullyQualifiedName": "Vehicle.VIN",
    "dataType": "STRING",
    "defaultValue": "",
    "description": "Vehicle Identification Number",
    "unit": "",
    "allowedValues": [],
    "min": 0,
    "max": 0
  }
}
```

**Supported Data Types:**
- `STRING`
- `INT8`, `UINT8`, `INT16`, `UINT16`, `INT32`, `UINT32`, `INT64`, `UINT64`
- `FLOAT`, `DOUBLE`
- `BOOLEAN`
- `UNIX_TIMESTAMP`
- `INT8_ARRAY`, `UINT8_ARRAY`, `INT16_ARRAY`, `UINT16_ARRAY`, `INT32_ARRAY`, `UINT32_ARRAY`, `INT64_ARRAY`, `UINT64_ARRAY`
- `FLOAT_ARRAY`, `DOUBLE_ARRAY`
- `BOOLEAN_ARRAY`
- `UNIX_TIMESTAMP_ARRAY`
- `UNKNOWN`

#### Branch Signals (Hierarchical Structure)
```json
{
  "fullyQualifiedName": "Vehicle.Powertrain",
  "branch": {
    "fullyQualifiedName": "Vehicle.Powertrain",
    "description": "Powertrain branch containing engine and transmission signals"
  }
}
```

#### Sensor Signals (Dynamic Data)
```json
{
  "fullyQualifiedName": "Vehicle.Speed",
  "sensor": {
    "fullyQualifiedName": "Vehicle.Speed",
    "dataType": "DOUBLE",
    "description": "Vehicle speed in km/h",
    "unit": "km/h",
    "min": 0,
    "max": 300
  }
}
```

#### Actuator Signals (Device States)
```json
{
  "fullyQualifiedName": "Vehicle.Door.Driver.IsLocked",
  "actuator": {
    "fullyQualifiedName": "Vehicle.Door.Driver.IsLocked",
    "dataType": "BOOLEAN",
    "description": "Driver door lock status",
    "allowedValues": ["true", "false"]
  }
}
```

### UpdateSignalCatalog Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String | Yes | Name of the signal catalog to update |
| `description` | String | No | Updated description |
| `nodesToAdd` | Array | No | New signal nodes to add |
| `nodesToUpdate` | Array | No | Existing nodes to modify |
| `nodesToRemove` | Array | No | Nodes to remove (by fully qualified name) |

---

## 2. Vehicle Configuration

### CreateVehicle API Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `vehicleName` | String | Yes | Unique identifier for the vehicle (max 100 characters) |
| `modelManifestArn` | String | Yes | ARN of the vehicle model manifest |
| `decoderManifestArn` | String | Yes | ARN of the decoder manifest for signal decoding |
| `attributes` | Map | No | Key-value pairs of vehicle attributes |
| `associationBehavior` | String | No | `CreateIotThing` or `ValidateIotThingExists` |
| `tags` | Array | No | Resource tags for organization |
| `stateTemplates` | Array | No | State template configurations |

### Vehicle Attributes (Custom Key-Value Pairs)
```json
{
  "attributes": {
    "Make": "Tesla",
    "Model": "Model 3",
    "Year": "2024",
    "VIN": "5YJ3E1EA8KF123456",
    "Color": "Pearl White",
    "TrimLevel": "Long Range AWD",
    "BatteryCapacity": "75kWh",
    "Odometer": "15234",
    "Owner": "Fleet-001",
    "Location": "US-West-2",
    "FirmwareVersion": "2024.2.10"
  }
}
```

### Association Behavior Options
- **`CreateIotThing`**: Creates a new AWS IoT Thing when creating the vehicle
- **`ValidateIotThingExists`**: Validates an existing AWS IoT Thing as a vehicle

### BatchCreateVehicle Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `vehicles` | Array | Yes | Array of vehicle configurations (up to 10 per batch) |

### UpdateVehicle Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `vehicleName` | String | Yes | Name of vehicle to update |
| `modelManifestArn` | String | No | Updated model manifest ARN |
| `decoderManifestArn` | String | No | Updated decoder manifest ARN |
| `attributes` | Map | No | Updated attributes |
| `attributeUpdateMode` | String | No | `Overwrite` or `Merge` |

---

## 3. Model Manifest Configuration

### CreateModelManifest API Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String | Yes | Unique name for the model manifest |
| `description` | String | No | Description of the vehicle model |
| `nodes` | Array | Yes | List of signal node paths from signal catalog |
| `signalCatalogArn` | String | Yes | ARN of the signal catalog to use |
| `tags` | Array | No | Resource tags |

### Model Status
- `DRAFT`: Editable state
- `ACTIVE`: Published and in use
- `INVALID`: Contains validation errors

### UpdateModelManifest Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String | Yes | Model manifest name |
| `description` | String | No | Updated description |
| `nodesToAdd` | Array | No | Signal nodes to add |
| `nodesToRemove` | Array | No | Signal nodes to remove |
| `status` | String | No | Update status (DRAFT or ACTIVE) |

---

## 4. Decoder Manifest Configuration

### CreateDecoderManifest API Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String | Yes | Unique name for decoder manifest |
| `description` | String | No | Description |
| `modelManifestArn` | String | Yes | ARN of associated model manifest |
| `signalDecoders` | Array | No | Signal decoding rules |
| `networkInterfaces` | Array | No | Network interface definitions |
| `tags` | Array | No | Resource tags |

### Network Interface Types

#### CAN Interface
```json
{
  "interfaceId": "1",
  "type": "CAN_INTERFACE",
  "canInterface": {
    "name": "can0",
    "protocolName": "CAN",
    "protocolVersion": "2.0A"
  }
}
```

#### OBD Interface
```json
{
  "interfaceId": "2",
  "type": "OBD_INTERFACE",
  "obdInterface": {
    "name": "obd",
    "requestMessageId": 2015,
    "obdStandard": "J1979",
    "pidRequestIntervalSeconds": 5,
    "dtcRequestIntervalSeconds": 10,
    "useExtendedIds": false,
    "hasTransmissionEcu": true
  }
}
```

#### Vehicle Middleware Interface
```json
{
  "interfaceId": "3",
  "type": "VEHICLE_MIDDLEWARE",
  "vehicleMiddleware": {
    "name": "ROS2",
    "protocolName": "ROS_2"
  }
}
```

### Signal Decoder Types

#### CAN Signal Decoder
```json
{
  "fullyQualifiedName": "Vehicle.Speed",
  "type": "CAN_SIGNAL",
  "interfaceId": "1",
  "canSignal": {
    "messageId": 123,
    "isBigEndian": false,
    "isSigned": false,
    "startBit": 0,
    "offset": 0.0,
    "factor": 0.1,
    "length": 16,
    "name": "VehicleSpeed"
  }
}
```

#### OBD Signal Decoder
```json
{
  "fullyQualifiedName": "Vehicle.Powertrain.EngineSpeed",
  "type": "OBD_SIGNAL",
  "interfaceId": "2",
  "obdSignal": {
    "pidResponseLength": 2,
    "serviceMode": 1,
    "pid": 12,
    "scaling": 0.25,
    "offset": 0.0,
    "startByte": 0,
    "byteLength": 2,
    "bitRightShift": 0,
    "bitMaskLength": 16
  }
}
```

#### Message Signal Decoder
```json
{
  "fullyQualifiedName": "Vehicle.Chassis.SteeringAngle",
  "type": "MESSAGE_SIGNAL",
  "interfaceId": "3",
  "messageSignal": {
    "topicName": "/vehicle/steering",
    "structuredMessage": {
      "primitiveMessageDefinition": {
        "ros2PrimitiveMessageDefinition": {
          "primitiveType": "FLOAT64",
          "offset": 0.0,
          "scaling": 1.0,
          "upperBound": 0
        }
      }
    }
  }
}
```

---

## 5. Campaign Configuration

### CreateCampaign API Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String | Yes | Unique campaign name |
| `description` | String | No | Campaign description |
| `signalCatalogArn` | String | Yes | ARN of signal catalog |
| `targetArn` | String | Yes | ARN of vehicle or fleet |
| `collectionScheme` | Object | Yes | Data collection rules (time-based or event-based) |
| `signalsToCollect` | Array | Yes | List of signals to collect |
| `dataDestinationConfigs` | Array | No | Where to send collected data |
| `compression` | String | No | Compression algorithm: `OFF` or `SNAPPY` (default) |
| `priority` | Integer | No | Campaign priority (0-10, default 0) |
| `startTime` | Timestamp | No | Campaign start time |
| `expiryTime` | Timestamp | No | Campaign expiration time |
| `postTriggerCollectionDuration` | Long | No | Milliseconds to collect after trigger (0-4294967295) |
| `diagnosticsMode` | String | No | `OFF` or `SEND_ACTIVE_DTCS` |
| `spoolingMode` | String | No | `OFF` or `TO_DISK` |
| `dataExtraDimensions` | Array | No | Additional vehicle attributes to include |
| `dataPartitions` | Array | No | Data partition configurations |
| `tags` | Array | No | Resource tags |

### Compression Options
- **`OFF`**: No compression (faster, larger size)
- **`SNAPPY`**: Fast compression algorithm (default, good balance)

### Diagnostics Mode Options
- **`OFF`**: Do not collect diagnostic trouble codes (default)
- **`SEND_ACTIVE_DTCS`**: Include active DTCs in collected data

### Spooling Mode Options
- **`OFF`**: Do not store data when offline (default)
- **`TO_DISK`**: Store data to disk when vehicle loses connectivity, upload when reconnected

### Signals to Collect
```json
{
  "signalsToCollect": [
    {
      "name": "Vehicle.Speed",
      "maxSampleCount": 1000,
      "minimumSamplingIntervalMs": 100
    },
    {
      "name": "Vehicle.Powertrain.EngineSpeed",
      "maxSampleCount": 1000,
      "minimumSamplingIntervalMs": 100
    },
    {
      "name": "Vehicle.Chassis.SteeringAngle",
      "maxSampleCount": 500,
      "minimumSamplingIntervalMs": 200
    }
  ]
}
```

### Collection Scheme Types

#### Time-Based Collection
```json
{
  "collectionScheme": {
    "timeBasedCollectionScheme": {
      "periodMs": 10000
    }
  }
}
```
Collects data at fixed intervals (in milliseconds).

#### Condition-Based Collection (Event-Driven)
```json
{
  "collectionScheme": {
    "conditionBasedCollectionScheme": {
      "expression": "$variable.`Vehicle.Speed` > 100",
      "minimumTriggerIntervalMs": 5000,
      "triggerMode": "ALWAYS",
      "conditionLanguageVersion": 1
    }
  }
}
```

**Trigger Mode Options:**
- `ALWAYS`: Trigger every time condition is met
- `RISING_EDGE`: Trigger only when condition changes from false to true

**Expression Language:**
- Comparison operators: `>`, `<`, `>=`, `<=`, `==`, `!=`
- Logical operators: `&&`, `||`, `!`
- Arithmetic operators: `+`, `-`, `*`, `/`
- Signal reference: `` $variable.`Signal.Path` ``
- Functions: `last_change()`, `timeout()`, `avg()`, `min()`, `max()`

**Example Complex Expressions:**
```
# Speed over 100 AND engine RPM over 5000
$variable.`Vehicle.Speed` > 100 && $variable.`Vehicle.Powertrain.EngineSpeed` > 5000

# Temperature exceeds threshold
$variable.`Vehicle.Powertrain.Temperature` > 95

# Rapid deceleration
$variable.`Vehicle.Speed` - last_change($variable.`Vehicle.Speed`, 1000) < -20

# Low battery
$variable.`Vehicle.Powertrain.Battery.StateOfCharge` < 20
```

### Data Destination Configurations

#### Amazon S3 Destination
```json
{
  "dataDestinationConfigs": [
    {
      "s3Config": {
        "bucketArn": "arn:aws:s3:::my-fleetwise-data",
        "dataFormat": "JSON",
        "storageCompressionFormat": "GZIP",
        "prefix": "vehicle-data/"
      }
    }
  ]
}
```

**Data Format Options:**
- `JSON`: Human-readable JSON format
- `PARQUET`: Columnar format optimized for analytics

**Storage Compression Options:**
- `NONE`: No compression
- `GZIP`: Standard gzip compression

#### Amazon Timestream Destination
```json
{
  "dataDestinationConfigs": [
    {
      "timestreamConfig": {
        "timestreamTableArn": "arn:aws:timestream:us-east-1:123456789012:database/FleetwiseDB/table/VehicleData",
        "executionRoleArn": "arn:aws:iam::123456789012:role/FleetwiseTimestreamRole"
      }
    }
  ]
}
```

#### MQTT Topic Destination
```json
{
  "dataDestinationConfigs": [
    {
      "mqttTopicConfig": {
        "mqttTopicArn": "arn:aws:iot:us-east-1:123456789012:topic/fleetwise/vehicle/data",
        "executionRoleArn": "arn:aws:iam::123456789012:role/FleetwiseMqttRole"
      }
    }
  ]
}
```

### Data Extra Dimensions (Enrichment)
Add vehicle attributes as dimensions to collected data:
```json
{
  "dataExtraDimensions": [
    "Vehicle.VIN",
    "Vehicle.Make",
    "Vehicle.Model"
  ]
}
```

### Data Partitions
Organize collected data into partitions:
```json
{
  "dataPartitions": [
    {
      "id": "partition-1",
      "storageOptions": {
        "maximumSize": {
          "unit": "MB",
          "value": 100
        },
        "storageLocation": "arn:aws:s3:::my-bucket/partition-1/",
        "minimumTimeToLive": {
          "unit": "HOURS",
          "value": 24
        }
      }
    }
  ]
}
```

---

## 6. Fleet Configuration

### CreateFleet API Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `fleetId` | String | Yes | Unique fleet identifier |
| `description` | String | No | Fleet description |
| `signalCatalogArn` | String | Yes | ARN of signal catalog |
| `tags` | Array | No | Resource tags |

### AssociateVehicleFleet Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `vehicleName` | String | Yes | Vehicle to add to fleet |
| `fleetId` | String | Yes | Target fleet ID |

---

## 7. Data Destinations

### Supported Destinations
1. **Amazon S3**: For long-term storage and batch analytics
2. **Amazon Timestream**: For time-series data and real-time queries
3. **AWS IoT Core MQTT**: For real-time streaming and pub/sub patterns

### Multi-Destination Support
Campaigns can send data to multiple destinations simultaneously:
```json
{
  "dataDestinationConfigs": [
    {
      "s3Config": { /* S3 configuration */ }
    },
    {
      "timestreamConfig": { /* Timestream configuration */ }
    },
    {
      "mqttTopicConfig": { /* MQTT configuration */ }
    }
  ]
}
```

---

## 8. Collection Schemes

### Time-Based Collection Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `periodMs` | Long | Yes | Collection interval in milliseconds (10000-60000) |

### Condition-Based Collection Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `expression` | String | Yes | Boolean expression defining trigger condition |
| `minimumTriggerIntervalMs` | Long | No | Minimum time between triggers (0-4294967295) |
| `triggerMode` | String | No | `ALWAYS` or `RISING_EDGE` |
| `conditionLanguageVersion` | Integer | No | Version of expression language (default: 1) |

---

## 9. Fleet Size Configuration

### FleetSize Parameter
Control the number of simulated vehicles in your environment:
- Specified in CloudFormation templates
- Can be updated to scale fleet up or down
- Affects cost and data volume

### Example Configuration
```yaml
Parameters:
  FleetSize:
    Type: Number
    Default: 10
    MinValue: 1
    MaxValue: 1000
    Description: Number of vehicles to simulate
```

---

## 10. Edge Agent Configuration

### Connection Types
- **`iotCore`**: Direct connection to AWS IoT Core
- **`iotGreengrassV2`**: Connection via AWS IoT Greengrass

### Configuration File Options
```json
{
  "networkInterfaces": [],
  "staticConfig": {
    "mqttConnection": {
      "endpointUrl": "xxxxxx.iot.us-east-1.amazonaws.com",
      "clientId": "my-vehicle-001",
      "collectionSchemeListTopic": "$aws/iotfleetwise/vehicles/my-vehicle-001/collection_schemes",
      "decoderManifestTopic": "$aws/iotfleetwise/vehicles/my-vehicle-001/decoder_manifest",
      "canDataTopic": "$aws/iotfleetwise/vehicles/my-vehicle-001/signals",
      "checkinTopic": "$aws/iotfleetwise/vehicles/my-vehicle-001/checkin"
    },
    "bufferSizes": {
      "dtcBufferSize": 100,
      "socketCANBufferSize": 10000,
      "decodedSignalsBufferSize": 10000
    },
    "threadIdleTimes": {
      "inspectionThreadIdleTimeMs": 50,
      "socketCANThreadIdleTimeMs": 50,
      "canDecoderThreadIdleTimeMs": 50
    },
    "persistencyPath": "/var/aws-iot-fleetwise/",
    "logColor": "Auto",
    "logLevel": "Trace"
  }
}
```

---

## 11. Resource Tagging

All FleetWise resources support tagging for:
- Cost allocation
- Resource organization
- Access control
- Automation

### Tag Structure
```json
{
  "tags": [
    {
      "Key": "Environment",
      "Value": "Production"
    },
    {
      "Key": "Project",
      "Value": "ConnectedVehicle"
    },
    {
      "Key": "CostCenter",
      "Value": "Engineering"
    },
    {
      "Key": "Owner",
      "Value": "fleet-team@example.com"
    }
  ]
}
```

---

## 12. State Templates

Define expected states for vehicle signals:

```json
{
  "stateTemplates": [
    {
      "identifier": "normalOperation",
      "stateTemplateProperties": [
        {
          "fullyQualifiedName": "Vehicle.Speed",
          "stateValue": {
            "doubleValue": 60.0
          }
        }
      ]
    }
  ]
}
```

---

## Summary of Key Configuration Dimensions

| Category | Options Available |
|----------|-------------------|
| **Signal Types** | Attributes, Branches, Sensors, Actuators |
| **Data Types** | 20+ types including primitives, arrays, timestamps |
| **Network Protocols** | CAN, OBD-II, Vehicle Middleware (ROS2, SOME/IP) |
| **Collection Modes** | Time-based, Condition-based (event-driven) |
| **Compression** | OFF, SNAPPY |
| **Data Destinations** | S3 (JSON/Parquet), Timestream, MQTT |
| **Storage Formats** | JSON, Parquet with GZIP compression |
| **Diagnostics** | OFF, SEND_ACTIVE_DTCS |
| **Spooling** | OFF, TO_DISK |
| **Fleet Scaling** | 1 to 1000+ vehicles |
| **Trigger Modes** | ALWAYS, RISING_EDGE |

---

## Best Practices

1. **Signal Catalog Design**
   - Use VSS standard for signal naming
   - Organize signals in logical hierarchies
   - Define appropriate data types and ranges

2. **Campaign Configuration**
   - Start with time-based collection for testing
   - Use condition-based collection for production efficiency
   - Enable spooling for vehicles with intermittent connectivity
   - Set appropriate sampling rates to balance data quality and cost

3. **Data Destinations**
   - Use S3 for long-term storage and compliance
   - Use Timestream for real-time dashboards and queries
   - Use MQTT for real-time integrations and alerts

4. **Resource Management**
   - Tag all resources consistently
   - Use fleets to manage vehicle groups
   - Monitor campaign data volume and costs

5. **Security**
   - Use IAM roles with least privilege
   - Enable encryption at rest and in transit
   - Regularly rotate credentials

---

## Cost Optimization Tips

1. **Sampling Rates**: Adjust `minimumSamplingIntervalMs` based on signal importance
2. **Compression**: Enable SNAPPY compression for campaigns
3. **Selective Collection**: Only collect signals needed for specific use cases
4. **Condition-Based Triggers**: Use event-based collection to reduce data volume
5. **Data Lifecycle**: Configure S3 lifecycle policies to archive old data
6. **Fleet Sizing**: Right-size fleet for testing vs production

---

## Additional Resources

- AWS IoT FleetWise Developer Guide: https://docs.aws.amazon.com/iot-fleetwise/latest/developerguide/
- AWS IoT FleetWise API Reference: https://docs.aws.amazon.com/iot-fleetwise/latest/APIReference/
- Vehicle Signal Specification (VSS): https://covesa.github.io/vehicle_signal_specification/
- AWS IoT FleetWise Edge Agent: https://github.com/aws/aws-iot-fleetwise-edge
