'use strict';

// output service configuration
var config = require('config');
var out_service_name = config.get('out-service.name');
var hl_invoke_type = config.get('out-service.hl_invoke_type');

var db_utils = require('../utils/dbUtils.js');
var hl_utils = require('../utils/hlUtils.js');


/**
 * invoke a chaincode function 
 *
 * invokeSpec: Invoke-request-body Body in JSON
 * returns response
 **/
exports.InvokePOST = function (invokeSpec) {
  var message = {};
  return new Promise(function (resolve, reject) {

      if (out_service_name == 'fabric') {
        console.log("Calling hyperledger");
        if (hl_invoke_type != 'invoke') {
            console.log("Bad definition of parameter hl_invoke_type. Setting default to 'invoke'");
        }
        console.log("Calling api to 'invoke' a chaincode...");
        hl_utils.hl_invoke(invokeSpec, function (res) {
          if (Object.keys(res).length > 0) {
              message = JSON.stringify({ "message": res });
              console.log(message);
              resolve(message);
          } else {
              reject(message);
          }
        });
      }
  });
}
