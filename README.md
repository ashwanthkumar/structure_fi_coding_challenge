# structure_fi_coding_challenge

## Usage
This requires Go 1.17 toolchain to be installed on the machine.

```
$ make && ./structure_fi_coding_challenge
```

In case you have docker and not have Go toolchain installed, you can also use Docker to build and run the project using the following commands

```
# Build the docker image
$ docker build . -t structure_fi_by_ashwanth_kumar
$ docker run -p 8080:8080 -it structure_fi_by_ashwanth_kumar
```

This should start the application on port `8080` as default and we should be consuming the streams on the background.

Once the service starts, you can visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) for a swagger interface to test out the API implementation.