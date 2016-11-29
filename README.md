# go-stable [![Build Status](https://travis-ci.org/mcuadros/go-stable.svg?branch=master)](https://travis-ci.org/mcuadros/go-stable) [![codecov.io](https://codecov.io/github/mcuadros/go-stable/coverage.svg?branch=master)](https://codecov.io/github/mcuadros/go-stable?branch=master) [![GitHub release](https://img.shields.io/github/release/mcuadros/go-stable.svg)](https://github.com/mcuadros/go-stable) [![Docker Stars](https://img.shields.io/docker/pulls/mcuadros/go-stable.svg)](https://hub.docker.com/r/mcuadros/go-stable/tags/)


*go-stable* is a **self-hosted** service, that provides **versioned URLs** for any **Go package**, allowing to have fixed versions for your dependencies. *go-stable* is heavily inspired by [gopkg.in](http://labix.org/gopkg.in).

_How it works?_ Is a **proxy** between a git server, such as `github.com` or `bitbucket.com` and your git client, base on the requested URL (eg.: `example.com/repository.v1`) the most suitable available **tag is match** and used as **default branch**. 

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

**v0** contains more magic than expected, if none tag nor branch match, the *master* branch is returned.

## <a name="private" /> Using go-stable with private repositories

*go-stable* supports private repositories, since is based on HTTP protocol. The auth is done by [basic access authentication](https://en.wikipedia.org/wiki/Basic_access_authentication). 

That means that a _user_ and _password_ should be provided, you can use the GitHub user and password (if you are not using 2FA) but we really **recommend** using a [*personal token*](https://help.github.com/articles/creating-an-access-token-for-command-line-use/) that should be used as user. This methods allows you to not disclose your password, and invalidate the token whenever you want, revoking the access to the private repos.

When a package is installed through `go get` using a URL provided by *go-stable*, the terminal prompts are disabled. To provide the token (or the user and password...), you need to use a [`.netrc`](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) file. 

The line to add to ` ~/.netrc` (Linux or MacOS) or `%HOME%/_netrc` (Windows) is:
```
machine <your-go-stable-domain> login <your-personal-token>
```

If you are a bit paranoid, you can [encrypt](http://bryanwweber.com/writing/personal/2016/01/01/how-to-set-up-an-encrypted-.netrc-file-with-gpg-for-github-2fa-access/) your token or password using GPG.

## <a name="url" /> URL configuration

The URL router is based on [`gorilla/mux`](https://github.com/gorilla/mux), this enables `go-stable` with extremely flexible URLs patterns. By default, a couple of routes are configured, depending on the different flags provided.

### Same git provider and user/organization (recommended setup)

If all the packages are owned by the **same developer or organization** using the same provider (like *github.com*), you can specify the values for bot using the `--server` (for the provider part of the URL) and `--organization` flags.

In this case, the pattern *go-stable* uses is `/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}` (eg.: `example.com/repository.v1`). If we used the flag `--server github.com` and `--organization mcuadros`, the previous example will look for a version matching `v1` in the repo `github.com/mcuadros/repository`.

### If you need several organizations ... 

Leaving empty or not passing an `--organization` value will require the user to add a first segment in the URL to specify the developer or organization.

The pattern used in this case is `/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}`. If the flag `--server` was configured as `github.com`, `example.com/serabe/repository.v1` would look for a version matching `v1` in `github.com/serabe`.

### Multiple providers

Optionally, you can leave the `--server` empty too. In this case, a new segment would be needed at the beginning of the path specifying the provider. `example.com/github.com/mcuadros/go-stable.v1` will look for a version of `github.com/mcuadros/go-stable`.

The pattern used is `/{srv:[a-z0-9-.]+}/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}`.

### DIY

You can use any other pattern as long as you provide the router four variables: `srv`, `org`, `repository` and `version`. This feature is configured via the `--base-route` flag and the format for the pattern is specified by [`gorilla/mux`](https://github.com/gorilla/mux).


License
-------

MIT, see [LICENSE](LICENSE)
