
// output service configuration
var config = require('config');
var assert = require('assert');

var exec = require('ssh-exec');
var hl_user = config.get('out-service.hl_user');
var hl_pass = config.get('out-service.hl_password');
var hl_ip = config.get('out-service.hl_ip');
var hl_channel = config.get('out-service.hl_default_channel');
var hl_chaincode = config.get('out-service.hl_default_chaincode');
var hl_dockerid = config.get('out-service.hl_default_dockerid');
var hl_script_path = config.get('out-service.hl_default_script_path');
var hl_peer = config.get('out-service.hl_default_peer');

/* Hyperledger Fabric - PUT */
var hl_put = exports.hl_put =  function(myobj, callback) {	
	console.log('[hlUtils.js] Calling PUT toward Hyperledger Fabric...');    
	_put(myobj, function(res){
	      return callback(res);
	  });
}

function _put (myobj, callback){
	console.log('[hlUtils.js] Executing PUT with following key and value:');

	
	var key = myobj.key;
	var value = myobj.value;
	console.log("[hlUtils.js] PUT key:"+key+" value:"+value);
	var command = hl_script_path + 'hl_put.sh ' +' '+ hl_peer +' '+ hl_channel +' '+ hl_chaincode +' '+ key +' '+ value +' '+ hl_dockerid;
	console.log('[hlUtils.js] Ready to execute the PUT command: ' + command);

	// example call script on fabric vm: ./hl_put_test.sh 0 mychannel keyValueStore k10 v10 1db78d8$
	try {
		exec(command, {
						user: hl_user,
						host: hl_ip,
						password: hl_pass
					}).pipe(process.stdout);
	} catch(e) {
			console.log(e);
	}

	console.log('[hlUtils.js] PUT lunched toward Hyperledger Fabric');
	return callback("ok");
}


/* Hyperledger Fabric - INVOKE */
var hl_invoke = exports.hl_invoke =  function(myobj, callback) {
	console.log('[hlUtils.js] Calling INVOKE toward Hyperledger Fabric...');
	_invoke(myobj, function(res){
	      return callback(res);
	  });
}

