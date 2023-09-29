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


### Build the executable
> [!NOTE]
> Make sure that you already installed `gcc`, `make`, `sqlite3` and `go`.

```bash
# This command will build the project and sqlite database
# and outputs the executable to ./bin directory.
make all
```


### Build docker image (Recommended)

```bash
# build docker image
make docker
```

After you built the docker image, you can run the image with [Docker Compose](https://github.com/thehxdev/cutlink/tree/main/docs/docker-compose-examples).

> [!NOTE]
> For examples of docker usage go to [Docker Compose Examples and Istructions](https://github.com/thehxdev/cutlink/tree/main/docs/docker-compose-examples).


If you want to run the image with `docker` command directly without a reverse proxy and TLS support:
```bash
docker run -d --name cutlink -v "$PWD"/config.toml:/etc/cutlink/config.toml -p 5000:5000 cutlink:latest
```


### Build the executable with golang docker container
```bash
# remove .git directory
rm -rf .git/

# This is same as make command but builds the project inside golang docker container
# and outputs the executable to ./bin directory
make docker_exe
```


## How this could be better?

Since this is my first web development project after one year of Python and C programming, some bad practices or bad coding style
is expected. I try to learn more and make this project better.

- Check for security vulnerabilities (Sessions and CSRF protection are already present).
- Better form information checking.


## Features to implement

- Date or Click limit for links
- Admin control for better server management
