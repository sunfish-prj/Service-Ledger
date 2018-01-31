'use strict';

var db_utils = require('../../hyperledger/fabric/utils/dbUtils.js');

/**
 * Delete a stored key 
 *
 * body Delete-request-body 
 * returns response
 **/
exports.deletePOST = function(body) {
    var message = {};	
    return new Promise(function(resolve, reject) {
 	  db_utils.db_delete(body, function(res){	
		if(Object.keys(res).length > 0) {	
	  		message = JSON.stringify({"message" : res});
	  		console.log(message);
			  resolve(message);
		  }else{
			  reject(message);
		  }
	  });	  
    });
}

