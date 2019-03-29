const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
    devtool: 'eval-source-map',
    mode: 'development',
    entry: './typescript/src/main.tsx',
    output: {
        filename: 'js/bundle.js',
        path: path.resolve(__dirname, '../assets/')
    },
    resolve: {
        extensions: ['.ts', '.tsx', '.js']
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "css/index.css"
        })
    ],
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                loader: 'ts-loader'
            },
            {
                test: /\.scss$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                        options: {
                            // you can specify a publicPath here
                            // by default it use publicPath in webpackOptions.output
                            publicPath: '../'
                        }
                    },
                    'css-loader',
                    'sass-loader'
                ]
            }
        ]
    }
};
