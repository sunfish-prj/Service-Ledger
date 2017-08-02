'use strict';

var db_utils = require('../utils/dbUtils.js');

/**
 * Storing a key-value pair 
 *
 * putSpec Put-request-body Body in JSON
 * returns response
 **/
exports.putPOST = function(putSpec) {

    var result = {};

	var ack = db_utils.db_put(putSpec);
	
	console.log("res: " + ack._id);
		
	result['application/json'] = { "message" : ack };

    return new Promise(function(resolve, reject) {
      if (Object.keys(result).length > 0) {
        resolve(result[Object.keys(result)[0]]);
      } else {
        resolve();
      }
    });
}
