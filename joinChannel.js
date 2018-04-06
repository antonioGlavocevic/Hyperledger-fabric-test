'use strict';

let fs = require('fs');
let path = require('path');

let helper = require('./helper.js');

joinChannel();

async function joinChannel() {
	var fabric_client = await helper.getClient('user1');
	var channel = fabric_client.getChannel('mychannel');

	var tx_id = fabric_client.newTransactionID(true);
	let g_request = {
		txId: tx_id
	};
	let genesis_block = await channel.getGenesisBlock(g_request);

	tx_id = fabric_client.newTransactionID(true);
	let j_request = {
		targets: ['peer0.org1.example.com','peer1.org1.example.com'],
		block: genesis_block,
		txId: tx_id
	};
	let results = await channel.joinChannel(j_request);
	if(results && results.response && results.response.status == 200) {
		console.log('success');
	} else {
		console.log('fail');
	}
}
