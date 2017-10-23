'use strict'

var rp = require('request-promise');
var config = require('config');
var url = require('url');
var hyperledger_paras = config.get('hyperledger-service');

var debug = true;

module.exports.setupLedger= function setupLedger(req, res, next) {

  if(debug) console.log('---->ServiceLedger: setupLedger method called');

  var examples = {};
  examples['application/json'] = {
    "message": "example message of setupLedger"
  }

  // login
  rp({
      method: "POST",
      uri: url.format({
           protocol: 'http',
           hostname: hyperledger_paras.host,
           port: hyperledger_paras.port,
           pathname: hyperledger_paras.path.login
      }),
      body: {
           "username": req.body.username,
           "orgName": req.body.orgName,
      },
      headers: {
           "content-type": "application/x-www-form-urlencoded"
      },
      json: true
  }).then(response => {
     if(debug) console.log(response); 
     if(response.success == true) {
        examples['application/json'].message = response.token;
     }
     if (Object.keys(examples).length > 0) {
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
     } else {
        res.end();
     }

     // create channel
     rp({
         method: "POST",
         uri: url.format({
              protocol: 'http',
              hostname: hyperledger_paras.host,
              port: hyperledger_paras.port,
              pathname: hyperledger_paras.path.createChannel
              }),
         body: {
              "channelName": req.body.channelName,
              "channelConfigPath": req.body.channelConfigPath,
              },
         headers: {
              "authorization": "Bearer " + response.token,
              "content-type": "application/json"
              },
         json: true
     });

     // join channel
     rp({
         method: "POST",
         uri: url.format({
              protocol: 'http',
              hostname: hyperledger_paras.host,
              port: hyperledger_paras.port,
              pathname: hyperledger_paras.path.joinChannel
              }),
         body: {
              "peers": req.body.peers,
              },
         headers: {
              "authorization": "Bearer " + response.token,
              "content-type": "application/json"
              },
         json: true
     });


  }).catch(err => {
     if(debug) console.log("----->ServiceLedger: error when requesting hyperledger!"); 
     examples['application/json'].message = "error when requesting hyperledger";
  
     if (Object.keys(examples).length > 0) {
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
     } else {
        res.end();
     }
  }); 

}
