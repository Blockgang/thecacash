const express = require('express');
// var https = require('https');
const http = require('http');
const bodyParser = require('body-parser');
const cors = require('cors');
const store = require('store');

const SERVER_ADDRESS = process.env.INTERNAL_ADDRESS_SERVER || 'localhost';
const SERVER_PORT = process.env.INTERNAL_PORT_SERVER || 8000;
const PORT = process.env.INTERNAL_PORT_CLIENT || 3000;
const SERVER_PROTOCOL = 'http';


const app     = express();
app.use(express.static(__dirname + '/app'));

app.get('/', function(req, res) {
    res.sendFile('index.html');
});

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: true
}));

// Create an HTTP service.
http.createServer(app).listen(PORT);
// Create an HTTPS service identical to the HTTP service.
// https.createServer(options, app).listen(3001);


console.log('Running at Port ' + PORT);
