# go-stable [![Build Status](https://travis-ci.org/mcuadros/go-stable.svg?branch=master)](https://travis-ci.org/mcuadros/go-stable) [![codecov.io](https://codecov.io/github/mcuadros/go-stable/coverage.svg?branch=master)](https://codecov.io/github/mcuadros/go-stable?branch=master) 


## Deploying your private go-stable

The easist way to deploy a private go-stable is using Docker.

```sh
docker run -it -v <certificates-folder>:/certificates -p :443:443 mcuadros/go-stable
```

*go-stable* runs allways under a TLS server so you need to provide trusted key and certificate to run it, place it in a the `<certificate-folder>` with the names `go-stable.cert.pem` and `go-stable.key.pem`

If you don't have a trusted certificate you can generate ones free and very fast ussing the [Let's Encrypt](https://letsencrypt.org/) tool. 

```sh
git clone https://github.com/letsencrypt/letsencrypt
cd letsencrypt
./letsencrypt-auto certonly --standalone -d <your-domain>
```
