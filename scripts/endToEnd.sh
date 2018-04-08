#!/bin/bash

DIR="$(dirname "$0")"
source $DIR/helperFunctions.sh

if [ "$#" -gt 1 ]; then
  SLOWMODE="Y"
fi

echo "===================== SETUP BEGIN ====================="
createChannel 1 0

joinChannel 1 0
joinChannel 1 1
sleep 1
#joinAnchors 1

joinChannel 2 0
joinChannel 2 1
sleep 1
#joinAnchors 2

installChaincode 1 0
installChaincode 1 1
installChaincode 2 0
installChaincode 2 1

instantiateChaincode 1 0
sleep 5
echo "===================== SETUP END ====================="

echo "===================== TEST RUN BEGIN ===================== "
invoke 1 0 "[\"reg\",\"Test\"]"
sleep 5
query 1 0 "Test"
query 1 1 "Test"
query 2 0 "Test"
query 2 1 "Test"
echo "===================== TEST RUN END ===================== "

echo "===================== END TO END FINISHED!!! ===================== "