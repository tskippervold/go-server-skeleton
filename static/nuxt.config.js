export default {
    target: 'static',
    components: true,
    buildModules: [
        '@nuxtjs/tailwindcss',
    ],
    modules: [
        ['@nuxtjs/component-cache', { maxAge: 1000 * 60 * 60 }],
    ],
    head: {
        script: [
            {
                src: 'https://static.bambora.com/checkout-sdk-web/latest/checkout-sdk-web.min.js',
                defer: true
            }
        ]
    }
}