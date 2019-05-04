const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

const mode = process.env.NODE_ENV;
const outputPath = mode === 'production' ? '../assets/' : '../server/assets/';

module.exports = {
    devtool: 'eval-source-map',
    mode: mode,
    entry: './typescript/src/main.tsx',
    output: {
        filename: 'js/bundle.js',
        path: path.resolve(__dirname, outputPath),
        publicPath: "/"
    },
    resolve: {
        extensions: ['.ts', '.tsx', '.js']
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'css/index.css'
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
                        loader: MiniCssExtractPlugin.loader
                    },
                    'css-loader',
                    'sass-loader'
                ]
            },
            {
                test: /\.(png|jpg|gif|svg)$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: 'gen/[name].[ext]',
                        }
                    }
                ]
            },
        ]
    }
};
