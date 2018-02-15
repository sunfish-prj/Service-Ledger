'use strict';

var config = require('config');
var out_service_name = config.get('out-service.name'); 

var db_utils = require('../utils/dbUtils.js');
var hl_utils = require('../utils/hlUtils.js');

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
			
				console.log("[GetService.js] Calling Mongo - GET");
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
				console.log("[GetService.js] Calling Hyperledger Fabric - GET");
				hl_utils.hl_get(getId, function (res) {
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

