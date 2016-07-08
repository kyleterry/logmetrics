# Metrics From Nginx Logs

A small program to read from the nginx access log every 5 seconds and format the
results for use in Statsd.

Results will look similar to the following:

```
20x:10|s
30x:45|s
40x:15|s
50x:20|s
/route1:10s
/route2:10s
```

The above means that during the 5 second interval, 10 "200" status requests came
in, 45 "300" status requests came in, 15 "400" status requests came in, and 20
"500" status requests came in; 10 of which were from route1 and 10 were from route2.

## Running and testing in Docker

There are a couple Dockerfiles to get things up and running.

### Building log metrics

`docker build -t kyleterry/logmetrics .`

### Building the test nginx image

`cd nginx`  
`docker build -t kyleterry/nginx .`

### Running both

To run a test scenario, run the nginx container. /var/log/nginx is a
volume:

`docker run -p 8080:80 --name ng -d kyleterry/nginx`

Then run the logmetrics container and mount the ng volume:

`docker run -d --name logmetrics --volumes-from ng kyleterry/logmetrics`

The log metrics default input file is `/var/log/nginx/access.log`.

If you want to view the requests being formatted and written, mount the volume
in an ubuntu or debian container and tail the log:

`docker run -i -t --volumes-from logmetrics ubuntu /bin/bash`
`tail -f /var/log/stats.log`

Now on your Docker host, you can hit the nginx container via your browser or cli
using curl, wget, or ab:

`curl localhost:8080`

I setup an `/error_page` location that returns a 500 error for testing that too.
