# plotjitter
Generate side-by-side jitter [boxplots](https://en.wikipedia.org/wiki/Box_plot) from one or more packet capture files.

## Usage example
```sh
go get github.com/slavc/plotjitter
~/go/bin/plotjitter -bpf tcp aaa.pcap bbb.pcap
```
A `jitter.png` file will get generated in current working directory, with a side-by-side boxplots of jitter from aaa.pcap and bbb.pcap only of TCP traffic.
