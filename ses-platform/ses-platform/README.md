# SES Platform - Simulation Environment Specification

A comprehensive platform for developers, testers, and DevOps engineers to specify, provision, and manage simulation environments based on configurable capabilities and enablers.

## 🏗️ Architecture Overview

### System Components

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (React)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Dashboard   │  │   Wizard     │  │  Templates   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└────────────────────────────┬────────────────────────────────────┘
                             │ REST API
┌────────────────────────────┴────────────────────────────────────┐
│                    Backend (Go + Gin)                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Environment  │  │  Validation  │  │  Simulation  │          │
│  │   Manager    │  │   Engine     │  │   Engine     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└────────────────────────────┬────────────────────────────────────┘
                             │ SQL
┌────────────────────────────┴────────────────────────────────────┐
│                    Database (PostgreSQL)                         │
│  - Capabilities & Enablers   - State Management                 │
│  - Environments              - Audit Logs                        │
│  - Resource Allocations      - Cost Tracking                     │
└──────────────────────────────────────────────────────────────────┘
```

## 📋 Features

### ✅ Core Capabilities (18 Total)
- **C01**: Spec Authoring & Validation
- **C02**: Parsing & Internal Modeling
- **C03**: Planning Engine
- **C04**: Provisioning Automation
- **C05**: Orchestration & State Machine
- **C06**: Monitoring & Metrics
- **C07**: Logging & Audit
- **C08**: Cost Management
- **C09**: Security & Compliance
- **C10**: Spec-Kit & Reuse
- **C11**: Reservation & Scheduling
- **C12**: Simulation Execution
- **C13**: Error Handling & Recovery
- **C14**: Provider Abstraction
- **C15**: Environment Visualization
- **C16**: Messaging & Agent Coordination
- **C17**: State Persistence Layer
- **C18**: Access & Governance

### 🔧 Enablers (20 Total)
Each capability is powered by one or more enablers providing foundational services like:
- Core Platform Infrastructure
- Schema Validation
- Graph & Planning Algorithms
- Execution Framework
- Provider SDK Integrations
- Metrics & Logging Stacks
- Security/RBAC Components
- And more...

## 🚀 Getting Started

### Prerequisites

- **Node.js** 18+ and npm/yarn
- **Go** 1.21+
- **PostgreSQL** 14+
- **Git**

### 1. Database Setup

```bash
# Create PostgreSQL database and user
psql -U postgres

CREATE DATABASE ses_platform;
CREATE USER ses_user WITH PASSWORD 'ses_password';
GRANT ALL PRIVILEGES ON DATABASE ses_platform TO ses_user;

# Run schema creation
psql -U ses_user -d ses_platform -f database/schema.sql
```

### 2. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Initialize Go module
go mod init ses-platform

# Install dependencies
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Update database connection in main.go
# Edit the DSN string with your PostgreSQL credentials

# Run the backend
go run main.go
```

The API will start on `http://localhost:8080`

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm start
```

The frontend will start on `http://localhost:3000`

## 📡 API Endpoints

### Environments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/environments` | List all environments |
| POST | `/api/v1/environments` | Create new environment |
| GET | `/api/v1/environments/:id` | Get environment details |
| PUT | `/api/v1/environments/:id` | Update environment |
| DELETE | `/api/v1/environments/:id` | Delete environment |

### Environment Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/environments/:id/provision` | Provision environment |
| POST | `/api/v1/environments/:id/start` | Start environment |
| POST | `/api/v1/environments/:id/stop` | Stop environment |
| POST | `/api/v1/environments/:id/upload` | Upload binary/config |
| GET | `/api/v1/environments/:id/status` | Get current status |
| GET | `/api/v1/environments/:id/metrics` | Get metrics data |
| GET | `/api/v1/environments/:id/logs` | Get logs |

### Validation & Cost

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/validate` | Validate spec |
| POST | `/api/v1/cost/estimate` | Estimate cost |

### Reference Data

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/capabilities` | List all capabilities |
| GET | `/api/v1/enablers` | List all enablers |
| GET | `/api/v1/templates` | List templates |

### Audit

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/audit` | Get audit logs |
| GET | `/api/v1/environments/:id/history` | Get environment history |

## 🔄 Workflow

### Creating a New Environment

1. **Step 1: Project Details**
   - Enter name, description, owner, and tags

2. **Step 2: Capability Selection**
   - Choose from 18 available capabilities
   - System validates dependencies automatically
   - Blocked capabilities show dependency requirements

3. **Step 3: Enabler Configuration**
   - Auto-populated based on selected capabilities
   - Review required enablers
   - System ensures all dependencies are met

4. **Step 4: Resource Requirements**
   - Define compute (CPU, memory, instances)
   - Specify storage requirements
   - Configure network settings
   - Set scheduling priority and duration

5. **Step 5: Review & Submit**
   - Review all configurations
   - See estimated daily cost
   - Submit for provisioning

### Provisioning Flow

```
pending → validating → provisioning → running
                ↓
              error → rollback
