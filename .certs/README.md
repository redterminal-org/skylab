Generate a self-signed certificate and key for serving skylab on localhost:

```sh
cd .certs
openssl req -x509 -newkey rsa:4096 -nodes -days 365 \
  -keyout key.pem -out cert.pem \
  -subj /CN=localhost -addext subjectAltName=DNS:localhost
```

The generated `*.pem` files are ignored by git; don't commit them.
