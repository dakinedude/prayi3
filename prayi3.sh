#!/bin/bash

PRAYER_TIMES_PATH="/home/mats/go/bin/prayi3"

# Run the Go program and output prayer information
PRAYER_INFO=$($PRAYER_TIMES_PATH)
echo "$PRAYER_INFO"
