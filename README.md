# env [![GoDoc](https://godoc.org/github.com/sajari/env?status.svg)](http://godoc.org/github.com/sajari/env)
Environment variable management

Package Goals:
* Up-front check for existance *and* validity of *all* required config variables
* Can output *all* variables used by a process on startup
* Can export and import full env configurations
* Seamlessly transition local testing environments to [Kubernetes](https://kubernetes.io/)

Inspired by the simplicity and power of the standard library `flag` package, `env` works by first defining a set of variables in your program which will take the value of ENV variables. 

This is extended with the [`envsvc`](http://godoc.org/pkg/github.com/sajari/env/envsvc) package which is explicitly designed for use in service processes. 

Example below:

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/sajari/env"
	"github.com/sajari/env/envsvc"
)

func main() {
	// Define variables...
	listen := env.BindAddr("LISTEN", "bind address for gRPC server")
	debug := env.BindAddr("LISTEN_DEBUG", "bind address for http debug server")
	apiKey := env.String("API_KEY", "key id for api authorisation")
	workers := env.Int("WORKERS", "number of parallel workers to start")

	// Parse them from the environment and exit if all env vars are not set
	envsvc.Parse()

	// Debug server (optional)
	http.ListenAndServe(*debug, nil)

	// Your code begins here!
}
```

By default `String`, `Int`, `Bool`, `Duration` are defined.  We also have flag types `BindAddr` and `DialAddr` typically included in all of our services for exposing ports on server processes and dialing to other servers.

### Extending this pattern
It’s also possible to define new variable types by implementing the Value interface (which is the same as in `flag`). You can also define separate sets of variables, rather than using the global functions in the `env` package.  

### Config export, generation and more
Inheriting a project and getting configured is simple. The above code example will exit if the env vars are not set. From an engineer perspective this is great, you immediately see why the service won't start and what you need to fix to get running.

The built-in `--help` flag tells you your options: 

```shell
$ ./my-service --help
Usage of ./my-service:
  -env-check
    	check env variables
  -env-dump
    	dump env variables
  -env-dump-json
    	dump env variables in JSON format
  -env-dump-yaml
    	dump env variables in YAML format
```

So we can `check` which env vars are required:

```shell
$ ./my-service -env-check
missing env MY_SERVICE_LISTEN
missing env MY_SERVICE_LISTEN_DEBUG
missing env MY_SERVICE_API_KEY
missing env MY_SERVICE_WORKERS
```

Note: the env vars are prefixed with the service name to avoid clashes. 

Ok that's useful, now we know what we need to get this service up and running. I'm lazy, so i want this done for me:

```shell
$ ./my-service -env-dump
# bind address for gRPC server
export MY_SERVICE_LISTEN=""       

# bind address for http debug server
export MY_SERVICE_LISTEN_DEBUG="" 

# key id for api authorisation
export MY_SERVICE_API_KEY=""      

# number of parallel workers to start
export MY_SERVICE_WORKERS=""      
```

Excellent, now I have a workable set of environment parameters ready to be exported, note that they are currently blank. At this point the engineer has some options:

* Create workable values (particularly if the owner)
* Look for example values in the service repo (good practice)
* Ask the service owner

So after making these up and adding them to my environment i can re-run the `dump` commmand to get the following:

```shell
$ ./my-service -env-dump
export MY_SERVICE_LISTEN=":1234"         # bind address for gRPC server
export MY_SERVICE_LISTEN_DEBUG=":5678"   # bind address for http debug server
export MY_SERVICE_API_KEY="abc"          # key id for api authorisation
export MY_SERVICE_WORKERS="4"            # number of parallel workers to start
```

We now have a fully working environment that will be validated on service start and can be exported and shared with other engineers as needed. 
