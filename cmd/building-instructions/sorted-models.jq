#!/usr/local/bin/jq -rf

# Usage:
#   ./list-models.jq data/index.json | ./sorted-names.jq

[
  .[]
  | .model
]
| unique
| .[]
