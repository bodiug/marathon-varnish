package main

import (
"encoding/json"
"log"
"os"
"strings"
"text/template"
)

type Task struct {
    AppId        string `json:"appId"`
    Id           string `json:"id"`
    Host         string `json:"host"`
    Ports        []int  `json:"ports"`
    StartedAt    string `json:"startedAt"`
    StagedAt     string `json:"stagedAt"`
    Version      string `json:"version"`
    ServicePorts []int
}

func (t Task) Backend() string {
    return strings.Replace(strings.Replace(t.Id, ".", "_", -1), "-", "_", -1)
}

func (t Task) Director() string {
    return t.AppId
}

func (t Task) FirstPort() int {
    return t.Ports[0]
}

func (t Task) DirectorId() string {
    return strings.TrimLeft(t.AppId, "/")
}

type TaskList struct {
    Tasks []Task `json:"tasks"`
}

type Data struct {
    Directors map[string][]Task
    Backends  []Task
}

func (d *Data) Init(tasks []Task) {
    d.Directors = map[string][]Task{}
    for _, each := range tasks {
        ts, ok := d.Directors[each.DirectorId()]
        if !ok {
            ts = []Task{}
        }
        ts = append(ts, each)
        d.Directors[each.DirectorId()] = ts
    }
    d.Backends = tasks
}

var configTemplate = `
{{range .Backends}}
backend {{.Backend}} {
  .host = "{{.Host}}";
  .port = "{{.FirstPort}}";
  .probe = { .url = "/"; .interval = 5s; .timeout = 1s; .window = 5; .threshold = 3; }
}
{{end}}

{{range $k, $v := .Directors}}
director {{$k}} round-robin {   
{{range $v}}
  { .backend = {{.Backend}}; }
{{end}}
}
{{end}}

sub vcl_error {
  # Restart request flow on status 503
  if (obj.status == 503 && req.restarts < 4) {
    return (restart);
  }
}

sub vcl_recv {
{{range $k, $v := .Directors}}  
  if (req.http.host == "{{$k}}") {
    set req.backend = {{$k}};
    return (pass);
  }
{{end}}
  error 405;
}
`

var config = template.Must(template.New("config").Parse(configTemplate))

func main() {
    list := TaskList{}
    //strings.NewReader(test)
    if err := json.NewDecoder(os.Stdin).Decode(&list); err != nil {
        log.Fatal(err)
    }
    data := Data{}
    data.Init(list.Tasks)
    if err := config.Execute(os.Stdout, data); err != nil {
        log.Fatal(err)
    }
}

