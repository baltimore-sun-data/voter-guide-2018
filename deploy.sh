#! /bin/sh
set -eux -o pipefail

echo "Deploying site to S3"
s3deploy -config s3deploy.yaml -source public -bucket elections2018-news-baltimoresun-com -region us-east-1 -v

echo "Clearing CloudFront cache"
boreas -dist E2AXRGC2L5V6PH
