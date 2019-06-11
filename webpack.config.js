const path = require('path');

module.exports = {
    entry: './front/src/index.ts',
    mode: 'development',
    devtool: 'source-map',
    module: {
        rules: [
            {
                include: /\.tsx?$/,
                exclude: /node_modules/,
                use: 'ts-loader',
            }
        ]
    },
    resolve: {
        extensions: [ '.tsx', '.ts', '.js' ]
    },
    output: {
        filename: 'index.js',
        path: path.resolve(__dirname, 'front')
    }
};

