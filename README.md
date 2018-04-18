## Generate crypo

cryptogen generate --config=./crypto-config.yaml

export FABRIC_CFG_PATH=$PWD
mkdir channel-artifacts
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

export CHANNEL_NAME=mychannel
configtxgen -profile OneOrgChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

## Add yourself to network NOT UP TO DATE

docker-compose -f docker-compose5.yml up -d

docker exec -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer1.org1.example.com peer channel fetch config \
  -o ec2-13-59-126-206.us-east-2.compute.amazonaws.com:7050 -c mychannel --tls \
  --cafile /etc/hyperledger/crypto/orderer/ec2-13-59-126-206.us-east-2.compute.amazonaws.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker exec -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer1.org1.example.com peer channel join \
  -b mychannel_config.block --tls \
  --cafile /etc/hyperledger/crypto/orderer/ec2-13-59-126-206.us-east-2.compute.amazonaws.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker logs -f peer1.org1.example.com
