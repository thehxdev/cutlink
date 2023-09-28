# Docker Compose with Nginx

Example and documentation for running Cutlink with nginx as a reverse proxy with TLS.


## Issue SSL certificate for TLS Encryption

### Self-signed certificate

> [!WARNING]
> Use self-signed SSL certificates ONLY for development environment and NOT for production.

while you are in this directory, run following command to generate `cert.pem` and `key.pem` files in `./nginx/ssl` directory.
```bash
mkdir -p ./nginx/ssl

openssl req -x509 -newkey rsa:4096 -nodes -keyout ./nginx/ssl/key.pem -out ./nginx/ssl/cert.pem -days 365 -subj '/CN=localhost'
```


### ACME script

For production, buy a domain name and setup DNS recordes for that. then issue a SSL certificate for that with ACME script.

> [!WARNING]
> Run these commands as `root` user.

> [!NOTE]
> Make sure that `curl` and `socat` programs are already installed.

```bash
# install official ACME installer script
curl https://get.acme.sh | bash

# use LetsEncrypt service
/root/.acme.sh/acme.sh --set-default-ca --server letsencrypt

# register your email (you can use temp mails)
/root/.acme.sh/acme.sh --register-account -m YOUR_EMAIL

# issue certificates (Port 80 MUST be open)
# replace `YOUR_DOMAIN` with your actual domain pointing to your server
/root/.acme.sh/acme.sh --issue -d YOUR_DOMAIN --standalone

# make a directory to store SSL certificates
# replace `YOUR_DOMAIN` with your actual domain pointing to your server
mkdir -p /root/.ssl/YOUR_DOMAIN

# install SSL certificates to /root/.ssl directory
# replace `YOUR_DOMAIN` with your actual domain pointing to your server
/root/.acme.sh/acme.sh --installcert -d YOUR_DOMAIN --key-file /root/.ssl/YOUR_DOMAIN/key.pem --fullchain-file /root/.ssl/YOUR_DOMAIN/cert.pem
```

When you installed certificates, copy them to `./nginx/ssl` directory.
```bash
mkdir -p ./nginx/ssl

# replace `YOUR_DOMAIN` with your actual domain pointing to your server
cp /root/ssl/YOUR_DOMAIN/*.pem ./nginx/ssl/
```


## Configuration

By default if you run `docker compose` without configuration, the app will listen on port `8443`.
If you want to change the listening port for server (for example port 443), edit `compose.yaml` file:

```yaml
  # ...
  nginx:
    image: nginx
    volumes:
      - $PWD/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - $PWD/nginx/default.conf:/etc/nginx/conf.d/default.conf
      - $PWD/nginx/ssl:/ssl
    ports:
      - 443:443 # changed the left side number to 443
  # ...
```

That's it!


## Run Cutlink

When you have SSL certificates in correct directory (`./nginx/ssl`) and configured `compose.yaml` file, you can run the container:

```
docker compose up -d
```
