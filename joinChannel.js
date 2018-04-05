'use strict';

//var log4js = require('log4js');
//var logger = log4js.getLogger('Helper');
//logger.level = 'debug';

let Fabric_Client = require('fabric-client');
//Fabric_Client.setLogger(logger);
let fs = require('fs');
let path = require('path');

var fabric_client = Fabric_Client.loadFromConfig('network-config.yaml');
fabric_client.loadFromConfig('org1.yaml');

var channel = fabric_client.newChannel('mychannel');
var orderer = fabric_client.newOrderer('grpc://localhost:7050');
channel.addOrderer(orderer);
var peer = fabric_client.newPeer('grpc://localhost:7051');
channel.addPeer(peer);

fabric_client.initCredentialStores()
.then((nothing) => {
	fabric_client.setUserContext({username:'admin', password:'adminpw'})
	.then((admin) => {

		var tx_id = fabric_client.newTransactionID();
		let g_request = {
			txId: tx_id
		};
		channel.getGenesisBlock(g_request)
		.then((block) => {
			var genesis_block = block;
			var tx_id = fabric_client.newTransactionID();
			let j_request = {
				//targets: ['peer0.org1.example.com'],
				block: genesis_block,
				txId: tx_id
			};
			console.log(genesis_block);
			return channel.joinChannel(j_request);
		}).then((results) => {
			if(results && results.response && results.response.status == 200) {
				console.log('success');
			} else {
				console.log('fail');
			}
		});
	});
});
