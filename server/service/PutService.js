'use strict';

var db_utils = require('../utils/dbUtils.js');

/**
 * Storing a key-value pair 
 *
 * putSpec Put-request-body Body in JSON
 * returns response
 **/
exports.putPOST = function(putSpec) {
    var message = {};	
    return new Promise(function(resolve, reject) {
 	  db_utils.db_put(putSpec, function(res){	
		  if (Object.keys(res).length > 0) {	
	  		message = JSON.stringify({"message" : res});
	  		console.log(message);
			resolve(message);
		  }else{
			reject(message);
		  }
	  });	  
    });
}
