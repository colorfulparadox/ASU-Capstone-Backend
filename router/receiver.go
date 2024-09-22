package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouteType uint16

const (
	RouteGet RouteType = iota
	RoutePost
)

type senderFunc func(gc *gin.Context, pool *pgxpool.Pool)
type middlewareFunc func(gc *gin.Context, pool *pgxpool.Pool)

type Receiver struct {
	Route      string
	RouteType  RouteType
	Middleware middlewareFunc
	Sender     senderFunc
}

func default_middleware(gc *gin.Context, pool *pgxpool.Pool) { gc.Next() }
