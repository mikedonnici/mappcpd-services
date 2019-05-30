#!/bin/bash

echo "Running pubmedr #####################################################################"
pubmedr

echo "Running fixr #########################################################"
fixr -b 1 -t "pubmedData,fixResources"

echo "Running syncr #####################################################################"
syncr -b 1 -c all

echo "Running fixr #########################################################"
fixr -b 1 -t "pubmedData,fixResources"

echo "Running algr #####################################################################"
algr -c all

echo "Running backupdb ##################################################################"
backupdb

echo "All done!"
