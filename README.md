# go-stable [![Build Status](https://travis-ci.org/mcuadros/go-stable.svg?branch=master)](https://travis-ci.org/mcuadros/go-stable) [![codecov.io](https://codecov.io/github/mcuadros/go-stable/coverage.svg?branch=master)](https://codecov.io/github/mcuadros/go-stable?branch=master) [![GitHub release](https://img.shields.io/github/release/mcuadros/go-stable.svg)](https://github.com/mcuadros/go-stable) [![Docker Stars](https://img.shields.io/docker/pulls/mcuadros/go-stable.svg)](https://hub.docker.com/r/mcuadros/go-stable/tags/)


*go-stable* is a **self-hosted** service, that provides **versioned URLs** for any **Go package**, allowing to have a stable APIs. *go-stable* is heavily inspired by [gopkg.in](http://labix.org/gopkg.in).

_How it works?_ Is a proxy between a git server, such as `github.com` or `bitbucket.com` and your git client, base on the requested URL (eg.: `example.com/repository.v1`) the most suitable available tag is match. 

The key features are:
- [Self-hosted service](#self-hosted)
- [Semantic Versioning](#semantic)
- [Private repositories support](#private)
- [Custom URLs](#url)



## <a name="self-hosted" /> Deploying your private go-stable

*go-stable* is a self-hosted server so **you need your own domain** and server to run the service.

The easiest way to deploy a *go-stable* is using Docker. Every time a new version is released our *CI* builds a new docker image, you can check all the available releases at [Docker Hub](https://hub.docker.com/r/mcuadros/go-stable/tags/) 

This is the bare minimum command to run a `go-stable` server:

```sh
docker run -d \
    -v <certificates-folder>:/certificates -p 443
    mcuadros/go-stable:<version> \
    stable server --host <my-domain>
```

Just run `stable server --help` to read all the available configuration options.

Since *go-stable* runs always under a TLS server, a trusted key and certificate is required to run it. 

By default a new certificate is issued to your domain at [*Let's Encrypt*](https://letsencrypt.org/) using [acmewrapper](https://github.com/dkumor/acmewrapper). In order to perform the domain validation a `go-stable server` running in the port `443` is required. After the first execution you can use another port, but this is not recommended, because the _auto-renovation_ happens every two weeks.

If you want to use a custom TLS key/certificate pair, maybe because your are in a private network or because you have already a valid certificated, you can place the files at the `<certificate-folder>` with the names `cert.pem` and `key.pem`

## <a name="semantic" /> Semantic Versioning
_Semantic Versioning_ is fully supported. The `version` variable from a URL as *example.com/org/repository*.**v1** is translated to a [`go-version`](https://github.com/mcuadros/go-version) constrain, like `v1.*`

This means that for example: if the repository has the following tags: `1.0`, `1.0rc1`. `1.10-dev`, the `1.0` will be chosen, but if the tags were: `1.0rc1`. `1.10-dev` the result is `1.0rc1`.

*go-stable* is not very strict with the tag format, you can use `v1.0` or just `1.0`.

## <a name="private" /> Using go-stable with private repositories

*go-stable* supports private repositories, since is based on HTTP protocol, the auth is done by [basic access authentication](https://en.wikipedia.org/wiki/Basic_access_authentication). 

That means that a _user_ and _password_ should be provided, you can use the GitHub user and password (if you are not using 2FA) but we really **recommend** use a [*personal token*](https://help.github.com/articles/creating-an-access-token-for-command-line-use/), the token should be used as user, this allows to you to not disclosure your password, and invalidate the token whenever you want.

When a package using *go-stable* is installed through `go get` the terminal prompts are disabled, to provide the the token (or the user and password...), you can do it with a [`.netrc`](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) file. 

You need to add to ` ~/.netrc` (Linux or MacOS) or `%HOME%/_netrc` (Windows), the following line:
```
machine <your-go-stable-domain> login <your-personal-token>
```

If you are a bit paranoid, you can [encrypt](http://bryanwweber.com/writing/personal/2016/01/01/how-to-set-up-an-encrypted-.netrc-file-with-gpg-for-github-2fa-access/) your token or password using GPG.

## <a name="url" /> URL configuration

The URL router is based on [`gorilla/mux`](https://github.com/gorilla/mux), this enables `go-stable` with extremely flexible URLs patterns, by default a couple of routes are configured, depending on the different flags provided.

#### Same git provider and user/organization (most interesting setup)
When all the packages are owned by the **same developer or organization**, in the same provider (eg.: *github.com*), you can configure the default ones using the `--server` and `--organization` flags, in this case the following pattern is used: `/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}` (eg.: `example.com/repository.v1`)

#### If you need several organizations ... 
Just leave empty the `--organization` flag, this configures: `/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}` (eg.: `example.com/org/repository.v1`).

#### Multiple providers
If you want to using different providers, I don't know why but... you can set `--server` to an empty value, then the pattern used is: `/{srv:[a-z0-9-.]+}/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}` (eg.: `example.com/github.com/org/repository.v1`)

You can use any pattern, the only requirements are four variables: *srv*, *org*, *repository* and *version*. You wrote your own and configure it using `--base-route` flag.



License
-------

MIT, see [LICENSE](LICENSE)