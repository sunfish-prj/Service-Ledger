'use strict'

var debug = true;

module.exports.installChaincode = function installChaincode(req, res, next) {

  if(debug) console.log('---->ServiceLedger: installChaincode method called');

  var examples = {};
  examples['application/json'] = {
    "message": "example message of installChaincode"
  }
  
  if (Object.keys(examples).length > 0) {
    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
  } else {
    res.end();
  }
}
