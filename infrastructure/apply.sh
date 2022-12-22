#! /bin/sh
aws cloudformation create-stack --stack-name $1 --template-body file://$2 --region=us-east-1 --parameter file://$3 --capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND