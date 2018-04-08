'use strict';

let Fabric_Client = require('fabric-client');
let copService = require('fabric-ca-client');
let fs = require('fs');
let path = require('path');

async function getClient(organisation, username) {
    var fabric_client = Fabric_Client.loadFromConfig('network-config.yaml');
    fabric_client.loadFromConfig(organisation+'.yaml');
    await fabric_client.initCredentialStores();

    if (username) {
        let user = await fabric_client.getUserContext(username,true);
        if(!user) {
            throw new Error(username + ' was not found');
        } else {
            console.log(username + 'was found to be registered and enrolled');
        }
    }
    console.log('getClient done');
    return fabric_client;
}

var getRegisterUser = async function(organisation, username) {
    var fabric_client = await getClient(organisation);

    var admin = await fabric_client.setUserContext({username: 'admin', password: 'adminpw'});
    var caClient = fabric_client.getCertificateAuthority();
    let secret = await caClient.register({
        enrollmentID: username,
        affiliation: organisation+'.department1'
    }, admin);
    var user  = await fabric_client.setUserContext({username:username, password:secret});
    if(user && user.isEnrolled) {
        console.log('secret: \n'+user._enrollmentSecret+'\n');
    } else {
        throw new Error('User was not enrolled ');
    }
}

exports.getClient = getClient;
exports.getRegisterUser = getRegisterUser;