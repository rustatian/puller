package main

import (
	"flag"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"os"
	"strings"
	"sync"
)

func main() {
	fs := flag.NewFlagSet("dockerMegaUtility", flag.ExitOnError)
	nrun := fs.Bool("run", true, "run all containers")
	nstop := fs.Bool("stop", true, "stop all containers")
	ndelete := fs.Bool("delete", true, "delete all containers")


	err := fs.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}

	// Args stores filter arguments as map key:{map key: bool}.
	// It contains an aggregation of the map of arguments (which are in the form
	// of -f 'key=value') based on the key, and stores values for the same key
	// in a map with string keys and boolean values.
	// e.g given -f 'label=label1=1' -f 'label=label2=2' -f 'image.name=ubuntu'
	// the args will be {"image.name":{"ubuntu":true},"label":{"label1=1":true,"label2=2":true}}
	args := filters.NewArgs()

	for i := 0; i < len(fs.Args()); i ++ {
		splitted := strings.Split(fs.Args()[i], "=")
		args.Add("name", splitted[1])
	}


	if *nstop {
		stop(ctx, cli, args)
	}

	if *ndelete {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All:true, Filters:args})
		if err != nil {
			panic(err)
		}

		for _, cont := range containers {
			remove(ctx, cli, cont.ID)
		}
	}

	wg := &sync.WaitGroup{}
	if *nrun {
		wg.Add(6)
		go rabbit(ctx, cli, wg)
		go zipkin(ctx, cli, wg)
		go consul(ctx, cli, wg)
		go postgresql(ctx, cli, wg)
		go redis(ctx, cli, wg)
		go elastic(ctx, cli, wg)
	}
	wg.Wait()
}
