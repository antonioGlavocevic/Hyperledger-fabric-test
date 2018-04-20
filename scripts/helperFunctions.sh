#!/bin/bash

pauseCheck() {
  if [ ! -z "$SLOWMODE" ]; then
    read -n 1 -p "ANY KEY TO CONTINUE"
  fi
}

setEnv() {
  ORG=$1
  PEER=$2

  CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org${ORG}.hyperfabric.xyz/users/Admin@org${ORG}.hyperfabric.xyz/msp
  CORE_PEER_ADDRESS=peer${PEER}.org${ORG}.hyperfabric.xyz:7051
  CORE_PEER_LOCALMSPID=Org${ORG}MSP
  CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org${ORG}.hyperfabric.xyz/peers/peer${PEER}.org${ORG}.hyperfabric.xyz/tls/ca.crt

  echo "#######"
  echo "Env: Org${ORG} peer${PEER}"
  echo "#######"
}

createChannel() {
  setEnv $1 $2
  peer channel create -o orderer.hyperfabric.xyz:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/hyperfabric.xyz/orderers/orderer.hyperfabric.xyz/msp/tlscacerts/tlsca.hyperfabric.xyz-cert.pem
  echo "===================== \"$CHANNEL_NAME\" channel created ===================== "
  pauseCheck
}

joinChannel() {
  setEnv $1 $2
  peer channel join -b mychannel.block
  echo "===================== \"$CORE_PEER_ADDRESS\" joined \"$CHANNEL_NAME\" channel ===================== "
  pauseCheck
}

joinAnchors() {
  setEnv $1 $2
  peer channel update -o orderer.hyperfabric.xyz:7050 -c $CHANNEL_NAME -f ./channel-artifacts/Org${1}MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/hyperfabric.xyz/orderers/orderer.hyperfabric.xyz/msp/tlscacerts/tlsca.hyperfabric.xyz-cert.pem
  echo "===================== \"$CORE_PEER_LOCALMSPID\" anchored to \"$CHANNEL_NAME\" channel ===================== "
  pauseCheck
}

installChaincode() {
  setEnv $1 $2
  peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/
  echo "===================== Chaincode installed to \"$CORE_PEER_ADDRESS\" ===================== "
  pauseCheck
}

instantiateChaincode() {
  setEnv $1 $2
  peer chaincode instantiate -o orderer.hyperfabric.xyz:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/hyperfabric.xyz/orderers/orderer.hyperfabric.xyz/msp/tlscacerts/tlsca.hyperfabric.xyz-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.member', 'Org2MSP.member')"
  echo "===================== Chaincode instantiated ===================== "
  pauseCheck
}

invoke() {
  setEnv $1 $2
  peer chaincode invoke -o orderer.hyperfabric.xyz:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/hyperfabric.xyz/orderers/orderer.hyperfabric.xyz/msp/tlscacerts/tlsca.hyperfabric.xyz-cert.pem  -C $CHANNEL_NAME -n mycc -c '{"Args":'$3'}'
  pauseCheck
}

query() {
  setEnv $1 $2
  peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","'$3'"]}'
  pauseCheck
}
