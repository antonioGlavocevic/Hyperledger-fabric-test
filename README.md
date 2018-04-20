## Generate crypo

# Certs and keys
cryptogen generate --config=./crypto-config.yaml

# Genesis block
export FABRIC_CFG_PATH=$PWD
mkdir channel-artifacts
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

# Channel
export CHANNEL_NAME=mychannel
configtxgen -profile OneOrgChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

# Anchor peers
export CHANNEL_NAME=mychannel
configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP

## Add yourself to network

docker-compose -f docker-compose4.yml up -d

docker exec -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.hyperfabric.xyz/msp" peer1.org1.hyperfabric.xyz peer channel fetch config \
  -o orderer.hyperfabric.xyz:7050 -c mychannel --tls \
  --cafile /etc/hyperledger/crypto/orderer/orderer.hyperfabric.xyz/msp/tlscacerts/tlsca.hyperfabric.xyz-cert.pem

docker exec -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.hyperfabric.xyz/msp" peer1.org1.hyperfabric.xyz peer channel join \
  -b mychannel_config.block --tls \
  --cafile /etc/hyperledger/crypto/orderer/orderer.hyperfabric.xyz/msp/tlscacerts/tlsca.hyperfabric.xyz-cert.pem

docker logs -f peer1.org1.hyperfabric.xyz
