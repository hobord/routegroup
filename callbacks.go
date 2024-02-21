package routegroup

import "log/slog"

// make log for each route registered
func RegisterRouteCallback(pattern string) {
	slog.Info("adding http route", "pattern", pattern)
}

// make log for each panic while registering a route
func RegisterPanicHandler(pattern string, err interface{}) {
	slog.Error("panic", "error", err, "pattern", pattern)
}