function _invoke (myobj, callback){
	console.log('[hlUtils.js] Preparing INVOKE params..');
	
	//peer info
	var peer = undefined;
	var peer_ip = undefined;
	var peer_user = undefined;
	var peer_pass = undefined;

	var channel = undefined;
	var chaincode_name = undefined;
	if ('peer' in myobj && myobj.peer != ''){
	  peer = myobj.peer;
	  //take the dockerid cli, username, id, password related to the peer inserted
	  dockerId = config.get('out-service.hl_peer_'+peer+'_cli_id');
	  peer_ip = config.get('out-service.hl_peer_'+peer+'_ip');
	  peer_user = config.get('out-service.hl_peer_'+peer+'_user');
	  peer_pass = config.get('out-service.hl_peer_'+peer+'_password');
	}
	else{
	  console.log("[hlUtils.js] No peer found in config file. Setting DEFAULT PEER");
	  peer_ip = hl_ip;
	  peer_user = hl_user;
	  peer_pass = hl_pass;
	  peer = hl_peer;	//default one
	  dockerId = hl_dockerid;	//default one
	}

	if ('channel' in myobj && myobj.channel != '' ){
	  channel = myobj.channel;
	}
	else{
	  console.log("[hlUtils.js] No channel found in config file. Setting DEFAULT CHANNEL");
	  channel = hl_channel;		//default one
	}

	if ('chaincodeName' in myobj && myobj.chaincodeName != '' ){
		chaincode_name = myobj.chaincodeName;
	}
	else{
		console.log("[hlUtils.js] No chaincode found in config file. Setting DEFAULT CHAINCODE");
		chaincode_name = hl_chaincode;     //default one
	}

	var fcn = myobj.fcn;
	var args = myobj.args;

	//merge of the two strings in myobj -> e.g. "put,key,val"
	var fcnargs = fcn+','+args;
	var command = hl_script_path + 'hl_invoke.sh' +' '+ peer +' '+ channel +' '+ chaincode_name +' '+fcnargs+ ' '+ dockerId;
	console.log('[hlUtils.js] Ready to execute the INVOKE command: ' + command);
	console.log('[hlUtils.js] Destination PEER: ' + peer_user +'@'+ peer_ip);
	// example call script on fabric vm: ./hl_put_test.sh 0 mychannel keyValueStore k10 v10 1db78d826131
	try {
	  exec(command, {
			user: peer_user,
			host: peer_ip,
			password: peer_pass
		}).pipe(process.stdout);
	} catch(e) {
		console.log(e);
	}

	console.log('[hlUtils.js] INVOKE lunched toward Hyperledger Fabric');
	
	// a timeout is needed to avoid a reading of an old file
	setTimeout( function(){
		
		var client = require('scp2')
		client.scp({
			host: peer_ip,
			username: peer_user,
			password: peer_pass,
			path: hl_script_path + 'result.log'
		}, './', function(err) {
			if (err) {
				console.log(err);
				return;
			}
			
			var fs = require('fs');
			//fs.readFile( hl_script_path + 'result.log', function (err, data) {
			fs.readFile( './result.log', function (err, data) {	
				if (err) {
					console.log(err);
					return;
				}
		
				//clean results
				var string_data = data.toString();
				var response_res = string_data.substr(string_data.indexOf('response:'), 1000);
                console.log('[hl_utils.js] response: ' + response_res);

				response_res = response_res.replace(/\"/g, "");
				response_res = response_res.replace(/\\/g, "'");
				
				//kind of response checking 
				try {
					response_res = response_res.split("response:<")[1].split(" >")[0].replace("\"", "").replace("\\","");
				}
				catch(err) {
					response_res = "status:200 message:OK";
				}	
				
				// GET
				/*if ( response_res.includes("payload") ) {
					response_res = response_res.split('payload:')[1];
					
				} else {
					response_res = response_res.split("message:")[1];
				}*/
				if ( !response_res.includes("payload") ) {
					response_res = response_res + ' payload:NULL';
				}
				return callback(response_res);
			
			})
			
		})}, 5000);

}


/* Hyperledger Fabric - GET */
var hl_get = exports.hl_get =  function(myobj, callback) {	
    console.log('[hlUtils.js] Calling GET operation toward Hyperledger Fabric...');
	_get(myobj, function(res){
	      return callback(res);
	  });
}

function _get (myobj, callback){
	console.log('[hlUtils.js] Executing GET with following key:');
	
	var key = myobj.key;
	console.log("[hlUtils.js] GET key:"+key);
	var command = hl_script_path + 'hl_get.sh ' +' '+ hl_peer +' '+ hl_channel +' '+ hl_chaincode +' '+ key +' '+ hl_dockerid;
	console.log('[hlUtils.js] Ready to execute the GET command: ' + command);

	// example call script on fabric vm: ./hl_put_test.sh 0 mychannel keyValueStore k10 v10 1db78d8$
	try {
		exec(command, {
						user: hl_user,
						host: hl_ip,
						password: hl_pass
					}).pipe(process.stdout);
	} catch(e) {
			console.log(e);
	}

	console.log('[hlUtils.js] GET lunched toward Hyperledger Fabric');
	
	// a timeout is needed to avoid a reading of an old file
	setTimeout( function(){
		
		var client = require('scp2')
		client.scp({
			host: hl_ip,
			username: hl_user,
			password: hl_pass,
			path: hl_script_path + 'result-get.log'
		}, './', function(err) {
			if (err) {
				throw err; 
			}
			
			var fs = require('fs');
			//fs.readFile( hl_script_path + 'result.log', function (err, data) {
			fs.readFile( './result-get.log', function (err, data) {	
				if (err) {
					throw err; 
				}
		
				//clean results
				var string_data = data.toString();
				var response_res = string_data.substr(string_data.indexOf('response:'), 100);
				response_res = response_res.replace(/\"/g, "");
				
				//kind of response checking 
				response_res = response_res.split("response:<")[1].split(" >")[0].replace("\"", "");
				// GET
				/*if ( response_res.includes("payload") ) {
					response_res = response_res.split('payload:')[1];
					
				} else {
					response_res = response_res.split("message:")[1];
					if (response_res === 'OK') {
						response_res = 'NULL';
					}	
				}*/
				if ( !response_res.includes("payload") ) {
					response_res = response_res + ' payload:NULL'
				}				
				return callback(response_res);
			
			})
			
		})}, 5000);
	
}


/* Hyperledger Fabric - DELETE */
var hl_delete = exports.hl_delete =  function(myobj, callback) {	
	console.log('[hlUtils.js] Calling DELETE toward Hyperledger Fabric...');    
	_delete(myobj, function(res){
	      return callback(res);
	  });
}

function _delete (myobj, callback){
	console.log('[hlUtils.js] Executing DELETE with following key:');

	var key = myobj.key;
	console.log("[hlUtils.js] DELETE key:"+key);
	var command = hl_script_path + 'hl_delete.sh ' +' '+ hl_peer +' '+ hl_channel +' '+ hl_chaincode +' '+ key +' '+ hl_dockerid;
	console.log('[hlUtils.js] Ready to execute the DELETE command: ' + command);

	// example call script on fabric vm: ./hl_put_test.sh 0 mychannel keyValueStore k10 v10 1db78d8$
	try {
		exec(command, {
						user: hl_user,
						host: hl_ip,
						password: hl_pass
					}).pipe(process.stdout);
	} catch(e) {
			console.log(e);
	}

	console.log('[hlUtils.js] DELETE lunched toward Hyperledger Fabric');
	return callback("ok");
}
