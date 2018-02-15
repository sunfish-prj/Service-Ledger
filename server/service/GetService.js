'use strict';

var config = require('config');
var out_service_name = config.get('out-service.name'); 

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
			
			if (out_service_name == 'mongo') {
			
				console.log("Calling mongo - put");

				
		 	  	db_utils.db_get(getId, function(res){	
				//db_utils.db_getKeys(getId, function(res){		   
				if (Object.keys(res).length > 0) {	
			  		message = JSON.stringify({"message" : res});
			  		console.log(message);
			   		resolve(message);
				  }else if (Object.keys(res).length == 0){
					  message = JSON.stringify({"message" : ''});
						console.log(message);
						resolve(message);
					}else{
					  console.log(JSON.stringify({"message" : res}));
		   			reject(message);
				  }
			  });
			}
		

			if (out_service_name == 'fabric') {
			
				console.log("Calling hyperledger");
			
				console.log("Calling api to 'put' in the keyValueStore chaincode...");
				hl_utils.hl_put(putSpec, function (res) {
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

