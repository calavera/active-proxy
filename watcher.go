package main

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type watcher struct {
	etcdLeader   string
	etcdMachines []string
	client       *etcd.Client
}

func NewWatcher() *watcher {
	w := &watcher{}
	w.Init()
	return w
}

func (w *watcher) Init() {
	w.client = etcd.NewClient([]string{})

	if len(w.etcdMachines) > 0 {
		w.client.SetCluster(w.etcdMachines)
	}
}

func (w *watcher) StartApplications(p *proxy) {
	w.loadApplications(p)

	go func() {
		appsChannel := make(chan *etcd.Response, 10)
		w.watchApplications(p, appsChannel)

		w.client.Watch("applications", 0, appsChannel, nil)
	}()
}

func (w *watcher) loadApplications(p *proxy) {
	values, err := w.client.Get("applications")

	if err == nil {
		for _, entry := range values {
			app := strings.Split(entry.Key, "/")[2]
			w.registerApp(app, p)
		}
	}
}

func (w *watcher) watchApplications(p *proxy, appsChannel chan *etcd.Response) {
	for {
		entry := <-appsChannel
		app := strings.Split(entry.Key, "/")[2]

		w.registerApp(app, p)
	}
}

func (w *watcher) registerApp(app string, p *proxy) {
	values, err := w.client.Get("applications/" + app)

	if err != nil {
		log.Printf("Error getting settings for: %s\nReason: %s", app, err.Error())
	} else {
		a := &application{Name: app}

		for _, value := range values {
			switch value.Key {
			case "/applications/" + app + "/port":
				a.Port = value.Value
			case "/applications/" + app + "/test":
				a.Test = value.Value
			}
		}

		p.Route(a)
	}
}
