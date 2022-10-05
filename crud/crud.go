package crud

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model interface{}

type Controller[T Model] gorm.DB

func New[T Model](orm *gorm.DB) *Controller[T] {
	return (*Controller[T])(orm)
}

func PrimaryKeyExpression(pk interface{}) clause.Expression {
	return clause.IN{
		Column: clause.Column{
			Table: clause.CurrentTable,
			Name:  clause.PrimaryKey,
			Raw:   false,
		},
		Values: []interface{}{pk},
	}
}

func getGuid(ctx *fasthttp.RequestCtx) string {
	id, ok := ctx.UserValue("guid").(string)
	if !ok {
		HandleError(ctx, http.StatusNotAcceptable, errors.New("guid is not a string"))
	}

	return id
}

func (c *Controller[T]) Get(ctx *fasthttp.RequestCtx) {
	var m T

	id := getGuid(ctx)
	err := (*gorm.DB)(c).Model(m).First(&m, PrimaryKeyExpression(id)).Error

	HandleError(ctx, http.StatusInternalServerError, err)
	Respond(ctx, m)
}

func (c *Controller[T]) List(ctx *fasthttp.RequestCtx) {
	var m T
	var slice []T
	err := (*gorm.DB)(c).Model(m).Find(&slice).Error
	HandleError(ctx, http.StatusInternalServerError, err)
	Respond(ctx, slice)
}

func (c *Controller[T]) Partial(ctx *fasthttp.RequestCtx) {
	id := getGuid(ctx)

	var m T
	db := (*gorm.DB)(c).Model(m)
	HandleError(ctx, http.StatusNotFound, db.First(&m, PrimaryKeyExpression(id)).Error)
	HandleError(ctx, http.StatusInternalServerError, json.Unmarshal(ctx.Request.Body(), &m))
	HandleError(ctx, http.StatusInternalServerError, (*gorm.DB)(c).Model(m).Select("*").Updates(m).Error)
	Respond(ctx, m)
}

func Respond(ctx *fasthttp.RequestCtx, m interface{}) {
	ctx.Response.Header.Add("content-type", "application/json")
	HandleError(ctx, http.StatusInternalServerError, json.NewEncoder(ctx).Encode(m))
}

func HandleError(ctx *fasthttp.RequestCtx, status int, err error) {
	if err == nil {
		return
	}

	ctx.SetStatusCode(status)
	lg := LoggerFromRequest(ctx)
	lg.Err(err).Int("status", status).Msg("found error")
	_, err = ctx.Write([]byte(err.Error()))

	// Stops execution
	panic(err)
}

func LoggerFromRequest(ctx *fasthttp.RequestCtx) zerolog.Logger {
	return log.With().Bytes("URL", ctx.URI().FullURI()).Bytes("method", ctx.Method()).Logger()
}
