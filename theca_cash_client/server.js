const webpack = require('webpack');
const WebpackDevServer = require('webpack-dev-server');
const config = require('./webpack.dev');
const bodyParser = require('body-parser');
// const cors = require('cors');
// const store = require('store');

const HOST_ADDRESS = process.env.HOST_ADDRESS || 'localhost';
const SERVER_PORT = process.env.SERVER_PORT || 3000;
const SERVER_PROTOCOL = process.env.SERVER_PROTOCOL || 'http';

new WebpackDevServer(webpack(config), {
    publicPath: config.output.publicPath,
    hot: true,
    historyApiFallback: true,
    disableHostCheck: true,
    // It suppress error shown in console, so it has to be set to false.
    quiet: false,
    // It suppress everything except error, so it has to be set to false as well
    // to see success build.
    noInfo: false,
    stats: {
        // Config for minimal console.log mess.
        assets: false,
        colors: true,
        version: false,
        hash: false,
        timings: false,
        chunks: false,
        chunkModules: false
    },
    setup: (app) => {
        app.use(bodyParser.json());
        app.use(bodyParser.urlencoded({
            extended: true
        }));

        // app.use(cors());
    }
}).listen(SERVER_PORT, HOST_ADDRESS, (err) => {
    if (err) {
        console.log(err);
    }

    console.log('Listening at ' + HOST_ADDRESS + ':' + SERVER_PORT);
});
