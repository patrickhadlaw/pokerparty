# Pokerparty

## Prerequisites

* Must have golang compiler installed
* Must have nodejs installed

## License

### This project is licensed under the license provided - see the [LICENSE](LICENSE) file for details

## Authors

* **Patrick Hadlaw** - [patrickhadlaw](https://github.com/patrickhadlaw)
* **Jack Capombassis** - [Jack-Capombassis](https://github.com/Jack-Capombassis)

## Build instructions

Clone repo:
```
$ git clone https://github.com/patrickhadlaw/pokerparty.git
```

To generate frontend:
```
$ cd pokerparty/frontend
$ ng build --prod
$ cp -R dist ../com/
```

To build pokerparty-server:
```
$ git clone https://bitbucket.org/patrickhadlaw/pokerparty.git
$ cd pokerparty
$ go build notion-server.go
```

To run:
```
$ ./pokerparty-server
```

## Commandline arguments

Run `$ ./pokerparty-server --help` for help