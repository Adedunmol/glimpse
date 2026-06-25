#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Usage: ./upload_photos.sh <upload_id> <folder> <token>
# Example: ./upload_photos.sh abc123 ./photos eyJhbGci...
# ---------------------------------------------------------------------------

UPLOAD_ID="${1:-}"
FOLDER="${2:-}"
TOKEN="${3:-}"
BASE_URL="http://localhost:8080/api/v1"

# --- Validate args ----------------------------------------------------------
if [[ -z "$UPLOAD_ID" || -z "$FOLDER" || -z "$TOKEN" ]]; then
  echo "Usage: $0 <upload_id> <folder> <token>"
  exit 1
fi

if [[ ! -d "$FOLDER" ]]; then
  echo "Error: '$FOLDER' is not a directory."
  exit 1
fi

if ! command -v jq &>/dev/null; then
  echo "Error: 'jq' is required but not installed."
  echo "  macOS:  brew install jq"
  echo "  Ubuntu: sudo apt install jq"
  exit 1
fi

# --- Collect image files ----------------------------------------------------
mapfile -t FILES < <(find "$FOLDER" -maxdepth 1 -type f \
  \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \
     -o -iname "*.webp" -o -iname "*.heic" \) | sort)

if [[ ${#FILES[@]} -eq 0 ]]; then
  echo "Error: No image files found in '$FOLDER'."
  exit 1
fi

echo "Found ${#FILES[@]} image(s) in '$FOLDER'."

# --- Build the files JSON array for the init request -----------------------
# e.g. [{"name":"photo1.jpg"},{"name":"photo2.png"}]
FILES_JSON="["
for i in "${!FILES[@]}"; do
  NAME=$(basename "${FILES[$i]}")
  FILES_JSON+=$(printf '{"name":"%s"}' "$NAME")
  [[ $i -lt $(( ${#FILES[@]} - 1 )) ]] && FILES_JSON+=","
done
FILES_JSON+="]"

# --- Step 1: Request presigned upload URLs ----------------------------------
echo ""
echo "[1/3] Requesting presigned URLs..."

INIT_RESPONSE=$(curl --silent --fail \
  --request POST \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer $TOKEN" \
  --data "{\"files\": $FILES_JSON}" \
  "$BASE_URL/uploads/$UPLOAD_ID/photos")

echo "Init response received."

# --- Parse the response -----------------------------------------------------
# Extract keys and URLs as "key|url" lines using jq
mapfile -t KV_PAIRS < <(echo "$INIT_RESPONSE" | jq -r '.uploads[] | .key + "|" + .url')

if [[ ${#KV_PAIRS[@]} -ne ${#FILES[@]} ]]; then
  echo "Error: Server returned ${#KV_PAIRS[@]} URLs but we have ${#FILES[@]} files."
  exit 1
fi

# --- Step 2: Upload each file to S3 (parallel) ------------------------------
echo ""
echo "[2/3] Uploading ${#FILES[@]} file(s) to S3 in parallel..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

PIDS=()

for i in "${!FILES[@]}"; do
  FILE="${FILES[$i]}"
  NAME=$(basename "$FILE")
  KEY="${KV_PAIRS[$i]%%|*}"
  URL="${KV_PAIRS[$i]##*|}"
  MIME=$(file --mime-type -b "$FILE")

  (
    HTTP_STATUS=$(curl --silent --output /dev/null --write-out "%{http_code}" \
      --request PUT \
      --header "Content-Type: $MIME" \
      --upload-file "$FILE" \
      "$URL")

    if [[ "$HTTP_STATUS" -lt 200 || "$HTTP_STATUS" -ge 300 ]]; then
      echo "  ✗ $NAME (HTTP $HTTP_STATUS)"
      exit 1
    fi

    echo "  ✓ $NAME (HTTP $HTTP_STATUS)"
  ) &

  PIDS+=($!)
done

# Wait for all uploads and check for failures
FAILED=0
for i in "${!PIDS[@]}"; do
  if ! wait "${PIDS[$i]}"; then
    FAILED=1
  fi
done

if [[ "$FAILED" -ne 0 ]]; then
  echo ""
  echo "Error: One or more uploads failed. Aborting."
  exit 1
fi

# Build KEYS_JSON now that all uploads succeeded
KEYS_JSON="["
for i in "${!KV_PAIRS[@]}"; do
  KEY="${KV_PAIRS[$i]%%|*}"
  KEYS_JSON+=$(printf '{"key":"%s"}' "$KEY")
  [[ $i -lt $(( ${#KV_PAIRS[@]} - 1 )) ]] && KEYS_JSON+=","
done
KEYS_JSON+="]"

# --- Step 3: Mark upload complete to trigger processing --------------------
echo ""
echo "[3/3] Completing upload..."

COMPLETE_RESPONSE=$(curl --silent --fail \
  --request POST \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer $TOKEN" \
  --data "{\"files\": $KEYS_JSON}" \
  "$BASE_URL/uploads/$UPLOAD_ID/complete")

echo "Done. Server response:"
echo "$COMPLETE_RESPONSE" | jq . 2>/dev/null || echo "$COMPLETE_RESPONSE"

echo ""
echo "All ${#FILES[@]} photo(s) uploaded and processing started."