# remy
CLI for interacting with the WebLogic Server RESTful Management Services.

[![wercker status](https://app.wercker.com/status/fe26b6defa4f97dde747ae1f1fbcb815/m "wercker status")](https://app.wercker.com/project/bykey/fe26b6defa4f97dde747ae1f1fbcb815)

`remy` is an application I wrote to learn a little about [golang](http://www.golang.org) as well as to itch a little
scratch I had with WebLogic servers.  WebLogic has a feature called *[RESTful Management Extensions](http://docs.oracle.com/cd/E23943_01/web.1111/e24682/toc.htm#RESTS149)*, which lets you query
an AdminServer's domain for various pieces of information it knows about Servers, Clusters, Datasources, and Applications
that it is resposnible for.  This is pretty useful if you are interested in writing a lot of small scripts to quickly get
the status of some datasource in your cluster without having to resort to connecting, setting up, and maintaining your own
copy of a WebLogic Scripting Tool (WLST) interface.

This application provides a fancy command-line interface over top of this REST API, allowing you to quickly ascertain
how a server or cluster is doing without having to maintain all the complexity of an entire WLS domain locally just
to `connect(url="t3://server:7001", username="weblogic", password="welcome1")`, etc., etc.

# Installation

## Releases

The latest version I put out there is `v0.2.1`, and that can be found on the releases page: [v0.2.1 release](https://github.com/klauern/remy/releases/tag/v0.2.1)

## From Source

For the die-hard coder in you:

```sh
go get -u github.com/klauern/remy
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
4. Home directory (~/, $HOME) `.wlsrest.toml` config file

An example `wlsrest.toml` or `.wlsrest.toml` config file:

```
[/home/user/] $ cat ~/.wlsrest.toml
AdminURL = "http://homeserver:7001"
Username = "homeuser"
Password = "homepassword"
```

### Generating Configuration for the above

Both the local directory and Home (`~/`) directory config files can be generated for you with `remy config`.  This
provides the added benefit of encrypting the password for you:

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
Password = "{AES}VHQVFwN72jWgRYzWbnJQugUfCa6LAU0W"
```

NOTE: Do not attempt to assume there's some inherent security in this.  The password is generated with AES, but the default
`key` is stored in code.  You can optionally override the default key with an environment variable `WLS_REMYKEY`:

```
$ export WLS_REMYKEY="My very very very awesome key!!!" # (MUST be 32 bytes in length EXACTLY)
$ remy config --local
$ cat .\wlsrest.toml
AdminURL = "http://localhost:7001"
Username = "weblogic"
Password = "{AES}BM1uj9uv1bD7KV6BXapCf1kucxDYbCU6"
```

# Query Examples

Below are sample outputs provided by the tool itself.  This is a rudimentary **1.0** of the output.  I hope to provide
 some options for different output at some point, such as rendering a widget for use with [`termui`](https://github.com/gizak/termui).
 Any help would be appreciated.

## Servers

### All Servers (short form)

```sh
$ remy servers
Finding all Servers
Using Full Format? false
Name:        AdminServer   | State:           RUNNING       | Health:        HEALTH_OK
Cluster:                   | CurrentMachine:                | JVM Load:      0
Sockets #:   0             | Heap Sz Cur:     0             | Heap Free Cur: 0
Java Ver:                  | OS Name:                       | OS Version:
WLS Version:

Name:        WLS_WSM1      | State:           RUNNING       | Health:        HEALTH_OK
Cluster:                   | CurrentMachine:                | JVM Load:      0
Sockets #:   0             | Heap Sz Cur:     0             | Heap Free Cur: 0
Java Ver:                  | OS Name:                       | OS Version:
WLS Version:

Name:        WLS_SOA1      | State:           RUNNING       | Health:        HEALTH_OK
Cluster:                   | CurrentMachine:                | JVM Load:      0
Sockets #:   0             | Heap Sz Cur:     0             | Heap Free Cur: 0
Java Ver:                  | OS Name:                       | OS Version:
WLS Version:

Name:        WLS_OSB1      | State:           RUNNING       | Health:        HEALTH_OK
Cluster:                   | CurrentMachine:                | JVM Load:      0
Sockets #:   0             | Heap Sz Cur:     0             | Heap Free Cur: 0
Java Ver:                  | OS Name:                       | OS Version:
WLS Version:
```

### Individual Server (always full-format)

```sh
$ remy servers AdminServer
Finding Server information for AdminServer
Server AdminServer:
Name:        AdminServer   | State:           RUNNING       | Health:        HEALTH_OK
Cluster:                   | CurrentMachine:  localhost     | JVM Load:      0
Sockets #:   8             | Heap Sz Cur:     3151495168    | Heap Free Cur: 423742888
Java Ver:    1.7.0_80      | OS Name:         Linux         | OS Version:    2.6......
WLS Version: WebLogic Server 10.3.6.0  Tue Nov 15 08:52:36 PST 2011 1441050
```

## Clusters

### All Clusters (short form)

```sh
$ remy clusters
Finding All Clusters
Using Full Format? false
Name: WSM-PM_Cluster
State:             RUNNING       | Health:               HEALTH_OK     | Cluster Master?       false         | Drop Out Freq:
Resend Req. Count: 0             | Fragments Sent Count: 0             | Fragments Recv Count: 0

Name: SOA_Cluster
State:             RUNNING       | Health:               HEALTH_OK     | Cluster Master?       false         | Drop Out Freq:
Resend Req. Count: 0             | Fragments Sent Count: 0             | Fragments Recv Count: 0

Name: OSB_Cluster
State:             RUNNING       | Health:               HEALTH_OK     | Cluster Master?       false         | Drop Out Freq:
Resend Req. Count: 0             | Fragments Sent Count: 0             | Fragments Recv Count: 0
```

### Individual Cluster (always full-format)

```sh
$ remy clusters SOA_Cluster
Finding Cluster information for SOA_Cluster
Name: SOA_Cluster
State:             RUNNING       | Health:               HEALTH_OK     | Cluster Master?       false         | Drop Out Freq:
Resend Req. Count: 0             | Fragments Sent Count: 127690        | Fragments Recv Count: 0
```


## Datasources

TODO

## Applications

### All Applications (short form)

```
$ remy applications
Finding All Applications
Using Full Format? false
Name: FileAdapter                                       |AppType: rar  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: DbAdapter                                         |AppType: rar  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: JmsAdapter                                        |AppType: rar  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: AqAdapter                                         |AppType: rar  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: FtpAdapter                                        |AppType: rar  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: SocketAdapter                                     |AppType: rar  |State: STATE_ACTIVE|Health: HEALTH_OK
...
Name: b2bui                                             |AppType: ear  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: Healthcare UI                                     |AppType: ear  |State:             |Health:
Name: DefaultToDoTaskFlow                               |AppType: ear  |State: STATE_ACTIVE|Health: HEALTH_OK
Name: composer                                          |AppType: ear  |State: STATE_ACTIVE|Health: HEALTH_OK
...
```

### Individual Application (always full-format)

In an application, we don't output sections where there is nothing, hence there are pieces missing from this particular
application.  You may see a lot more data or a lot less depending.

```
$ remy applications composer
Finding application information for composer
Name: composer                                          |AppType: ear  |State: STATE_ACTIVE|Health: HEALTH_OK
Target States
        Target: SOA_Cluster                             |State: STATE_ACTIVE
Work Managers
        Name: default                                   |Server: WLS_SOA1      |Pending Requests: 0             |Completed Requests: 0
        Name: wm/SOAWorkManager                         |Server: WLS_SOA1      |Pending Requests: 0             |Completed Requests: 0
```

# TODO List Prior to `1.0.0`

- [ ] TONS more tests (test-first is hard for me, sorry guys)
  - [ ] Command-Line Flag parsing Tests
  - [X] Configuration Parsing / Flag handling
- [ ] Pretty print formatting for responses
  - [ ] Possible `termui` implementation
- [ ] Enrich documentation across the board
- [X] Configure downloadable releases

# Contributing

Pull requests are welcome.  If you find this useful, please share and share alike.

# Contact

I can be reached on Twitter [@klauern](https://twitter.com/klauern) as well as on this repo.

