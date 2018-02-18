#!/bin/bash

echo "Running pubmedr #####################################################################"
pubmedr

echo "Running fixr -t pubmedData #########################################################"
fixr -b 1 -t "pubmedData"

echo "Running mongr #####################################################################"
mongr -b 1 -c all

echo "Running fixr -t fixResources #####################################################"
fixr -b 1 -t "fixResources"

echo "Running algr #####################################################################"
algr -b 1

echo "All done!"
