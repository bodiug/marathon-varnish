marathon-varnish
================

This is an example to create a Varnish configuration from the Mesos Marathon API.

The output from `http://localhost:8080/v2/tasks` will output something like:

```
{
  "tasks": [
    {
      "appId": "/foo",
      "id": "foo.4cef049d-643e-11e4-bdcd-9cb65491f714",
      "host": "10.98.214.12",
      "ports": [
        31001
      ],
      "startedAt": "2014-12-06T09:17:40.523Z",
      "stagedAt": "2014-11-04T16:19:06.948Z",
      "version": "2014-10-24T14:14:52.870Z",
      "servicePorts": [
        0
      ]
    },
    {
      "appId": "/foo",
      "id": "foo.49651a86-643e-11e4-bdcd-9cb65491f714",
      "host": "10.98.214.12",
      "ports": [
        31000
      ],
      "startedAt": "2014-12-06T09:17:40.531Z",
      "stagedAt": "2014-11-04T16:19:01.075Z",
      "version": "2014-10-24T14:14:52.870Z",
      "servicePorts": [
        0
      ]
    },
    {
      "appId": "/hello",
      "id": "hello.4cf0160f-643e-11e4-bdcd-9cb65491f714",
      "host": "10.98.214.13",
      "ports": [
        31001
      ],
      "startedAt": "2014-12-06T09:17:40.538Z",
      "stagedAt": "2014-11-04T16:19:06.954Z",
      "version": "2014-10-24T14:14:00.053Z",
      "servicePorts": [
        0
      ]
    },
    {
      "appId": "/hello",
      "id": "hello.4971c4ba-643e-11e4-bdcd-9cb65491f714",
      "host": "10.98.214.13",
      "ports": [
        31000
      ],
      "startedAt": "2014-12-06T09:17:40.544Z",
      "stagedAt": "2014-11-04T16:19:01.094Z",
      "version": "2014-10-24T14:14:00.053Z",
      "servicePorts": [
        0
      ]
    }
  ]
}
```

Using `curl -s http://localhost:8080/v2/tasks | ./marathon-varnish` this will generate a Varnish VCL for you:

```
backend foo_4cef049d_643e_11e4_bdcd_9cb65491f714 {
  .host = "10.98.214.12";
  .port = "31001";
  .probe = { .url = "/"; .interval = 5s; .timeout = 1s; .window = 5; .threshold = 3; }
}

backend foo_49651a86_643e_11e4_bdcd_9cb65491f714 {
  .host = "10.98.214.12";
  .port = "31000";
  .probe = { .url = "/"; .interval = 5s; .timeout = 1s; .window = 5; .threshold = 3; }
}

backend hello_4cf0160f_643e_11e4_bdcd_9cb65491f714 {
  .host = "10.98.214.13";
  .port = "31001";
  .probe = { .url = "/"; .interval = 5s; .timeout = 1s; .window = 5; .threshold = 3; }
}

backend hello_4971c4ba_643e_11e4_bdcd_9cb65491f714 {
  .host = "10.98.214.13";
  .port = "31000";
  .probe = { .url = "/"; .interval = 5s; .timeout = 1s; .window = 5; .threshold = 3; }
}

director foo round-robin {   
  { .backend = foo_4cef049d_643e_11e4_bdcd_9cb65491f714; }
  { .backend = foo_49651a86_643e_11e4_bdcd_9cb65491f714; }
}

director hello round-robin {   
  { .backend = hello_4cf0160f_643e_11e4_bdcd_9cb65491f714; }
  { .backend = hello_4971c4ba_643e_11e4_bdcd_9cb65491f714; }
}

sub vcl_error {
  # Restart request flow on status 503
  if (obj.status == 503 && req.restarts < 4) {
    return (restart);
  }
}

sub vcl_recv {
  if (req.http.host == "foo") {
    set req.backend = foo;
    return (pass);
  }
  
  if (req.http.host == "hello") {
    set req.backend = hello;
    return (pass);
  }

  error 405;
}
```
