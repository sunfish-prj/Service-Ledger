
// output service configuration
var config = require('config');
var assert = require('assert');

var exec = require('ssh-exec');


/* Hyperledger Fabric - PUT */
var hl_put = exports.hl_put =  function(myobj, callback) {	
    exec('ls -lh', {
	  user: 'ubuntu',
	  host: 'localhost:2222',
	  password: 'ubuntu'
	}).pipe(process.stdout)
}

function _put (myobj, db, callback){
	console.log('Executing put');
	var collection = db.collection(db_collection);
	// Insert a document
	collection.update({"key": myobj.key},
		 	myobj, {upsert: true}, function(err, res){
    	if (err) {
			console.log(err);
			return new Error();
		}
		console.log('Put succeeded');
		return callback("ok");
   })
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
