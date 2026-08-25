package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"minebelt/internal/audit"
	"minebelt/internal/belt"
	"minebelt/internal/brake"
	"minebelt/internal/coal"
	"minebelt/internal/gas"
	"minebelt/internal/interlock"
	"minebelt/internal/motor"
	"minebelt/internal/ns"
	"minebelt/internal/quota"
	"minebelt/internal/roller"
	"minebelt/internal/store"
)

type Server struct {
	st       *store.Store
	reg      *ns.Registry
	line     *belt.Line
	ctrl     *motor.Controller
	brakes   *brake.Manager
	feed     *coal.FeederManager
	transit  *coal.Transit
	transfer *belt.TransferState
	gasm     *gas.Manager
	xlock    *interlock.Interlock
	monitor  *roller.Monitor
	ledger   *quota.Ledger
	journal  *audit.Journal
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", s.handleHealth)
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/belts", s.handleListBelts)
		api.Post("/belts", s.handleCreateBelt)
		api.Get("/belts/{beltID}", s.handleBeltStatus)
		api.Post("/belts/{beltID}/start", s.handleStart)
		api.Post("/belts/{beltID}/stop", s.handleStop)
		api.Post("/belts/{beltID}/extend", s.handleExtend)
		api.Post("/belts/{beltID}/speed", s.handleSetSpeed)
		api.Post("/belts/{beltID}/load", s.handleSetLoad)
		api.Post("/belts/{beltID}/scale", s.handleSetScale)
		api.Post("/belts/{beltID}/roller-sample", s.handleBeltRollerSample)
		api.Post("/rollers/{rollerID}/sample", s.handleRollerSample)
		api.Get("/rollers/{rollerID}", s.handleRollerStatus)
		api.Get("/rollers", s.handleListRollers)
		api.Post("/rollers/reset-filter", s.handleResetRollerFilter)
		api.Post("/flow/tick", s.handleFlowTick)
		api.Get("/flow", s.handleFlowStatus)
		api.Post("/flow/reset", s.handleFlowReset)
		api.Post("/interlock/cascade", s.handleCascade)
		api.Post("/gas/{face}/calibrate", s.handleCalibrateGas)
		api.Post("/gas/{face}/sample", s.handleGasSample)
		api.Get("/gas/{face}/verdict", s.handleGasVerdict)
		api.Get("/gas/alarms", s.handleGasAlarms)
		api.Post("/maintenance/{beltID}/tag", s.handleTag)
		api.Post("/maintenance/{beltID}/clear", s.handleClearTag)
		api.Post("/maintenance/{beltID}/release", s.handleRelease)
		api.Post("/motors/{beltID}/replace", s.handleReplaceMotor)
		api.Get("/motors/{beltID}", s.handleMotorStatus)
		api.Post("/brakes/{beltID}/trip", s.handleBrakeTrip)
		api.Post("/brakes/{beltID}/reset", s.handleBrakeReset)
		api.Get("/brakes/{beltID}", s.handleBrakeStatus)
		api.Post("/brakes/{beltID}/pressure", s.handleSetPressure)
		api.Get("/quota", s.handleQuota)
		api.Post("/quota/{beltID}/reset", s.handleQuotaReset)
		api.Post("/quota/{beltID}/limit", s.handleQuotaLimit)
		api.Post("/quota/{beltID}/rollover", s.handleQuotaRollover)
		api.Get("/quota/usage", s.handleQuotaUsage)
		api.Get("/quota/snapshots", s.handleQuotaSnapshotLoad)
		api.Post("/quota/snapshots", s.handleQuotaSnapshotSave)
		api.Get("/transfer", s.handleTransferList)
		api.Post("/transfer/{pointID}/block", s.handleTransferBlock)
		api.Get("/audit", s.handleAudit)
		api.Get("/snapshots", s.handleSnapshotLoad)
		api.Post("/snapshots", s.handleSnapshotSave)
	})
	return r
}
