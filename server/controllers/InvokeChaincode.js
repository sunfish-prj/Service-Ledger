'use strict'

var rp = require('request-promise');
var config = require('config');
var url = require('url');
var hyperledger_paras = config.get('hyperledger-service');

var debug = true;

module.exports.invokeChaincode = function invokeChaincode(req, res, next) {

  if(debug) console.log('---->ServiceLedger: invokeChaincode method called');

  var examples = {};
  examples['application/json'] = {
    "message": "example message of invokeChaincode"
  }


  rp({
      method: "POST",
      uri: url.format({
           protocol: 'http',
           hostname: hyperledger_paras.host,
           port: hyperledger_paras.port,
           pathname: hyperledger_paras.path.invoke + '/' + req.body.chaincodeName
      }),
      body: {
           "fcn": req.body.fcn,
           "args": req.body.args,
      },
      headers: {
           "authorization": "Bearer " + req.body.authorization,
           "content-type": "application/json"
      },
      json: true
  }).then(response => {
     if(debug) console.log(response); 
     examples['application/json'].message = response;
  
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
