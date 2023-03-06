#!/bin/bash

kubectl create --save-config cm settings-cas --from-file config/config.yaml -o yaml -n juanji --dry-run=client  | kubectl apply -f -