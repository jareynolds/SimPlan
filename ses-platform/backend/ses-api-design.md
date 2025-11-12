POST   /api/v1/environments              # Create new environment spec
GET    /api/v1/environments              # List all environments
GET    /api/v1/environments/{id}         # Get environment details
PUT    /api/v1/environments/{id}         # Update environment
DELETE /api/v1/environments/{id}         # Delete environment
POST   /api/v1/environments/{id}/provision    # Trigger provisioning
POST   /api/v1/environments/{id}/upload       # Upload binaries/configs
GET    /api/v1/environments/{id}/status       # Get current status
GET    /api/v1/environments/{id}/metrics      # Get metrics data
GET    /api/v1/environments/{id}/logs         # Get logs

GET    /api/v1/capabilities              # List all capabilities (C01-C18)
GET    /api/v1/enablers                  # List all enablers (E01-E20)
GET    /api/v1/templates                 # List spec templates
POST   /api/v1/validate                  # Validate spec (C01)
GET    /api/v1/cost/estimate             # Get cost estimation (C08)

