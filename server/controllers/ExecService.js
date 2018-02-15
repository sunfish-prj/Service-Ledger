'use strict';

// output service configuration
var config = require('config');
var out_service_name = config.get('out-service.name');
// var hl_invoke_type = config.get('out-service.hl_invoke_type');

var db_utils = require('../utils/dbUtils.js');
var hl_utils = require('../utils/hlUtils.js');

exports.invokePOST = function(args, res, next) {
  /**
   * invoke the functions in the chaincode
   *
   * invokeSpec Invoke-chaincode-body
   * returns response
   **/
//   var examples = {};
//   examples['application/json'] = {
//   "message" : "aeiou"
// };
//   if (Object.keys(examples).length > 0) {
//     res.setHeader('Content-Type', 'application/json');
//     res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
//   } else {
//     res.end();
//   }

  if (out_service_name == 'fabric') {

    console.log("[ExecService.js] Hyperledger Fabric - INVOKE");

    var payload = args.invokeSpec.value;

    console.log(payload);

    hl_utils.hl_invoke(payload, function (result) {
      if (Object.keys(result).length > 0) {
        res.writeHead(200,{'Content-Type':'application/json'});
        response['application/json'] = {
          "message" : result
         };
        res.end(JSON.stringify(response[Object.keys(response)[0]] || {}, null, 2));
      } else {
          res.writeHead(400,{'Content-Type':'application/json'});
    		  res.end(JSON.stringify({'message': 'error'}));
      }
    });
  }

}
