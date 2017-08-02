
// output service configuration
var config = require('config');
var assert = require('assert');
var out_service_ip = config.get('out-service.ip');
var out_service_port = config.get('out-service.port');
var db_name = config.get('out-service.dbname');
var db_collection = config.get('out-service.dbcollection');

var MongoClient = require( 'mongodb' ).MongoClient;
var url = 'mongodb://' + out_service_ip + ':' + out_service_port + '/' + db_name;

var db_put = exports.db_put =  function(myobj, callback) {	
    MongoClient.connect(url, function(err, db) {
	  assert.equal(null, err);
	  console.log(myobj);
      console.log("Connected successfully to mongodb");
	  _put(myobj, db);
      db.close();
    });
}

function _put (myobj, db){
	console.log('Executing put');
	var collection = db.collection(db_collection);
	// Insert a document
	var _res;
	return collection.insertOne(myobj, function(err){
    	if (err) return;
		console.log('Put succeeded ' + myobj._id);
//		console.log('Put succeeded ' + res.ops[0]._id);  
//		console.log('Put succeeded ' + JSON.stringify(res,null,4));  
   })
}
