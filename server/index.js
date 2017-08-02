'use strict';

var fs = require('fs'),
    path = require('path'),
    http = require('http');

var app = require('connect')();
var assert = require('assert');
var swaggerTools = require('swagger-tools');
var jsyaml = require('js-yaml');
var serverPort = 8081;

// output service configuration
var config = require('config');
var out_service_name = config.get('out-service.name');
var out_service_ip = config.get('out-service.ip');
var out_service_port = config.get('out-service.port');
var db_name= config.get('out-service.dbname');
var db_collection = config.get('out-service.dbcollection');

// swaggerRouter configuration
var options = {
  swaggerUi: path.join(__dirname, '/swagger.json'),
  controllers: path.join(__dirname, './controllers'),
  useStubs: process.env.NODE_ENV === 'development' // Conditionally turn on stubs (mock mode)
};

// The Swagger document (require it, build it programmatically, fetch it from a URL, ...)
var spec = fs.readFileSync(path.join(__dirname,'api/swagger.yaml'), 'utf8');
var swaggerDoc = jsyaml.safeLoad(spec);


// Initialize the Swagger middleware
swaggerTools.initializeMiddleware(swaggerDoc, function (middleware) {

  // Interpret Swagger resources and attach metadata to request - must be first in swagger-tools middleware chain
  app.use(middleware.swaggerMetadata());

  // Validate Swagger requests
  app.use(middleware.swaggerValidator());

  // Route validated requests to appropriate controller
  app.use(middleware.swaggerRouter(options));

  // Serve the Swagger documents and Swagger UI
  app.use(middleware.swaggerUi());

  // Start the server
  http.createServer(app).listen(serverPort, function () {
    
    if (out_service_name == 'mongo') {

        var MongoClient = require('mongodb').MongoClient
          , assert = require('assert');

        // Connection URL
        var url = 'mongodb://' + out_service_ip + ':' + out_service_port + '/' + db_name;

        // Use connect method to connect to the server
        MongoClient.connect(url, function(err, db) {
          assert.equal(null, err);
          console.log("Connected successfully to mongodb");
		  var myobj = { name: "Company Inc", address: "Highway 37" };
          db_put(db, myobj, function() {
            db.close();
  	  	  });
        });
    }

    console.log('Your server is listening on port %d (http://localhost:%d)', serverPort, serverPort);
    console.log('Swagger-ui is available on http://localhost:%d/docs', serverPort);
  });

});




