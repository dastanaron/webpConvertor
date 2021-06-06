WebpConvertor
==========================

Convert image to WebP format.

How it working
--------------------------

This works with the [CWEBP](https://developers.google.com/speed/webp/docs/cwebp) library. The code calls methods of the external library to convert the image format. To increase the speed of work, the transfer of the binary content of files is not saved to disk, but is processed only using buffers and an I / O stream.

See https://storage.googleapis.com/downloads.webmproject.org/releases/webp/index.html


How it use
--------------------------

You can configure the server using docker as seen in the [docker-compose](./docker-compose.yml) example.
Or you can simply run the program with a configuration file similar to the following

```yaml
webpLibPath: /usr/local/bin #You location to cwebp library
port: 8080 #Port for listening request
```

and run 

```bash
. /path-to-application/webpConvertor -c /path-to-your-config

```

For example see [Dockerfile](./Dockerfile)

Useful Commands:
--------------------------

Compile programm

```bash
make compile
```