WebpConvertor
==========================

Convert image to WebP format.

See [Swagger Documentation](https://app.swaggerhub.com/apis-docs/dastanaron/webpConvertor/1.0.2)

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
mode: ram #or tmp
port: 8080 #Port for listening request
```

If you set `ram` mode, picture will be download to RAM, converted and writen to response from RAM.

If you set `tmp` mode, picture will be downloaded to system temp, converted to new file, writen to response and deleted.
This mode loads less RAM

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

Project development
--------------------------

The project will develop according to my free time and desire to develop it. If you want to stimulate my desire to develop the project, you can buy me a coffee or [donate me](https://my.qiwi.com/Dmytryi-RPq_-kS82_)