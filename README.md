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
    my-cacus:
        default: true
        base_url: https://cacus.example.org/debian
        token: <token generated using "cacus --gen-token">
        ca_cert: /etc/ssl/certs/ca-certificates.crt
    other-cacus:
        base_url: https://cacus.example.org/debian
        token: <token generated using "cacus --gen-token">
        ca_cert: /etc/ssl/certs/my-ca.crt
```
Once you have it, you can try:
1) Upload the packages:
  ```sh
  $ caca upload debs/*.deb my-distro/testing            # upload to default cacus instance
  Uploading debs/example_1.0_all.deb... SUCCESS: Package example_1.0 was uploaded to my-distro/testing
  Uploading debs/other_1.0_all.deb... SUCCESS: Package other_1.0 was uploaded to my-distro/testing
  ...
  $ caca -instance other-cacus upload debs/*.deb other-distro/testing
  Uploading debs/example_1.0_all.deb... SUCCESS: Package example_1.0 was uploaded to other-distro/testing
  Uploading debs/other_1.0_all.deb... SUCCESS: Package other_1.0 was uploaded to other-distro/testing
  ```

2) Show information about distro:
  ```sh
  $ caca show test-vcad
  Distro 'test-vcad':
        Name: test-vcad
        Description: test vcad distro
        Components: [unstable testing stable]
        Number of packages: 17
        Type: general
        Origin: N/A
        Last updated at: 2017-06-28T11:53:10.983000 
   ```

3) Search packages:
  ```sh
  $ caca  search -distro cacus-jessie -pkg 'python-cacus$' -ver "0.7-\d"
  ==== Results for distro cacus-jessie ====
        Package: python-cacus
        Version: 0.7-1
        Maintainer: Danila Migalin <me@miga.me.uk>
        Architecure: all
        Components: [unstable]
        Description: Distributed Debian repository manager
   ```
4) Copy packages between components:
  ```sh
  $ caca copy -distro common -pkg python-cacus -ver 0.7.12 -from unstable -to stable 
  Package 'python-cacus_0.7.12' was copied in distro 'common' from 'unstable' to 'stable'
  ```
