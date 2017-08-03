'use strict';

var db_utils = require('../utils/dbUtils.js');

/**
 * Retrieving a value by its key 
 *
 * getId Get-request-body Body in JSON
 * returns response
 **/
exports.getPOST = function(getId) {
    var message = {};	
    return new Promise(function(resolve, reject) {
 	  db_utils.db_get(getId, function(res){	
		  if (Object.keys(res).length > 0) {	
	  		message = JSON.stringify({"message" : res});
	  		console.log(message);
			resolve(message);
		  }else{
			  console.log(message);
			reject(message);
		  }
	  });	  
    });
}

