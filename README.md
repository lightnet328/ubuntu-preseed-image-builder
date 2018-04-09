# ubuntu-preseed-image-builder

## Requirements

- Ubuntu 16.04.3 ISO Image

## Dependencies

- Docker

## Usage

### Creating a custom image

1. Clone this repository

```
$ git clone git@github.com:lightnet328/ubuntu-preseed-image-builder.git
```

2. Change current directory

```
$ cd ubuntu-preseed-image-builder
```

3. Download Ubuntu 16.04.3 ISO Image

```
$ wget http://old-releases.ubuntu.com/releases/xenial/ubuntu-16.04.3-server-amd64.iso
```

4. Copy env files

```
$ cp env.example.yml env.yml
$ cp env.secret.example.yml env.secret.yml
```

`env.secret.yml` overrides values in `env.yml`.

5. Build a image builder

```
$ docker build -t uibuild .
```

6. Build a custom image

```
$ docker run --rm --privileged -v $PWD:/builder uibuild build --suffix custom
```

The following is the same as above.

```
$ docker run --rm --privileged -v $PWD:/builder uibuild build --config env.yml --secret env.secret.yml --suffix custom
```

### Creating a live usb

#### macOS

2. Check the disk to write a custom image

```
$ diskutil list
```

3. Initialize $DISK and write a custom image

**Warning:**

The following command erases the contents of the disk. (In this case `/dev/disk2` will be erased)

Please check whether there is any mistake.

```
$ DISK=disk2
$ diskutil list
$ diskutil eraseDisk MS-DOS UNTITLED /dev/$DISK
$ diskutil unmountDisk /dev/$DISK
$ sudo dd if=ubuntu-16.04.3-server-amd64-custom.iso of=/dev/r$DISK bs=1m
```
