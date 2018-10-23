const path = require('path');
const webpack = require("webpack");
const HtmlWebPackPlugin = require("html-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CleanWebpackPlugin = require('clean-webpack-plugin');
const Dotenv = require('dotenv-webpack');

const HOST_ADDRESS = process.env.HOST_ADDRESS || 'localhost';
const SERVER_PORT = process.env.SERVER_PORT || 3000;
const SERVER_PROTOCOL = process.env.SERVER_PROTOCOL || 'http';

const query={
    mozjpg: {
        progressive: true,
    },
    gifsicle:{
        interlaced: false,
    },
    optipng:{
        optimizationLevel: 4,
    },
    pngquant:{
        quality: '75-90',
        speed: 3,
    },
};

const svgoConfig = JSON.stringify({
    plugins: [
        {removeTitle: true},
        {convertColors: {shorthex: false}},
        {convertPathData: false}
    ]
});


module.exports = {
    resolve: {
        extensions: ['*', '.js', '.jsx'],
        alias: {
            'load-image': 'blueimp-load-image/js/load-image.js',
            'load-image-meta': 'blueimp-load-image/js/load-image-meta.js',
            'load-image-exif': 'blueimp-load-image/js/load-image-exif.js',
            'canvas-to-blob': 'blueimp-canvas-to-blob/js/canvas-to-blob.js',
            'jquery-ui/widget': 'blueimp-file-upload/js/vendor/jquery.ui.widget.js',
            'load-image-scale': 'blueimp-load-image/js/load-image-scale.js',
            'jquery.ui.widget': 'node_modules/jquery.ui.widget/jquery.ui.widget.js',
            'jquery-ui':'jquery-ui/ui'
        },
    },
    output: {
        path: path.join(__dirname, '/dist/'),
        filename: '[name].js',
        publicPath: '/'
    },
    module: {
        rules: [
            {
                enforce: "pre",
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                loader: "eslint-loader",
            },
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ['babel-loader']
            },
            {
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader",
                        options: { minimize: true }
                    }
                ]
            },
            {
                test: /\.s?css$/,
                loaders: ['style-loader', 'css-loader'],
            },
            {
                test: /\.woff(2)?(\?[a-z0-9#=&.]+)?$/,
                loader: 'url-loader?limit=10000&mimetype=application/font-woff'
            },
            {
                test: /\.(ttf|eot)(\?[a-z0-9#=&.]+)?$/,
                loader: 'file-loader'
            },
            {
                test: /\.svg(\?[a-z0-9#=&.]+)?$/,
                loaders: [
                    'file-loader',
                    'svgo-loader?' + svgoConfig
                ]
            },
            {
                test: /\.(jpe?g|png|gif)$/i,
                loaders: [
                    'file-loader?hash=sha512&digest=hex&name=[hash].[ext]',
                    `image-webpack-loader?${JSON.stringify(query)}`
                ],
            }
        ]
    },
    plugins: [
        new CleanWebpackPlugin(['dist']),
        new HtmlWebPackPlugin({
            title: "theca.cash",
            template: "./app/index.tpl.html",
            inject: 'body',
            filename: 'index.html'
        }),
        new MiniCssExtractPlugin({
            filename: "[name].css",
            chunkFilename: "[id].css"
        }),
        new webpack.optimize.OccurrenceOrderPlugin(),
        new webpack.ProvidePlugin({
            $: "jquery",
            jQuery: "jquery"
        }),
    ],
    devtool: 'inline-source-map',
    devServer: {
        contentBase: './dist',
        host: HOST_ADDRESS,
        port: SERVER_PORT,
        historyApiFallback: true,
        open: true,
        hot: true
    }
};