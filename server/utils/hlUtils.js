
// output service configuration
var config = require('config');
var assert = require('assert');

var exec = require('ssh-exec');
var hl_user = config.get('out-service.hl_user');
var hl_pass = config.get('out-service.hl_password');
var hl_ip = config.get('out-service.hl_ip');


/* Hyperledger Fabric - PUT */
var hl_put = exports.hl_put =  function(myobj, callback) {	
	console.log('invoking ssh');    
	_put(myobj, function(res){
	      return callback(res);
	  });
}

function _put (myobj, callback){
	console.log('Executing put');
	
	var channel = "mychannel";
	var chaincode = "mycc";
	var key = myobj.key;
	var value = myobj.value;
	var command = './hl_put.sh ' + channel + ' ' + chaincode + ' ' + key + ' ' + value;
	console.log('command: ' + command);
	
	exec('ls -l', {
	  user: hl_user,
	  host: hl_ip,
	  password: hl_pass
	}).pipe(process.stdout)
	
	console.log('Put succeeded');
	return callback("ok");

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
	// Insert a document
	collection.findOne(_id, function(err, res){
    	if (err) return;
		try{
			console.log('Get succeeded! Value: ' + res.value);
			return callback(res.value);
		} catch (err) {
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
