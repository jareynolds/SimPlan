# C17 State Persistence Layer
Reference: Capability C17
Description: ACID store for environment state, history, resources, metrics, logs, spec hash indexing and PITR.
Linked Enablers: E01 Core Platform Infra; E19 Data Protection Layer; E07 Logging & Tracing Stack; E09 Security/RBAC Components.
Rationale: Guarantees durability and auditability for transitions (C05) and recovery (C13).
Dependencies: Security (C09), data classification (E19).
