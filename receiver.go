package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouteType int32

const (
	RouteGet RouteType = iota
	RoutePost
)

type senderFunc func(gc *gin.Context, pool *pgxpool.Pool)

type Receiver struct {
	route        string
	authRequired bool
	routeType    RouteType
	sender       senderFunc
}
