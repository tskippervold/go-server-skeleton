const path = require('path');
const CopyPlugin = require('copy-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

module.exports = {
    entry: [/*'./src/styles.scss', */'./src/checkout-inline/index.js'],
    output: {
        filename: 'index.js'
    },
    /*module: {
        rules: [
            {
                test:/\.(s*)css$/,
                use: [
                  // Creates `style` nodes from JS strings
                  'style-loader',
                  // Translates CSS into CommonJS
                  'css-loader',
                  // Compiles Sass to CSS
                  'sass-loader',
                ],
            }
        ]
    },*/
    plugins: [
        new CleanWebpackPlugin(),
        new CopyPlugin({
            patterns: [
                {
                    from: 'src/**/*',
                    transformPath(target) {
                        return target.replace('src/', '')
                    }
                }
            ]
        })
    ]
}