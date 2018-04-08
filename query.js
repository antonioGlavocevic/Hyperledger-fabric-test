'use strict';

let fs = require('fs');
let path = require('path');

let helper = require('./helper.js');

query();

async function query() {
	var fabric_client = await helper.getClient('org1','user1');
	var channel = fabric_client.getChannel('mychannel');

	var request = {
		targets: ['peer0.org1.example.com','peer1.org1.example.com','peer0.org2.example.com','peer1.org2.example.com'],
		chaincodeId: 'mycc',
		fcn: 'query',
		args: ['TOTAL']
	};
	
	let query_responses = await channel.queryByChaincode(request);

	console.log("Query has completed, checking results");
	// query_responses could have more than one  results if there multiple peers were used as targets
	if (query_responses && query_responses.length >= 1) {
		query_responses.forEach((query_response) => {
			if (query_response instanceof Error) {
				console.error("error from query = ", query_response);
			} else {
				console.log("Response is ", query_response.toString());
			}
		})
	} else {
		console.log("No payloads were returned from query");
	}
}