```

### State Management

The system tracks all state transitions in the database:
- Each transition logged with reason and metadata
- Full audit trail for compliance
- Rollback capability to previous states

## 💾 Database Schema

### Core Tables

- **capabilities**: All 18 capabilities with dependencies
- **enablers**: All 20 enablers with descriptions
- **environments**: User-created environments
- **state_transitions**: Lifecycle state changes
- **resource_allocations**: Compute/storage/network resources
- **reservations**: Scheduling and time windows
- **metrics_snapshots**: Time-series monitoring data
- **cost_records**: Detailed cost tracking
- **audit_logs**: Compliance and governance
- **structured_logs**: Application logs
- **uploads**: Binary and configuration files
- **error_records**: Failure tracking
- **rollback_operations**: Recovery operations

## 🎯 Use Cases

### For Developers
- Quickly spin up development environments
- Test code against specific capability configurations
- Upload and deploy new binaries for testing

### For Testers
- Create isolated test environments
- Run load, stress, and integration tests
- Validate against multiple configurations

### For DevOps Engineers
- Provision infrastructure on-demand
- Monitor resource utilization and costs
- Manage environment lifecycle and compliance

## 🔒 Security Features

- **RBAC**: Role-based access control (admin, operator, developer, viewer)
- **Audit Logging**: Complete audit trail of all operations
- **Compliance Checks**: Automated policy validation
- **Encryption**: Data protection at rest and in transit
- **MFA Support**: Multi-factor authentication ready

## 💰 Cost Management

- **Pre-provisioning Estimation**: See costs before creating
- **Real-time Tracking**: Monitor actual costs as they accrue
- **Cost Breakdown**: By component (compute, storage, capabilities)
- **Optimization Tips**: Suggestions for cost reduction
- **Budget Alerts**: Threshold-based notifications

## 📊 Monitoring & Observability

### Metrics Collected
- CPU, memory, disk, network utilization
- Request count and error rates
- Latency percentiles (P50, P95, P99)
- Cost per hour
- Health scores

### Logging
- Structured logs with trace IDs
- Distributed tracing support
- Multiple severity levels
- Full-text search capability

## 🔄 Future Integration Points

The current implementation uses simulated provisioning. To integrate with real backend systems:

### 1. Replace Simulation Functions

In `main.go`, replace:
```go
func simulateProvisioning(envID string) {
    // Current: Simulated delays and state changes
}
```

With actual provider SDK calls:
```go
func provisionEnvironment(envID string) {
    // Call AWS/Azure/GCP APIs
    // Update state based on actual provisioning
}
```

### 2. Add Provider Integrations

Implement the Provider Abstraction layer (C14):
- AWS SDK integration
- Azure SDK integration
- GCP SDK integration
- On-premises orchestration

### 3. Connect Monitoring Systems

Integrate with existing monitoring:
- Prometheus for metrics collection
- ELK/Splunk for log aggregation
- Jaeger/Zipkin for distributed tracing

### 4. Implement Message Queue

Add real message broker for agent coordination (C16):
- RabbitMQ / Kafka integration
- Event-driven architecture
- Async job processing

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Integration Tests
```bash
# Start all services
docker-compose up -d

# Run integration test suite
./scripts/integration-tests.sh
```

## 📦 Deployment

### Docker Deployment

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: ses_platform
      POSTGRES_USER: ses_user
      POSTGRES_PASSWORD: ses_password
    volumes:
      - ./database/schema.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
    ports:
      - "8080:8080"
    depends_on:
      - postgres

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend
```

### Kubernetes Deployment

```bash
# Deploy to Kubernetes cluster
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n ses-platform
```

## 📚 Additional Resources

### Documentation Structure
```
docs/
├── architecture/
│   ├── system-design.md
│   ├── capability-reference.md
│   └── enabler-reference.md
├── api/
│   ├── rest-api.md
│   └── openapi.yaml
├── guides/
│   ├── getting-started.md
│   ├── user-guide.md
│   └── admin-guide.md
└── development/
    ├── contributing.md
    ├── coding-standards.md
    └── testing-guide.md
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

Built using:
- React + TypeScript for frontend
- Go + Gin for backend
- PostgreSQL for database
- Tailwind CSS for styling
- Lucide React for icons

---

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432

# Verify user permissions
psql -U ses_user -d ses_platform -c "\dt"
```

### Backend Not Starting
```bash
# Check Go version
go version

# Verify dependencies
go mod tidy
go mod download
```

### Frontend Build Errors
```bash
# Clear cache
rm -rf node_modules package-lock.json
npm install

# Check Node version
node --version  # Should be 18+
```

## 📞 Support

For issues and questions:
- GitHub Issues: [github.com/yourorg/ses-platform/issues](https://github.com/yourorg/ses-platform/issues)
- Documentation: [docs.ses-platform.io](https://docs.ses-platform.io)
- Email: support@ses-platform.io
