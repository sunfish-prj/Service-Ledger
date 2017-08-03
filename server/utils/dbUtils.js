
// output service configuration
var config = require('config');
var assert = require('assert');
var out_service_ip = config.get('out-service.ip');
var out_service_port = config.get('out-service.port');
var db_name = config.get('out-service.dbname');
var db_collection = config.get('out-service.dbcollection');

var MongoClient = require( 'mongodb' ).MongoClient;
var url = 'mongodb://' + out_service_ip + ':' + out_service_port + '/' + db_name;


/* DB - PUT */
var db_put = exports.db_put =  function(myobj, callback) {	
    MongoClient.connect(url, function(err, db) {
	  assert.equal(null, err);
      console.log("Connected successfully to mongodb");
	  _put(myobj, db, function(res){
	      db.close();
 	  	  return callback(res);
	  });
    });
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


/* DB - GET */
var db_get = exports.db_get =  function(_id, callback) {	
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


/* DB - DELETE */
var db_delete = exports.db_delete =  function(myobj, callback) {	
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