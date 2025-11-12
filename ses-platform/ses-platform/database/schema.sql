-- SES Platform Database Schema
-- PostgreSQL 14+

-- Create database
CREATE DATABASE ses_platform;
\c ses_platform;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================
-- CORE REFERENCE TABLES
-- =============================================

-- Capabilities table (C01-C18)
CREATE TABLE capabilities (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enablers JSONB NOT NULL DEFAULT '[]', -- Array of enabler IDs
    dependencies JSONB NOT NULL DEFAULT '[]', -- Array of capability IDs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_capabilities_name ON capabilities(name);

-- Enablers table (E01-E20)
CREATE TABLE enablers (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    supported_capabilities JSONB NOT NULL DEFAULT '[]', -- Array of capability IDs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_enablers_name ON enablers(name);

-- =============================================
-- ENVIRONMENT MANAGEMENT
-- =============================================

-- Environments table - Core entity for simulation environments
CREATE TABLE environments (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner VARCHAR(255),
    tags TEXT[], -- Array of tags
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
        -- Valid statuses: pending, validating, provisioning, running, stopped, error, deleted
    capabilities JSONB NOT NULL DEFAULT '[]', -- Selected capability IDs
    enablers_config JSONB NOT NULL DEFAULT '{}', -- Enabler configurations
    
    -- Resource Configuration
    compute_config JSONB NOT NULL DEFAULT '{}',
        -- {cpu: int, memory: int, instances: int}
    storage INTEGER NOT NULL DEFAULT 0, -- GB
    network VARCHAR(50) DEFAULT 'private',
        -- Valid: private, public, hybrid
    
    -- Scheduling
    priority VARCHAR(20) DEFAULT 'medium',
        -- Valid: low, medium, high, critical
    duration INTEGER DEFAULT 24, -- hours
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    
    -- Cost Tracking
    estimated_cost DECIMAL(10, 2) DEFAULT 0.00,
    actual_cost DECIMAL(10, 2) DEFAULT 0.00,
    
    -- Health & Monitoring
    health INTEGER DEFAULT 100, -- 0-100 percentage
    uptime VARCHAR(50) DEFAULT '0h',
    last_health_check TIMESTAMP WITH TIME ZONE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT chk_status CHECK (status IN ('pending', 'validating', 'provisioning', 'running', 'stopped', 'error', 'deleted')),
    CONSTRAINT chk_priority CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_health CHECK (health >= 0 AND health <= 100)
);

CREATE INDEX idx_environments_status ON environments(status);
CREATE INDEX idx_environments_owner ON environments(owner);
CREATE INDEX idx_environments_created_at ON environments(created_at DESC);
CREATE INDEX idx_environments_tags ON environments USING GIN(tags);

-- =============================================
-- STATE MANAGEMENT (C05, C17)
-- =============================================

-- State transitions - Track all lifecycle changes
CREATE TABLE state_transitions (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    from_state VARCHAR(50) NOT NULL,
    to_state VARCHAR(50) NOT NULL,
    reason TEXT,
    metadata JSONB DEFAULT '{}',
    user_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_state_transitions_env ON state_transitions(environment_id, created_at DESC);
CREATE INDEX idx_state_transitions_states ON state_transitions(from_state, to_state);

-- =============================================
-- RESOURCE ALLOCATION (C04, C11)
-- =============================================

-- Resource allocations - Track specific resource assignments
CREATE TABLE resource_allocations (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL,
        -- Valid: compute, storage, network, database, cache
    resource_id VARCHAR(255) NOT NULL, -- Provider-specific resource identifier
    provider VARCHAR(50) NOT NULL, -- aws, azure, gcp, on-prem
    region VARCHAR(50),
    status VARCHAR(50) DEFAULT 'allocating',
        -- Valid: allocating, active, releasing, released, failed
    config JSONB NOT NULL DEFAULT '{}',
    allocated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    released_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT chk_resource_status CHECK (status IN ('allocating', 'active', 'releasing', 'released', 'failed'))
);

CREATE INDEX idx_resource_allocations_env ON resource_allocations(environment_id);
CREATE INDEX idx_resource_allocations_status ON resource_allocations(status);

-- Reservations - Scheduling and resource locking (C11)
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(50) DEFAULT 'scheduled',
        -- Valid: scheduled, active, completed, cancelled, preempted
    conflict_resolution VARCHAR(50) DEFAULT 'queue',
        -- Valid: queue, preempt, fail
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_reservation_times CHECK (end_time > start_time),
    CONSTRAINT chk_reservation_status CHECK (status IN ('scheduled', 'active', 'completed', 'cancelled', 'preempted'))
);

CREATE INDEX idx_reservations_env ON reservations(environment_id);
CREATE INDEX idx_reservations_time ON reservations(start_time, end_time);
CREATE INDEX idx_reservations_status ON reservations(status);

-- =============================================
-- MONITORING & METRICS (C06, C08)
-- =============================================

-- Metrics snapshots - Time-series data aggregated periodically
CREATE TABLE metrics_snapshots (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- System Metrics
    cpu_usage DECIMAL(5, 2), -- Percentage
    memory_usage DECIMAL(5, 2), -- Percentage
    disk_usage DECIMAL(5, 2), -- Percentage
    network_in DECIMAL(15, 2), -- MB
    network_out DECIMAL(15, 2), -- MB
    
    -- Application Metrics
    request_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    latency_p50 DECIMAL(10, 2), -- Milliseconds
    latency_p95 DECIMAL(10, 2),
    latency_p99 DECIMAL(10, 2),
    
    -- Cost Metrics
    hourly_cost DECIMAL(10, 4),
    
    -- Custom metrics (flexible)
    custom_metrics JSONB DEFAULT '{}',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_metrics_env_time ON metrics_snapshots(environment_id, timestamp DESC);
CREATE INDEX idx_metrics_timestamp ON metrics_snapshots(timestamp DESC);

-- Cost tracking - Detailed cost breakdown (C08)
CREATE TABLE cost_records (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    cost_type VARCHAR(50) NOT NULL,
        -- Valid: compute, storage, network, data_transfer, licensing, other
    amount DECIMAL(10, 4) NOT NULL,
    unit VARCHAR(20), -- e.g., per_hour, per_gb, per_request
    resource_id VARCHAR(255),
    provider VARCHAR(50),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_cost_amount CHECK (amount >= 0)
);

CREATE INDEX idx_cost_records_env ON cost_records(environment_id, timestamp DESC);
CREATE INDEX idx_cost_records_type ON cost_records(cost_type);

-- =============================================
-- LOGGING & AUDIT (C07)
-- =============================================

-- Audit logs - Compliance and governance tracking
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) REFERENCES environments(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
        -- e.g., created, updated, deleted, provisioned, started, stopped
    entity_type VARCHAR(50), -- environment, resource, reservation, etc.
    entity_id VARCHAR(255),
    user_id VARCHAR(255) NOT NULL,
    user_email VARCHAR(255),
    ip_address INET,
    details JSONB DEFAULT '{}',
    result VARCHAR(50) DEFAULT 'success',
        -- Valid: success, failure, partial
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_audit_result CHECK (result IN ('success', 'failure', 'partial'))
);

CREATE INDEX idx_audit_logs_env ON audit_logs(environment_id, created_at DESC);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(created_at DESC);

-- Structured logs - Application and system logs
CREATE TABLE structured_logs (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) REFERENCES environments(id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    level VARCHAR(20) NOT NULL,
        -- Valid: DEBUG, INFO, WARN, ERROR, FATAL
    source VARCHAR(100), -- Component or service name
    message TEXT NOT NULL,
    trace_id VARCHAR(100), -- For distributed tracing
    span_id VARCHAR(100),
    error_code VARCHAR(50),
    stack_trace TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_log_level CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL'))
);

CREATE INDEX idx_structured_logs_env ON structured_logs(environment_id, timestamp DESC);
CREATE INDEX idx_structured_logs_level ON structured_logs(level, timestamp DESC);
CREATE INDEX idx_structured_logs_trace ON structured_logs(trace_id);

-- =============================================
-- ARTIFACT MANAGEMENT (C12)
-- =============================================

-- Uploads - Binary and configuration artifacts
CREATE TABLE uploads (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
        -- Valid: binary, config, script, data, other
    version VARCHAR(50),
    size BIGINT NOT NULL, -- Bytes
    checksum VARCHAR(128), -- SHA-256 hash
    storage_path TEXT, -- S3/GCS/local path
    status VARCHAR(50) DEFAULT 'pending',
        -- Valid: pending, processing, completed, failed
    uploaded_by VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT chk_upload_status CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

CREATE INDEX idx_uploads_env ON uploads(environment_id, created_at DESC);
CREATE INDEX idx_uploads_status ON uploads(status);

-- =============================================
-- TEMPLATES & CATALOG (C10)
-- =============================================

-- Environment templates - Reusable configurations
CREATE TABLE environment_templates (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50), -- e.g., load-test, integration, production-mirror
    capabilities JSONB NOT NULL DEFAULT '[]',
    enablers_config JSONB NOT NULL DEFAULT '{}',
    compute_config JSONB NOT NULL DEFAULT '{}',
    storage INTEGER DEFAULT 0,
    network VARCHAR(50) DEFAULT 'private',
    tags TEXT[],
    popularity INTEGER DEFAULT 0, -- Usage count
    estimated_cost DECIMAL(10, 2),
    is_public BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_category ON environment_templates(category);
CREATE INDEX idx_templates_popularity ON environment_templates(popularity DESC);
CREATE INDEX idx_templates_public ON environment_templates(is_public);

-- =============================================
-- ERROR HANDLING & RECOVERY (C13)
-- =============================================

-- Error records - Tracking failures and recovery
CREATE TABLE error_records (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) REFERENCES environments(id) ON DELETE CASCADE,
    error_type VARCHAR(50) NOT NULL,
        -- Valid: validation, provisioning, execution, network, timeout, resource_exhausted
    severity VARCHAR(20) NOT NULL,
        -- Valid: warning, recoverable, fatal
    error_code VARCHAR(50),
    message TEXT NOT NULL,
    stack_trace TEXT,
    context JSONB DEFAULT '{}',
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    retry_strategy VARCHAR(50), -- exponential_backoff, linear, fixed
    last_retry_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolution_action TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_error_severity CHECK (severity IN ('warning', 'recoverable', 'fatal'))
);

CREATE INDEX idx_error_records_env ON error_records(environment_id, created_at DESC);
CREATE INDEX idx_error_records_severity ON error_records(severity);
CREATE INDEX idx_error_records_unresolved ON error_records(resolved_at) WHERE resolved_at IS NULL;

-- Rollback operations - Track rollback history
CREATE TABLE rollback_operations (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    trigger_type VARCHAR(50) NOT NULL,
        -- Valid: manual, error, timeout, policy_violation
    rollback_target VARCHAR(50), -- State or snapshot to rollback to
    steps JSONB NOT NULL DEFAULT '[]', -- Array of rollback steps
    status VARCHAR(50) DEFAULT 'pending',
        -- Valid: pending, in_progress, completed, failed
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    
    CONSTRAINT chk_rollback_status CHECK (status IN ('pending', 'in_progress', 'completed', 'failed'))
);

CREATE INDEX idx_rollback_operations_env ON rollback_operations(environment_id, started_at DESC);

-- =============================================
-- SECURITY & COMPLIANCE (C09, C18)
-- =============================================

-- User roles and permissions
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
        -- Valid: admin, operator, developer, viewer
    team VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_user_role CHECK (role IN ('admin', 'operator', 'developer', 'viewer'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Access policies
CREATE TABLE access_policies (
    id SERIAL PRIMARY KEY,
    resource_type VARCHAR(50) NOT NULL, -- environment, template, etc.
    resource_id VARCHAR(255),
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    team VARCHAR(255),
    permissions JSONB NOT NULL DEFAULT '[]',
        -- Array of: read, write, delete, provision, execute
    granted_by VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_access_policies_resource ON access_policies(resource_type, resource_id);
CREATE INDEX idx_access_policies_user ON access_policies(user_id);

-- Compliance checks
CREATE TABLE compliance_checks (
    id SERIAL PRIMARY KEY,
    environment_id VARCHAR(50) REFERENCES environments(id) ON DELETE CASCADE,
    check_type VARCHAR(100) NOT NULL, -- e.g., encryption, retention, tagging
    status VARCHAR(50) NOT NULL,
        -- Valid: passed, failed, warning, skipped
    policy_name VARCHAR(255),
    details JSONB DEFAULT '{}',
    evidence TEXT,
    checked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_compliance_status CHECK (status IN ('passed', 'failed', 'warning', 'skipped'))
);

CREATE INDEX idx_compliance_checks_env ON compliance_checks(environment_id, checked_at DESC);
CREATE INDEX idx_compliance_checks_status ON compliance_checks(status);

-- =============================================
-- VIEWS FOR COMMON QUERIES
-- =============================================

-- Active environments with current costs
CREATE VIEW active_environments AS
SELECT 
    e.id,
    e.name,
    e.owner,
    e.status,
    e.health,
    e.uptime,
    e.estimated_cost,
    e.actual_cost,
    COALESCE(jsonb_array_length(e.capabilities), 0) as capability_count,
    e.created_at,
    e.updated_at
FROM environments e
WHERE e.status IN ('running', 'provisioning')
    AND e.deleted_at IS NULL
ORDER BY e.created_at DESC;

-- Environment summary with metrics
CREATE VIEW environment_summary AS
SELECT 
    e.id,
    e.name,
    e.status,
    e.owner,
    COUNT(DISTINCT ra.id) as resource_count,
    SUM(cr.amount) as total_cost,
    MAX(ms.timestamp) as last_metric_time,
    e.health,
    e.created_at
FROM environments e
LEFT JOIN resource_allocations ra ON e.id = ra.environment_id
LEFT JOIN cost_records cr ON e.id = cr.environment_id
LEFT JOIN metrics_snapshots ms ON e.id = ms.environment_id
WHERE e.deleted_at IS NULL
GROUP BY e.id;

-- Recent audit trail
CREATE VIEW recent_audit_trail AS
SELECT 
    al.id,
    al.environment_id,
    e.name as environment_name,
    al.action,
    al.user_id,
    al.result,
    al.created_at
FROM audit_logs al
LEFT JOIN environments e ON al.environment_id = e.id
ORDER BY al.created_at DESC
LIMIT 1000;

-- =============================================
-- FUNCTIONS AND TRIGGERS
-- =============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at trigger to relevant tables
CREATE TRIGGER update_environments_updated_at
    BEFORE UPDATE ON environments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_capabilities_updated_at
    BEFORE UPDATE ON capabilities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_enablers_updated_at
    BEFORE UPDATE ON enablers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically create audit log on environment changes
CREATE OR REPLACE FUNCTION log_environment_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO audit_logs (environment_id, action, entity_type, entity_id, user_id, details)
        VALUES (NEW.id, 'created', 'environment', NEW.id, NEW.owner, 
                jsonb_build_object('name', NEW.name, 'status', NEW.status));
    ELSIF TG_OP = 'UPDATE' THEN
        IF OLD.status != NEW.status THEN
            INSERT INTO audit_logs (environment_id, action, entity_type, entity_id, user_id, details)
            VALUES (NEW.id, 'status_changed', 'environment', NEW.id, NEW.owner,
                    jsonb_build_object('old_status', OLD.status, 'new_status', NEW.status));
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit_logs (environment_id, action, entity_type, entity_id, user_id, details)
        VALUES (OLD.id, 'deleted', 'environment', OLD.id, OLD.owner,
                jsonb_build_object('name', OLD.name));
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER environment_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON environments
    FOR EACH ROW
    EXECUTE FUNCTION log_environment_change();

-- =============================================
-- INITIAL DATA SEEDING
-- =============================================

-- Insert sample capabilities
INSERT INTO capabilities (id, name, description, enablers, dependencies) VALUES
('C01', 'Spec Authoring & Validation', 'Web editor and services to author SES', '["E02","E17","E11","E20"]', '[]'),
('C02', 'Parsing & Internal Modeling', 'Convert validated SES into normalized models', '["E03","E04","E17","E01"]', '["C01"]'),
('C03', 'Planning Engine', 'Generates execution plan', '["E03","E04","E08","E16"]', '["C01","C02"]'),
('C04', 'Provisioning Automation', 'Executes plan to create infrastructure', '["E05","E04","E01","E18","E20"]', '["C03","C09"]'),
('C05', 'Orchestration & State Machine', 'Manages lifecycle phases', '["E04","E03","E18","E01"]', '["C04","C17"]'),
('C06', 'Monitoring & Metrics', 'Collects infrastructure metrics', '["E06","E07","E11","E19"]', '["C04","C12"]'),
('C07', 'Logging & Audit', 'Structured logs with audit trail', '["E07","E19","E01","E09","E10"]', '["C16","C17"]'),
('C08', 'Cost Management', 'Real-time cost tracking', '["E08","E16","E06","E11"]', '["C06","C03"]'),
('C09', 'Security & Compliance', 'Credential vault, RBAC', '["E09","E10","E19","E01"]', '["C07","C17"]'),
('C10', 'Spec-Kit & Reuse', 'Library of templates', '["E17","E02","E11","E20"]', '["C17","C18"]'),
('C11', 'Reservation & Scheduling', 'Time-window validation', '["E14","E04","E03","E09"]', '["C02","C17"]'),
('C12', 'Simulation Execution', 'Runs scenarios', '["E12","E13","E04","E06","E20"]', '["C04","C05"]'),
('C13', 'Error Handling & Recovery', 'Automated rollback', '["E15","E07","E04","E18","E13"]', '["C05","C07"]'),
('C14', 'Provider Abstraction', 'Unified adapter layer', '["E05","E18","E01","E20"]', '["C09","C04"]'),
('C15', 'Environment Visualization', 'Topology dashboards', '["E11","E06","E07","E03","E08"]', '["C06","C08"]'),
('C16', 'Messaging & Agent Coordination', 'Internal message bus', '["E18","E01","E07","E09"]', '["C17","C09"]'),
('C17', 'State Persistence Layer', 'ACID store for state', '["E01","E19","E07","E09"]', '["C09"]'),
('C18', 'Access & Governance', 'Role management', '["E09","E10","E19","E07","E11"]', '["C17","C07"]')
ON CONFLICT (id) DO NOTHING;

-- Insert sample enablers
INSERT INTO enablers (id, name, description) VALUES
('E01', 'Core Platform Infra', 'DB, storage, vault, message broker'),
('E02', 'Schema & Validation', 'JSON/YAML validators, rule engine'),
('E03', 'Graph & Planning', 'DAG construction, topological sort'),
('E04', 'Execution Framework', 'Workflow engine, task dispatcher'),
('E05', 'Provider SDK', 'Cloud/on-prem SDK wrappers'),
('E06', 'Metrics Stack', 'Prometheus, time-series storage'),
('E07', 'Logging & Tracing', 'Central collector, distributed tracing'),
('E08', 'Cost Engine', 'Pricing cache, estimation algorithms'),
('E09', 'Security/RBAC', 'Auth, MFA, policy enforcement'),
('E10', 'Compliance Engine', 'Rule DSL, scheduled checks'),
('E11', 'UI Components', 'Web editor, dashboards, visualizers'),
('E12', 'Workload Drivers', 'Load/stress test integrations'),
('E13', 'Health & Validation', 'Endpoint pollers, status aggregators'),
('E14', 'Reservation Scheduler', 'Time-slot index, priority queue'),
('E15', 'Rollback Manager', 'LIFO unwinder, deallocator'),
('E16', 'Optimization Advisory', 'Right-sizing recommendations'),
('E17', 'Template & Catalog', 'Versioned spec templates'),
('E18', 'Message Protocol', 'Serialization, correlation IDs'),
('E19', 'Data Protection', 'Encryption, data classification'),
('E20', 'Testing & QA', 'Spec fixtures, provider mocks')
ON CONFLICT (id) DO NOTHING;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ses_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ses_user;
