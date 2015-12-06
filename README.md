# gop.kg [![Build Status](https://travis-ci.org/mcuadros/gop.kg.svg?branch=master)](https://travis-ci.org/mcuadros/gop.kg) [![codecov.io](https://codecov.io/github/mcuadros/gop.kg/coverage.svg?branch=master)](https://codecov.io/github/mcuadros/gop.kg?branch=master) 


## Deploying your private gop.kg

The easist way to deploy a private gop.kg is using Docker.

```sh
docker run -it -v <certificates-folder>:/certificates -p :443:443 mcuadros/gop.kg
```

*gop.kg* runs allways under a TLS server so you need to provide trusted key and certificate to run it, place it in a the `<certificate-folder>` with the names `gop.kg.cert.pem` and `gop.kg.key.pem`

If you don't have a trusted certificate you can generate ones free and very fast ussing the [Let's Encrypt](https://letsencrypt.org/) tool. 

```sh
git clone https://github.com/letsencrypt/letsencrypt
cd letsencrypt
./letsencrypt-auto certonly --standalone -d <your-domain>
```

> You must execute this commands where `<your-domain>` is hosted and the port 80 is open
