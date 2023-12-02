package net

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
)

var (
	connInst *netconn.NetConn
)

func Init() {
	c := config.Get()
	connInst = netconn.New(netconn.Config{
		StreamAddr:     c.Net.StreamAddr,
		ClientID:       c.Net.ClientID,
		ClientKey:      c.Net.ClientKey,
		ProgramVersion: config.ProgramVersion,
		Insecure:       c.Net.Insecure,
	})
}

func GetInst() *netconn.NetConn {
	if connInst == nil {
		Init()
	}

	return connInst
}
