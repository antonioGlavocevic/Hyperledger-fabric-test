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

fabric_client.initCredentialStores()
.then((nothing) => {
	fabric_client.setUserContext({username:'admin', password:'adminpw'})
	.then((admin) => {
		let envelope = fs.readFileSync(path.join(__dirname,'config/channel.tx'));
		var channelConfig = fabric_client.extractChannelConfig(envelope);

		var signature = fabric_client.signChannelConfig(channelConfig);
		var orderer = fabric_client.newOrderer('grpc://localhost:7050')
		let tx_id = fabric_client.newTransactionID();

		let request = {
			config: channelConfig,
			signatures: [signature],
			name: 'mychannel',
			orderer: orderer,
			txId: tx_id
		};
		fabric_client.createChannel(request);
	})
})
