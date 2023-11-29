# Cutlink!

A privacy and anonymity focused link shortener written in Golang.


## Features

Link shortener services are simple. But Cutlink is different in some points.

- Safe Mode! (See the target URL before redirect process BUT this is optional).
- No email or personal information needed for signing up. Your username is a randomly generated UUIDv4 that will show you only **one** time after you submited your password in signup page.
- You own your data. if you delete your account, **ALL** of your data will be deleted (Real delete :D).
- Only server-side errors are logged (for troubleshooting).
- You can build your own server if you don't trust others.


## Configuration

Cutlink uses `toml` file format for it's config file.
When you execute Cutlink, it will check `/etc/cutlink` and `./` directories in order to find `config.toml` file.
If it's not found, the hole program panics. [You can find default config file here.](https://github.com/thehxdev/cutlink/tree/main/config.toml)

If you're using docker compose, you don't need to change the config file. Otherwise take a look at it to fit your needs.


## Build

### Build docker image (Recommended)
```bash
# build docker image
make docker
```

After you built the docker image, you can run the image with [Docker Compose and Nginx](https://github.com/thehxdev/cutlink/tree/main/docs/docker-compose-examples/with-nginx).

If you want to run the image with `docker` command directly without a reverse proxy and TLS support:
```bash
docker run -d --name cutlink -v "$PWD"/config.toml:/etc/cutlink/config.toml -p 5000:5000 cutlink:latest
```

> [!NOTE]
> For examples of docker compose usage, go to [Docker Compose Examples and Instructions](https://github.com/thehxdev/cutlink/tree/main/docs/docker-compose-examples).


### Build the executable
> [!NOTE]
> Make sure that you already installed `gcc`, `make`, `sqlite3` and `go`.

```bash
# This command will build the project and outputs the executable to ./bin directory.
make
```


### Build the executable with Golang docker image

This will use Golang docker image to compile the executable and saves the output to `./bin` directory in local environment.
Same as running `make` but inside a docker container.

```bash
# remove .git directory
rm -rf .git/

make docker_exe
```


## Features to implement

- Date or Click limit for links
- Admin control panel for better server management
