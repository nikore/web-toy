# web-toy
Web Toy application.  Frequently used by DevOps to test deployment pipelines and otherwise cause mayhem.

## The rules:

1. No PR review needed, go ahead and hack
2. Don't remove existing functionality
3. If you break it, fix it (this applies to deployment as well as the code itself)

## Useful paths on port 8080:

* `/`—hello, world!
* `/disconnect`—disconnects the TCP/IP connection before any response is sent
* `/hello`—ditto
* `/env`—prints environment variables (except a blacklist)
* `/health`—health check
* `/status`—loadgen status
* `/timeout`—never returns

## Other functionality

* Junk data (currently the output of /usr/games/fortune) is served on port 1332.

## Environment variables:

* `WEB_TOY_LOAD_CPU`—If set and not `0`, generate 100% CPU load for the lifetime of the process
* `WEB_TOY_LOAD_MEM_MB`—If set to a valid non-0 integer, allocate that many megabytes of memory.  Will touch the memory once a minute to ensure it's not paged out.  (Make sure you raise the task memory allocation in Marathon to match, or else it'll get killed repeatedly.)
* `FIVEOHTHREE_FUSE_TIMER`—If set and not `0`, make the `/health` check return "503 Fuse timer exceeded" after this many seconds
* `HTTP_PORT`—Start webserver on this port (8080 default)
* `JUNK_PORT`—Serve non-HTTP data on this port (outputs `/usr/games/fortune`)
