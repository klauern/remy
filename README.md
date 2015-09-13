# remy
CLI for interacting with the WebLogic Server RESTful Management Extensions.

[![wercker status](https://app.wercker.com/status/fe26b6defa4f97dde747ae1f1fbcb815/m "wercker status")](https://app.wercker.com/project/bykey/fe26b6defa4f97dde747ae1f1fbcb815)

`remy` is an application I wrote to learn a little about [golang](http://www.golang.org) as well as to itch a little
scratch I had with WebLogic servers.  WebLogic has a feature called *RESTful Management Extensions*, which lets you query
an AdminServer's domain for various pieces of information it knows about Servers, Clusters, Datasources, and Applications
that it is resposnible for.  This is pretty useful if you are interested in writing a lot of small scripts to quickly get
the status of some datasource in your cluster without having to resort to connecting, setting up, and maintaining your own
copy of a WebLogic Scripting Tool (WLST) interface.

This application provides a fancy command-line interface over top of this REST API, allowing you to quickly ascertain
how a server or cluster is doing without having to maintain all the complexity of an entire WLS domain locally just
to `connect(url="t3://server:7001", username="weblogic", password="welcome1")`, etc., etc.

# Installation

```sh
go get github.com/klauern/remy
```

# Usage

Running `remy help` will show you your options:

```
C:\> remy help
Query a WebLogic Domain's resources, including Datasources, Applications, Clusters, and Servers by using the WebLogic RESTful Management Extensions API

Usage:
  remy [command]

Available Commands:
  applications Query applications deployed under AdminServer
  config       Configure the credentials and server to default REST connections to
  clusters     Query clusters under AdminServer
  datasources  Query datasources under AdminServer
  servers      Display Server information
  version      Show the version of this command

Flags:
  -s, --adminurl="http://localhost:7001": Url for the Admin Server
  -f, --full-format[=false]: Return full format from REST server
  -p, --password="welcome1": Password for the user
  -u, --username="weblogic": Username with privileges to access AdminServer

Use "remy [command] --help" for more information about a command.
```

## Configuration and Authentication

Every command to query for information must inclue a set of credentials to

* The AdminServer the WebLogic Domain resides on
* Username
* Password

These can be provided by a number of options:

1. Command-line flags:
   * `--adminurl="http://server:7001"`
   * `--username="weblogic"`
   * `--password="welcome1"`
2. Environment Variables:
   * `export WLS_ADMINURL='http://server:7001'`
   * `export WLS_USERNAME='weblogic'`
   * `export WLS_PASSWORD='welcome1'`
3. Local directory `wlsrest.toml` config file
```sh
[/home/user/] $ cat ~/.wlsrest.toml
AdminURL = "http://homeserver:7001"
Username = "homeuser"
Password = "homepassword"
```
4. Home directory (~/, $HOME) `.wlsrest.toml` config file

### Generating Configuration for the above

Both the local directory and Home (`~/`) directory config files can be generated for you with `remy config`:

```
$ remy config -h
Configure what Username, Password, and Admin Server:Port you want to send REST requests to when submitting calls on any of the other commands

Usage:
  remy config [flags]

Flags:
      --environment[=false]: Set the WLS_* environment variables
      --home[=false]: Generate/Update the ~/$HOME config file
      --local[=false]: Generate/Update the local directory's config file

Global Flags:
  -s, --adminurl="http://localhost:7001": Url for the Admin Server
  -f, --full-format[=false]: Return full format from REST server
  -p, --password="welcome1": Password for the user
  -u, --username="weblogic": Username with privileges to access AdminServer
```

Using it is pretty straightforward:

```
$ remy config --local --adminurl="http://localserver:7001" --username="weblogic" --password="welcome1"
Using the Local directory to set the ./wlsrest.toml file
$ cat wlsrest.toml
AdminURL = "http://localserver:7001"
Username = "weblogic"
Password = "welcome1"
```

# Examples

Below are sample outputs provided by 

## Servers

TODO

## Clusters

TODO

## Datasources

TODO

## Applications

TODO

# TODO List Prior to `0.1`, `1.0`, `whatever.0`

[ ] - Pretty print formatting for responses
[ ] - Enrich documentation across the board
[ ] - Configure downloadable releases

# Contributing

Pull requests are welcome.  If you find this useful, please share and share alike.

# Contact

I can be reached on Twitter [@klauern](https://twitter.com/klauern) as well as on this repo.

