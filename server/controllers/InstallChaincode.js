'use strict'

var rp = require('request-promise');
var config = require('config');
var url = require('url');
var hyperledger_paras = config.get('hyperledger-service');

var debug = true;

module.exports.installChaincode = function installChaincode(req, res, next) {

  if(debug) console.log('---->ServiceLedger: installChaincode method called');

  var examples = {};
  examples['application/json'] = {
    "message": "example message of installChaincode"
  }


  rp({
      method: "POST",
      uri: url.format({
           protocol: 'http',
           hostname: hyperledger_paras.host,
           port: hyperledger_paras.port,
           pathname: hyperledger_paras.path.install
      }),
      body: {
           "peers": req.body.peers,
           "chaincodeName": req.body.chaincodeName,
           "chaincodePath": req.body.chaincodePath,
           "chaincodeType": req.body.chaincodeType,
           "chaincodeVersion": req.body.chaincodeVersion
      },
      header: {
           "authorization": req.body.authorization,
           "content-type": "application/json"
      },
      json: true
  }).then(response => {
     if(debug) console.log(response); 
  
     if (Object.keys(examples).length > 0) {
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
     } else {
        res.end();
     }
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
