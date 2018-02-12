'use strict';

// output service configuration
var config = require('config');
// var out_service_name = config.get('out-service.name');
// var hl_invoke_type = config.get('out-service.hl_invoke_type');

var db_utils = require('../utils/dbUtils.js');
var hl_utils = require('../utils/hlUtils.js');

var debug = true;

/**
 * invoke a chaincode function 
 *
 * invokeSpec: Invoke-request-body Body in JSON
 * returns response
 **/
exports.InvokePOST = function (invokeSpec) {
  //var message = {};

  if(debug) console.log("----> InvokePOST function in InvokeService.js called");

  return new Promise(function (resolve, reject) {

        // if (hl_invoke_type != 'invoke') {
        //    console.log("Bad definition of parameter hl_invoke_type. Setting default to 'invoke'");
        // }
        if(debug) console.log("Calling hyperledger api to 'invoke' a chaincode...");
        hl_utils.hl_invoke(invokeSpec, function (res) {
          if (Object.keys(res).length > 0) {
              var message = JSON.stringify({ "message": res});
              console.log(message);
              resolve(message);
          } else {
              reject(message);
          }
        });
  });
}
