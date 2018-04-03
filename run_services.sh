#!/bin/bash

echo "Running pubmedr #####################################################################"
pubmedr

echo "Running fixr #########################################################"
fixr -b 1 -t "pubmedData,fixResources"

echo "Running mongr #####################################################################"
mongr -b 1 -c all

echo "Running algr #####################################################################"
algr -c all

echo "All done!"
