package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func stop(ctx context.Context, cli *client.Client, f filters.Args) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: f})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}
