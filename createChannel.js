'use strict';

let fs = require('fs');
let path = require('path');

let helper = require('./helper.js');

createChannel();

async function createChannel() {
	var fabric_client = await helper.getClient('org1');
	console.log('got client');
	let envelope = fs.readFileSync(path.join(__dirname,'config/channel.tx'));
	let channelConfig = fabric_client.extractChannelConfig(envelope);
	let signature = fabric_client.signChannelConfig(channelConfig);
	let tx_id = fabric_client.newTransactionID(true);

	let request = {
		config: channelConfig,
		signatures: [signature],
		name: 'mychannel',
		txId: tx_id
	};

	var response = await fabric_client.createChannel(request);
	console.log('response' + response);
	if (response && response.status === 'SUCCESS') {
		console.log('Successfully created the channel.');
	} else {
		throw new Error('Failed to create the channel');
	}
}
