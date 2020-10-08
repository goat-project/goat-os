# Goat-os

Openstack client for [goat](https://github.com/goat-project/goat) - Go Accounting Tool.

The Goat-os client is a command-line tool that connects to an Openstack cloud, 
extracts data about _virtual machines/servers_, _virtual networks_, _users_ and _images_, filters them 
accordingly and then sends them to a [goat server](https://github.com/goat-project/goat) for 
further processing. 

The data are filtered by time. Filter allows the records **from time**, **to time** or 
**for a period**. It cannot filter the records from/to and records for a period in 
the same time. Time from and time to can be used independently. Time from has to be 
earlier than time to.

See [goat wiki](https://github.com/goat-project/goat/wiki) for more info.

## Requirements
* Go 1.12 or newer
* Openstack instance
* [Goat server](https://github.com/goat-project/goat)

## Installation
The recommended way to install this tool is using `go get`:
```
go get -u github.com/goat-project/goat-os
```

## Configuration
Usage of goat-os:
```
Usage:
  goat-os [flags]
  goat-os [command]

Available Commands:
  help        Help about any command
  network     Extract network data
  storage     Extract storage data
  vm          Extract virtual machine data

Flags:
  -d, --debug string                 debug
  -e, --endpoint string              goat server [GOAT_SERVER_ENDPOINT] (required)
  -h, --help                         help for goat-os
  -i, --identifier string            goat identifier [IDENTIFIER] (required)
      --log-path string              path to log file
  -o, --openstack-endpoint string    Openstack endpoint [OPENSTACK_ENDPOINT] (required)
  -s, --openstack-secret string      Openstack secret [OPENSTACK_SECRET] (required)
  -p, --records-for-period string    records for period [TIME PERIOD]
  -f, --records-from string          records from [TIME]
  -t, --records-to string            records to [TIME]
      --version                      version for goat-os

Use "goat-os [command] --help" for more information about a command.
```

## Example
Extract virtual machine data from the last 5 years and save it with the identifier 'goat-vm'.
```
go run goat-os.go vm -p 5y -i goat-vm
```

## Container
The goat should run into the container described in [Dockerfile](https://github.com/goat-project/goat-os/blob/master/Dockerfile). 
Build and run commands:
```
docker image build -t goat-os-image .
docker run --rm -it --network host --name goat-os --volume goat-os:/var/goat-os goat-os-image
```

## Contributing
1. Fork [goat-os](https://github.com/goat-project/goat-os/fork)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request