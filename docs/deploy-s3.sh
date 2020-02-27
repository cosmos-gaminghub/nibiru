#!/bin/sh

aws s3 rm s3://$DOCS_BUCKET/ --recursive
aws s3 cp ./.vuepress/dist s3://$DOCS_BUCKET/ --recursive
