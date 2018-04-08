'use strict';

let fs = require('fs');
let path = require('path');

let helper = require('./helper.js');

invoke();

async function invoke() {
	var fabric_client = await helper.getClient('org1','user1');
	var channel = fabric_client.getChannel('mychannel');

	var tx_id = fabric_client.newTransactionID();
	var tx_id_string = tx_id.getTransactionID();
	var request = {
		targets: ['peer0.org1.example.com', 'peer1.org2.example.com'],
		chaincodeId: 'mycc',
		fcn: 'reg',
		args: ['Test1'],
		chainId: 'mychannel',
		txId: tx_id
	};

	let results = await channel.sendTransactionProposal(request);
	var proposalResponses = results[0];
	var proposal = results[1];

	var all_good = true;
	for (var i in proposalResponses) {
		let one_good = false;
		if (proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200) {
			one_good = true;
			console.log('Transaction proposal was good');
		} else {
			console.error('Transaction proposal was bad');
		}
		all_good = all_good & one_good;
	}

	/*var promises = [];
	let event_hub = channel.getChannelEventHubsForOrg();
	event_hub.forEach((eh) => {
		let invokeEventPromise = new Promise((resolve, reject) => {
			let event_timeout = setTimeout(() => {
				let message = 'REQUEST_TIMEOUT:' + eh.getPeerAddr();
				console.log('ERROR:' + message);
				eh.disconnect();
			}, 3000);
			eh.registerTxEvent(tx_id_string, (tx, code, block_num) => {
				console.log('The chaincode invoke chaincode transaction has been committed on peer' + eh.getPeerAddr());
				console.log('Transaction '+tx+' has status of '+code+' in block '+block_num);
				clearTimeout(event_timeout);

				if (code !== 'VALID') {
					let message = util.format('The invoke chaincode transaction was invalid, code:%s',code);
					console.log(message);
					reject(new Error(message));
				} else {
					let message = 'The invoke chaincode transaction was valid.';
					console.log(message);
					resolve(message);
				}
			}, (err) => {
				clearTimeout(event_timeout);
				console.log(err);
				reject(err);
			},
				// the default for 'unregister' is true for transaction listeners
				// so no real need to set here, however for 'disconnect'
				// the default is false as most event hubs are long running
				// in this use case we are using it only once
				{unregister: true, disconnect: true}
			);
			eh.disconnect();
		});
		promises.push(invokeEventPromise);
	});*/

	var orderer_request = {
		txId: tx_id,
		proposalResponses: proposalResponses,
		proposal: proposal
	};

	channel.sendTransaction(orderer_request);
}
