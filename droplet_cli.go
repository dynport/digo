package main

import (
	"fmt"
	"github.com/dynport/gocli"
	"os"
	"strconv"
)

func init() {
	args := &gocli.Args{}
	args.RegisterInt("-i", false, CurrentAccount().ImageId, "Image id for new droplet")
	args.RegisterInt("-r", false, CurrentAccount().RegionId, "Region id for new droplet")
	args.RegisterInt("-s", false, CurrentAccount().SizeId, "Size id for new droplet")
	args.RegisterInt("-k", false, CurrentAccount().SshKey, "Ssh key to be used")

	cli.Register(
		"droplet/create",
		&gocli.Action{
			Description: "Create new droplet",
			Usage:       "<name>",
			Handler:     CreateDroplet,
			Args:        args,
		},
	)
}

func init() {
	cli.Register(
		"droplet/list",
		&gocli.Action{
			Handler:     ListDroplets,
			Description: "List active droplets",
		},
	)
}

func ListDroplets(args *gocli.Args) (e error) {
	logger.Debug("listing droplets")

	droplets, e := CurrentAccount().Droplets()
	if e != nil {
		return e
	}

	if _, e := CurrentAccount().CachedSizes(); e != nil {
		return e
	}

	table := gocli.NewTable()
	if len(droplets) == 0 {
		table.Add("no droplets found")
	} else {
		table.Add("Id", "Created", "Status", "Locked", "Name", "IPAddress", "Region", "Size", "Image")
		for _, droplet := range droplets {
			table.Add(
				strconv.Itoa(droplet.Id),
				droplet.CreatedAt.Format("2006-01-02T15:04"),
				droplet.Status,
				strconv.FormatBool(droplet.Locked),
				droplet.Name,
				droplet.IpAddress,
				fmt.Sprintf("%s (%d)", CurrentAccount().RegionName(droplet.RegionId), droplet.RegionId),
				fmt.Sprintf("%s (%d)", CurrentAccount().SizeName(droplet.SizeId), droplet.SizeId),
				fmt.Sprintf("%s (%d)", CurrentAccount().ImageName(droplet.ImageId), droplet.ImageId),
			)
		}
	}
	fmt.Fprintln(os.Stdout, table.String())
	return nil
}

func CreateDroplet(a *gocli.Args) error {
	logger.Debugf("would create a new droplet with %#v", a.Args)
	if len(a.Args) != 1 {
		return fmt.Errorf("USAGE: create droplet <name>")
	}
	droplet := &Droplet{Name: a.Args[0]}

	var e error
	if droplet.SizeId, e = a.GetInt("-s"); e != nil {
		return e
	}

	if droplet.ImageId, e = a.GetInt("-i"); e != nil {
		return e
	}

	if droplet.RegionId, e = a.GetInt("-r"); e != nil {
		return e
	}

	if droplet.SshKey, e = a.GetInt("-k"); e != nil {
		return e
	}

	droplet, e = CurrentAccount().CreateDroplet(droplet)
	if e != nil {
		return e
	}
	droplet.Account = CurrentAccount()
	logger.Infof("created droplet with id %d", droplet.Id)
	return WaitForDroplet(droplet)
}

func init() {
	cli.Register(
		"droplet/destroy",
		&gocli.Action{
			Description: "Destroy droplet",
			Handler:     DestroyDroplet,
			Usage:       "<droplet_id>",
		},
	)
}

func DestroyDroplet(args *gocli.Args) error {
	logger.Debugf("would destroy droplet with %#v", args)
	if len(args.Args) == 0 {
		return fmt.Errorf("USAGE: droplet destroy id1,id2,id3")
	}
	for _, id := range args.Args {
		if i, e := strconv.Atoi(id); e == nil {
			logger.Infof("destroying droplet %d", i)
			rsp, e := CurrentAccount().DestroyDroplet(i)
			if e != nil {
				return e
			}
			logger.Debugf("got response %+v", rsp)
		}
	}
	return nil
}

func init() {
	args := &gocli.Args{}
	args.RegisterInt("-i", false, 0, "Rebuild droplet")
	cli.Register(
		"droplet/rebuild",
		&gocli.Action{
			Description: "Rebuild droplet",
			Handler:     RebuildDroplet,
			Usage:       "<droplet_id>",
			Args:        args,
		},
	)
}

func RebuildDroplet(a *gocli.Args) error {
	if len(a.Args) != 1 {
		return fmt.Errorf("USAGE: droplet rebuild <id>")
	}
	i, e := strconv.Atoi(a.Args[0])
	if e != nil {
		return fmt.Errorf("USAGE: droplet rebuild <id>")
	}

	imageId, e := a.GetInt("-i")
	if e != nil {
		return e
	}

	rsp, e := account.RebuildDroplet(i, imageId)
	if e != nil {
		return e
	}
	logger.Debugf("got response %+v", rsp)
	droplet := &Droplet{Id: i, Account: account}
	return WaitForDroplet(droplet)
}
