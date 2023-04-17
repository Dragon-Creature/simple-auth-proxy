# Simple Auth Proxy
![Screenshot](Screenshot.png)

This project solves the problem when you want to expose an existing application to the internet but the application has no built-in authentication. You could configure basic auth in your load balancer but basic auth doesn't work well with password managers, and it doesn't remember that you were logged in as soon as you close the browser.

In that spirit, the tool is compatible with htpasswd, the common tool used to generate a password file for basic auth with Apache.

## Environment variable

* TARGET_PROTOCOL: specifics if the source is http:// or https://, default is `http://`.
* TARGET_URL: specifics the address to the source, default is `localhost:30085`.
* HTPASSWD_FILE: the password file, default is `htpasswd`.
* COOKIE_MAX_AGE: the max age of a session after login in, in seconds, default is `86400` or 1 day.

## Supported features

* Supports both http and https traffic.
* Supports WebSockets in both directions but only text for the client to server, server to client support all flags.

## Running project
### Required software
* Docker
* NPM (only for non-docker build)
* Go (+1.19) (only for non-docker build)

### Start
Change the `your-username` to your username of choice and change the environment variables to fit your needs.

    htpasswd -c -B htpasswd your-username
    make
    docker run -p 8080:8080 -e TARGET_PROTOCOL=http:// -e TARGET_URL=localhost:30085 -e HTPASSWD_FILE=htpasswd -e COOKIE_MAX_AGE=86400 -v $(pwd)/htpasswd:/app/htpasswd docker.ssns.se/frozendragon498/simple-auth-proxy