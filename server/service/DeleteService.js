'use strict';

var config = require('config');
var out_service_name = config.get('out-service.name'); 

var db_utils = require('../utils/dbUtils.js');
var hl_utils = require('../utils/hlUtils.js');

/**
 * Delete a stored key 
 *
 * body Delete-request-body 
 * returns response
 **/
exports.deletePOST = function(body) {
    var message = {};	
    
	return new Promise(function(resolve, reject) {
		
		if (out_service_name == 'mongo') {
			
			console.log("[DeleteService.js] Calling Mongo - DELETE");
			db_utils.db_delete(body, function(res){	
			if (Object.keys(res).length > 0) {	
				message = JSON.stringify({"message" : res});
				console.log(message);
				resolve(message);
			} else {
				reject(message);
			}
			});	
		}

		if (out_service_name == 'fabric') {
		
			console.log("[DeleteService.js] Calling Hyperledger Fabric - DELETE");
			hl_utils.hl_delete(body, function (res) {
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

