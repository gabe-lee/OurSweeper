package server_client_manager

import (
	"net"
	"sync"

	"github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/internal/user_token"
	"github.com/gabe-lee/OurSweeper/lock"
	"github.com/google/uuid"
)

type (
	AnonToken   = user_token.UserStats
	Conn        = net.Conn
	UUID        = uuid.UUID
	Mutex       = sync.Mutex
	ServerWorld = common.ServerWorld
	MiniLock    = lock.MiniLock

	Lock = *MiniLock
)

type Client struct {
	Stats       AnonToken
	IsAnon      bool
	LastMsgTime int64
}

type ClientManager struct {
	clientMapLock Lock
	clientMap     map[UUID]Client
}

func (c *ClientManager) UserIsLoggedOn(u UUID) bool {
	c.clientMapLock.Lock()
	defer c.clientMapLock.Unlock()
	_, loggedOn := c.clientMap[u]
	return loggedOn
}

func (c *ClientManager) LoginUser(u UUID) bool {
	c.clientMapLock.Lock()
	defer c.clientMapLock.Unlock()
	_, loggedOn := c.clientMap[u]
	return loggedOn
}
