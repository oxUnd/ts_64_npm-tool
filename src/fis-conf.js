fis.config.set('deploy', {
    "remote": [
        {
            from: "/static",
            to: "../public/"
        },
        {
            from: "/template",
            to: ".."
        }
    ]
});

fis.config.set('pack', {
    "static/aio.js": "**.js",
    "static/aio.css": "**.css"
});