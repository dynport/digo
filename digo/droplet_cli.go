package main

import (
	"fmt"
	"github.com/dynport/digo"
	"github.com/dynport/gocli"
	"os"
	"strconv"
	"time"
)

const RENAME_USAGE =  "<droplet_id> <new_name>"

func init() {
	cli.Register("droplet/rename",
		&gocli.Action{
			Handler:     RenameDropletAction,
			Description: "Describe Droplet",
			Usage: RENAME_USAGE,
		},
	)
}

func RenameDropletAction(args *gocli.Args) error {
	if len(args.Args) != 2 {
		fmt.Errorf(RENAME_USAGE)
	}
	id, newName := args.Args[0], args.Args[1]
	i, e := strconv.Atoi(id)
	if e != nil {
		return e
	}
	logger.Infof("renaming droplet %d to %s", i, newName)
	_, e = CurrentAccount().RenameDroplet(i, newName)
	if e != nil {
		return e
	}
	logger.Infof("renamed droplet %d to %s", i, newName)
	return nil
}

func init() {
	cli.Register("droplet/info",
		&gocli.Action{
			Handler:     DescribeDropletAction,
			Description: "Describe Droplet",
		},
	)
}

func DescribeDropletAction(args *gocli.Args) error {
	if len(args.Args) != 1 {
		return fmt.Errorf("USAGE: <droplet_id>")
	}
	i, e := strconv.Atoi(args.Args[0])
	if e != nil {
		return e
	}
	droplet, e := CurrentAccount().GetDroplet(i)
	if e != nil {
		return e
	}
	table := gocli.NewTable()
	table.Add("Id", fmt.Sprintf("%d", droplet.Id))
	table.Add("Name", droplet.Name)
	table.Add("Status", droplet.Status)
	table.Add("Locked", strconv.FormatBool(droplet.Locked))
	fmt.Println(table)
	return nil
}

func init() {
	cli.Register(
		"droplet/list",
		&gocli.Action{
			Handler:     ListDropletsAction,
			Description: "List active droplets",
		},
	)
}

func ListDropletsAction(args *gocli.Args) (e error) {
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
			Handler:     CreateDropletAction,
			Args:        args,
		},
	)
}

func CreateDropletAction(a *gocli.Args) error {
	started := time.Now()
	logger.Debugf("would create a new droplet with %#v", a.Args)
	if len(a.Args) != 1 {
		return fmt.Errorf("USAGE: create droplet <name>")
	}
	droplet := &digo.Droplet{Name: a.Args[0]}

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
	e = digo.WaitForDroplet(droplet)
	logger.Infof("droplet %d ready, ip: %s. total_time: %.1fs", droplet.Id, droplet.IpAddress, time.Now().Sub(started).Seconds())
	return e
}

func init() {
	cli.Register(
		"droplet/destroy",
		&gocli.Action{
			Description: "Destroy droplet",
			Handler:     DestroyDropletAction,
			Usage:       "<droplet_id>",
		},
	)
}

func DestroyDropletAction(args *gocli.Args) error {
	logger.Debugf("would destroy droplet with %#v", args)
	if len(args.Args) == 0 {
		return fmt.Errorf("USAGE: droplet destroy id1,id2,id3")
	}
	for _, id := range args.Args {
		if i, e := strconv.Atoi(id); e == nil {
			logger.Prefix = fmt.Sprintf("droplet-%d", i)
			droplet, e := CurrentAccount().GetDroplet(i)
			if e != nil {
				logger.Errorf("unable to get droplet for %d", i)
				continue
			}
			logger.Infof("destroying droplet %d", droplet.Id)
			rsp, e := CurrentAccount().DestroyDroplet(droplet.Id)
			if e != nil {
				return e
			}
			logger.Debugf("got response %+v", rsp)
			started := time.Now()
			archived := false
			for i := 0; i < 300; i++ {
				droplet.Reload()
				if droplet.Status == "archive" || droplet.Status == "off" {
					archived = true
					break
				}
				logger.Debug("status " + droplet.Status)
				fmt.Print(".")
				time.Sleep(1 * time.Second)
			}
			fmt.Print("\n")
			logger.Info("droplet destroyed")
			if !archived {
				logger.Errorf("error archiving %d", droplet.Id)
			} else {
				logger.Debugf("archived in %.06f", time.Now().Sub(started).Seconds())
			}
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
			Handler:     RebuildDropletAction,
			Usage:       "<droplet_id>",
			Args:        args,
		},
	)
}

func RebuildDropletAction(a *gocli.Args) error {
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
	droplet := &digo.Droplet{Id: i, Account: account}
	return digo.WaitForDroplet(droplet)
}
