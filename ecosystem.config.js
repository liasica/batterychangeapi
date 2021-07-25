module.exports = {
    apps: [
        {
            script: 'start.sh',
            watch: ['app', 'boot', 'config', 'deploy', 'library', 'packed', 'router', 'swagger', 'main.go'],
            name: 'apiv2',
            watch_delay: 1000,
            cwd: '/var/www/apiv2',
            ignore_watch: ['uploads'],
            watch_options: {
                followSymlinks: false
            }
        },
        {
            script: 'appupdateserver',
            watch: false,
            name: 'update',
            cwd: '/var/www/app.shiguangjv.com',
        }
    ]
}