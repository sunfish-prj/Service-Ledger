'use strict';

exports.invokePOST = function(args, res, next) {
  /**
   * invoke the functions in the chaincode
   *
   * invokeSpec Invoke-chaincode-body 
   * returns response
   **/
  var examples = {};
  examples['application/json'] = {
  "message" : "aeiou"
};
  if (Object.keys(examples).length > 0) {
    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
  } else {
    res.end();
  }
}

