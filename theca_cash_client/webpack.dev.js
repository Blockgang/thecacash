const merge = require('webpack-merge');
const webpack = require('webpack');
const path = require('path');
const common = require('./webpack.config.js');

const HOST_ADDRESS = process.env.HOST_ADDRESS || 'localhost';
const SERVER_PORT = process.env.SERVER_PORT || 3000;
const SERVER_PROTOCOL = process.env.SERVER_PROTOCOL || 'http';
module.exports = merge(common, {
    entry: [
        'babel-polyfill',
        'webpack-dev-server/client?'+ SERVER_PROTOCOL+'://'+ HOST_ADDRESS + ':' + SERVER_PORT,
        'webpack/hot/only-dev-server',
        'react-hot-loader/patch',
        path.join(__dirname, 'app/index.js')
    ],
    plugins: [
        new webpack.HotModuleReplacementPlugin(),
        new webpack.NoEmitOnErrorsPlugin(),
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('development')
        }),

    ],
  devtool: 'inline-source-map',
  devServer: {
    contentBase: './dist'
  }
});
