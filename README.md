# CutLink!

A simple link shortener written in golang with focuse on privacy and anonymity.

> [!NOTE]
> This project is under active development.


## Features

Link shortener services are simple. But CutLink is different in some points.

- No email or personal information needed for sign-up. Your username is a randomly generated UUIDv4 that will show you only **one** time after you specified your password.
- You own your data. if you delete your account, **ALL** of your data will be deleted from database (Real delete :D).
- Only server-side errors are logged (for troubleshooting).
- You can build your own server if you don't trust others.
- Simple UI.


## Build

> [!NOTE]
> Make sure that you already installed `gcc`, `make`, `sqlite3` and `go`.

- Build the executable:
```bash
# This command will build the project and
# outputs the executable to ./bin directory.
make
```

- Build docker image:
```bash
# build docker image
make docker

# make executable with golang docker container.
# This is same as make command but builds the project inside golang docker container.
make docker_exe
```


## TLS Support

At this time, for TLS support you have to run the executable or docker image on the localhost (`127.0.0.1`) then use a
**reverse proxy** like [Caddy](https://caddyserver.com/) to pass all trafic from `0.0.0.0` to CutLink with TLS encryption.

- To run CutLink on localhost:
```bash
# executable (on port 5000)
./bin/cutlink -addr 'localhost:5000'

# docker image (after you build it)
docker run -d -p 127.0.0.1:5000:5000 cutlink:latest
```


## How this could be better?

Since this is my first web development project after one year of Python and C programming, some bad practices or bad coding style
is expected. I try to learn more and make this project better.

- Create docker-compose file for better experience with docker.
- Better error handling on client-side errors (All of them handled with `Internal Server Error`).
- Better (but VERY simple) UI (First time building front-ends).
- Check for security vulnerabilities (Sessions and CSRF are already present).
- Better form information checking.


## Features to implement

- Date or Click limit for links
- Password protection for links
- Admin control for better server management
