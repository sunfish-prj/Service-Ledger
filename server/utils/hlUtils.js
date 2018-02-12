
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
	console.log('invoking ssh');    
	_put(myobj, function(res){
	      return callback(res);
	  });
}

function _put (myobj, callback){
	console.log('Executing put');

	
	var key = myobj.key;
	var value = myobj.value;
	console.log("[hlUtils.js] PUT key:"+key+" value:"+value);
	var command = hl_script_path + 'hl_put.sh ' +' '+ hl_peer +' '+ hl_channel +' '+ hl_chaincode +' '+ key +' '+ value +' '+ hl_dockerid;
	console.log('command: ' + command);

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

	console.log('Put succeeded');
	return callback("ok");
}


/* Hyperledger Fabric - INVOKE */
var hl_invoke = exports.hl_invoke =  function(myobj, callback) {
	console.log('[hlUtils.js] start invoking hyperledger');
	_invoke(myobj, function(res){
	      return callback(res);
	  });
}

function _invoke (myobj, callback){
	console.log('[hlUtils.js] preparing invoke');
	
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
	  console.log("You are using DEFAULT PEER");
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
	  console.log("[hlUtils.js] setting DEFAULT CHANNEL");
	  channel = hl_channel;		//default one
	}

	if ('chaincodeName' in myobj && myobj.chaincodeName != '' ){
		chaincode_name = myobj.chaincodeName;
	}
	else{
		console.log("[hlUtils.js] setting DEFAULT CHAINCODE");
		chaincode_name = hl_chaincode;     //default one
	}

	var fcn = myobj.fcn;
	var args = myobj.args;

	//merge of the two strings in myobj -> e.g. "put,key,val"
	var fcnargs = fcn+','+args;
	var command = hl_script_path + 'hl_invoke.sh' +' '+ peer +' '+ channel +' '+ chaincode_name +' '+fcnargs+ ' '+ dockerId;
	console.log('[hlUtils.js] command: ' + command);

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

	
	// a timeout is needed to avoid a reading of an old file
	console.log('[hlUtils.js] Invoke succeeded');
	
	setTimeout( function(){
		
		var client = require('scp2')
		client.scp({
			host: peer_ip,
			username: peer_user,
			password: peer_pass,
			path: hl_script_path + 'result.log'
		}, './', function(err) {
			if (err) {
				throw err; 
			}
			
			var fs = require('fs');
			//fs.readFile( hl_script_path + 'result.log', function (err, data) {
			fs.readFile( './result.log', function (err, data) {	
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
				if ( response_res.includes("payload") ) {
					response_res = response_res.split('payload:')[1];
					
				} else {
					response_res = response_res.split("message:")[1];
				}
				return callback(response_res);
			
			})
			
		})}, 2000);


}


/* Hyperledger Fabric - GET */
var hl_get = exports.hl_get =  function(_id, callback) {	
    MongoClient.connect(url, function(err, db) {
	  assert.equal(null, err);
      console.log("Connected successfully to mongodb");
	  _get(_id, db, function(res){
	      db.close();
 	  	  return callback(res);
	  });
    });
}

function _get (_id, db, callback){
	console.log('Executing get');
	var collection = db.collection(db_collection);


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

	console.log('Invoke succeeded');
	//
	return callback("ok");

	//

	collection.findOne(_id, function(err, res){
    	if (err) return;
			try{
				console.log('Get succeeded! Value: ' + res.value);
				return callback(res.value);
			}	catch (err) {
				return callback(new Error());
			}
   })
}


/* Hyperledger Fabric - DELETE */
var hl_delete = exports.db_delete =  function(myobj, callback) {	
    MongoClient.connect(url, function(err, db) {
	  assert.equal(null, err);
      console.log("Connected successfully to mongodb");
	  _delete(myobj, db, function(res){
	      db.close();
 	  	  return callback(res);
	  });
    });
}

function _delete (myobj, db, callback){
	console.log('Executing delete');
	var collection = db.collection(db_collection);
	// Insert a document
	collection.deleteOne(myobj, function(err, res){
    	if (err){
    		return new Error();
    	}
		console.log('Delete succeeded');
		return callback("ok");
   })
}
