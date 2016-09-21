package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	context "golang.org/x/net/context"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetSourcesIDUsersUserIDExplorationsHandlerFunc turns a function with the right signature into a get sources ID users user ID explorations handler
type GetSourcesIDUsersUserIDExplorationsHandlerFunc func(context.Context, GetSourcesIDUsersUserIDExplorationsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetSourcesIDUsersUserIDExplorationsHandlerFunc) Handle(ctx context.Context, params GetSourcesIDUsersUserIDExplorationsParams) middleware.Responder {
	return fn(ctx, params)
}

// GetSourcesIDUsersUserIDExplorationsHandler interface for that can handle valid get sources ID users user ID explorations params
type GetSourcesIDUsersUserIDExplorationsHandler interface {
	Handle(context.Context, GetSourcesIDUsersUserIDExplorationsParams) middleware.Responder
}

// NewGetSourcesIDUsersUserIDExplorations creates a new http.Handler for the get sources ID users user ID explorations operation
func NewGetSourcesIDUsersUserIDExplorations(ctx *middleware.Context, handler GetSourcesIDUsersUserIDExplorationsHandler) *GetSourcesIDUsersUserIDExplorations {
	return &GetSourcesIDUsersUserIDExplorations{Context: ctx, Handler: handler}
}

/*GetSourcesIDUsersUserIDExplorations swagger:route GET /sources/{id}/users/{user_id}/explorations getSourcesIdUsersUserIdExplorations

GetSourcesIDUsersUserIDExplorations get sources ID users user ID explorations API

*/
type GetSourcesIDUsersUserIDExplorations struct {
	Context *middleware.Context
	Handler GetSourcesIDUsersUserIDExplorationsHandler
}

func (o *GetSourcesIDUsersUserIDExplorations) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetSourcesIDUsersUserIDExplorationsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(context.Background(), Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}