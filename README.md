Κάκα
====

CLI for managing [Cacus](https://github.com/beebeeep/cacus) repos.

Installation:
-------------
Assuming you have Go installed on your machine, do:
```
go get github.com/beebeeep/caca
```

I'll share compiled binaries for all platforms a bit later.

Usage:
------
First, you will need config file, `~/.cacarc` (or at other relevant path for your OS):
```yaml
instances:
    my-distro:
        base_url: https://cacus.example.org/debian
        token: <token generated using "cacus --gen-token">
    other-distro:
        base_url: https://cacus.example.org/debian
        token: <token generated using "cacus --gen-token">
```
Once you have it, you can try:
1) Upload the packages:
  ```sh
  $ caca upload debs/*.deb my-distro/testing
  Uploading debs/example_1.0_all.deb... SUCCESS: Package example_1.0 was uploaded to my_distro/testing
  Uploading debs/other_1.0_all.deb... SUCCESS: Package other_1.0 was uploaded to my_distro/testing
  ```
2) Other stuff TBD
