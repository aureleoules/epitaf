server {
    listen   80;
    listen   [::]:80 default ipv6only=on;

    root /usr/share/nginx/html;
    index index.html;

    server_name _; # all hostnames

    location / {
        try_files $uri /index.html;
    }
}