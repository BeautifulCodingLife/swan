package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes(mux *mux.Router) {
	routes := []*Route{
		NewRoute("GET", "/v1/apps", s.listApps),
		NewRoute("POST", "/v1/apps", s.createApp),
		NewRoute("GET", "/v1/apps/{app_id}", s.getApp),
		NewRoute("DELETE", "/v1/apps/{app_id}", s.deleteApp),
		NewRoute("POST", "/v1/apps/{app_id}/scale", s.scaleApp),
		NewRoute("PUT", "/v1/apps/{app_id}", s.updateApp),
		NewRoute("POST", "/v1/apps/{app_id}/rollback", s.rollback),
		NewRoute("PATCH", "/v1/apps/{app_id}/weights", s.updateWeights),

		NewRoute("GET", "/v1/apps/{app_id}/tasks", s.getTasks),
		NewRoute("GET", "/v1/apps/{app_id}/tasks/{task_id}", s.getTask),
		NewRoute("DELETE", "/v1/apps/{app_id}/tasks/{task_id}", s.deleteTask),
		NewRoute("DELETE", "/v1/apps/{app_id}/tasks", s.deleteTasks),
		NewRoute("PUT", "/v1/apps/{app_id}/tasks/{task_id}", s.updateTask),
		NewRoute("POST", "/v1/apps/{app_id}/tasks/{task_id}", s.rollbackTask),
		NewRoute("PATCH", "/v1/apps/{app_id}/tasks/{task_id}/weight", s.updateWeight),

		NewRoute("GET", "/v1/apps/{app_id}/versions", s.getVersions),
		NewRoute("GET", "/v1/apps/{app_id}/versions/{version_id}", s.getVersion),
		NewRoute("POST", "/v1/apps/{app_id}/versions", s.createVersion),

		NewRoute("POST", "/v1/compose", s.newCompose),
		NewRoute("POST", "/v1/compose/parse", s.parseYAML),
		NewRoute("GET", "/v1/compose", s.listComposes),
		NewRoute("GET", "/v1/compose/{compose_id}", s.getCompose),
		NewRoute("DELETE", "/v1/compose/{compose_id}", s.deleteCompose),

		NewRoute("GET", "/ping", s.ping),
		NewRoute("GET", "/v1/events", s.events),
		NewRoute("GET", "/v1/stats", s.stats),
		NewRoute("GET", "/version", s.version),
		NewRoute("GET", "/v1/leader", s.getLeader),
		NewRoute("POST", "/v1/purge", s.purge),

		NewRoute("GET", "/v1/debug/dump", s.dump),
		NewRoute("GET", "/v1/debug/load", s.load),
		NewRoute("GET", "/v1/fullsync", s.fullEventsAndRecords),

		NewRoute("GET", "/v1/agents", s.listAgents),
		NewRoute("GET", "/v1/agents/{agent_id}", s.getAgent),
		NewPrefixRoute("ANY", "/v1/agents/{agent_id}/proxy", s.redirectAgentProxy),
		NewPrefixRoute("ANY", "/v1/agents/{agent_id}/dns", s.redirectAgentDNS),
		NewPrefixRoute("ANY", "/v1/agents/{agent_id}/docker", s.redirectAgentDocker),
	}

	log.Debug("Registering HTTP route")

	for _, r := range routes {
		var (
			handler = s.makeHTTPHandler(r.Handler())
			path    = r.Path()
			methods = r.Methods()
		)

		log.Debugf("Registering %v, %s", methods, path)

		if r.prefix {
			mux.PathPrefix(path).Methods(methods...).Handler(handler)
		} else {
			mux.Path(path).Methods(methods...).Handler(handler)
		}
	}
}
