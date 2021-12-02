var sdk = require('./sdk.js');
module.exports = function(app){
  app.get('/api/initWallet', function (req, res) {
    var name = req.query.name;
    var id = req.query.id;

    let args = [name, id];
    sdk.send(true, 'initWallet', args, res);
  });
  app.get('/api/getWallet', function (req, res) {
    var id = req.query.id;

    let args = [id];

    sdk.send(false, 'getWallet', args, res);
  });
  app.get('/api/issueMileage', function (req, res) {
    var id = req.query.id;
    var value = req.query.value;

    let args = [id, value];
    sdk.send(true, 'issueMileage', args, res);
  });
  app.get('/api/purchaseWithMileage', function (req, res) {
    var source = req.query.source;
		var destination = req.query.destination;
    var value = req.query.value;
    
    let args = [source, destination, value];
    sdk.send(true, 'purchaseWithMileage', args, res);
  });
  app.get('/api/resetMileage', function (req, res) {
    var id = req.query.id;

    let args = [id];
    sdk.send(true, 'resetMileage', args, res);
  });
  app.get('/api/getLog', function (req, res) {
    var id = req.query.id;

    let args = [id];
    sdk.send(false, 'getLog', args, res);
  });
}