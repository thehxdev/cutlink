# CutLink!

A simple link shortener written in golang with focuse on privacy and anonymity.

> [!NOTE]
> This project is under active development.


> [!WARNING]
> I did not red any article or books about this kind of programs and also not even watched a tutorial video. I wrote this program with my own thoughts about how a link shortener _may_ work.


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
# This command will build the project and
# outputs the executable to ./bin directory.
make docker
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
- Better error handling on client-side errors (All of them handled with `Internal Server Error`)
- Better (but VERY simple) UI (First time building front-ends)
- Check for security vulnerabilities (Sessions and CSRF are already present)


## Features to implement

- Date or Click limit for links
- Password protection for links
- Admin control for better server management
