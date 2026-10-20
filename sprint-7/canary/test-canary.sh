#!/usr/bin/env bash

URL="http://127.0.0.1:8088/"
TOTAL=100
V1=0
V2=0
ERRORS=0

echo "Testing canary traffic distribution..."
echo "Sending $TOTAL requests to $URL"
echo

for ((i=1; i<=TOTAL; i++)); do
    RESPONSE=$(curl -sS --connect-timeout 3 --max-time 5 \
        -H "Host: canary.local" "$URL" 2>/dev/null) || {
        ERRORS=$((ERRORS + 1))
        continue
    }

    case "$RESPONSE" in
        *"Version 1.0 - Stable"*)
            V1=$((V1 + 1))
            ;;
        *"Version 2.0 - Canary"*)
            V2=$((V2 + 1))
            ;;
        *)
            ERRORS=$((ERRORS + 1))
            ;;
    esac
done

SUCCESS=$((V1 + V2))

echo "========== CANARY TEST RESULTS =========="
echo "Total requests: $TOTAL"
echo "Stable v1 responses: $V1"
echo "Canary v2 responses: $V2"
echo "Errors or unknown responses: $ERRORS"

if (( SUCCESS > 0 )); then
    awk -v v1="$V1" -v v2="$V2" -v total="$SUCCESS" \
        'BEGIN {
            printf "Stable v1: %.1f%%\n", v1/total*100
            printf "Canary v2: %.1f%%\n", v2/total*100
        }'
fi

echo "========================================="  