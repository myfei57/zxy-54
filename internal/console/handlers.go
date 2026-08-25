package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"minebelt/internal/audit"
	"minebelt/internal/belt"
	"minebelt/internal/ns"
	"minebelt/internal/store"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":            "ok",
		"sections":          s.reg.SectionCount(),
		"extended_sections": len(s.reg.SectionsByKind(ns.KindExtended)),
		"zones":             s.reg.ZoneCount(),
		"belt_count":        s.reg.BeltCount(),
		"running":           s.line.RunningCount(),
		"fault":             s.line.FaultCount(),
		"tagged":            s.xlock.TaggedCount(),
		"gas_alarms":        s.gasm.AlarmCount(),
		"audit_events":      s.journal.Total(),
	})
}

func (s *Server) handleListBelts(w http.ResponseWriter, r *http.Request) {
	belts := s.line.All()
	rows := make([]map[string]interface{}, 0, len(belts))
	for _, b := range belts {
		rows = append(rows, s.beltRow(b))
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"belts": rows})
}

func (s *Server) beltRow(b *belt.Belt) map[string]interface{} {
	sc, _ := s.line.Scale(b.ID)
	inst := 0.0
	if sc != nil {
		inst = sc.Read()
	}
	row := map[string]interface{}{
		"id":              b.ID,
		"name":            b.Name,
		"kind":            b.Kind,
		"section":         b.Section,
		"zone":            b.Zone,
		"status":          b.Status.String(),
		"speed":           b.Speed,
		"load":            b.Load,
		"scale_inst":      inst,
		"feeder_running":  s.feed.IsRunning(b.ID),
		"motor_powered":   s.ctrl.Powered(b.ID),
		"start_order":     s.feed.Order(b.ID),
		"feeder_attempts": s.feed.Attempts(b.ID),
		"feeder_rejected": s.feed.LastReject(b.ID) != nil,
		"ready_version":   s.ctrl.ReadyVersion(b.ID),
		"clearance_log":   s.ctrl.ClearanceLog(b.ID),
		"has_clearance":   s.ctrl.VerifyClearance(b.ID),
	}
	if f, ok := s.feed.Feeder(b.ID); ok {
		row["feeder_rate"] = f.Rate
	}
	if st, err := s.brakes.Status(b.ID); err == nil {
		row["brake_engaged"] = st.Engaged
		row["brake_overloaded"] = st.Overloaded
	}
	row["brake_latched"] = s.brakes.Latched(b.ID)
	row["brake_cause"] = s.brakes.Cause(b.ID)
	return row
}

func (s *Server) handleCreateBelt(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string  `json:"name"`
		Kind    string  `json:"kind"`
		Section string  `json:"section"`
		Zone    string  `json:"zone"`
		Rating  float64 `json:"rating_kw"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if body.Name == "" || body.Section == "" {
		writeError(w, http.StatusBadRequest, errNameOrSectionMissing)
		return
	}
	sec, ok := s.reg.SectionByName(body.Section)
	if !ok {
		sec = s.reg.AddSection(body.Section, nsKindFromString(body.Kind))
	}
	b := s.line.AddBelt(body.Name, body.Kind, body.Section)
	zoneName := body.Zone
	if zoneName == "" {
		zoneName = body.Name + "区"
	}
	_ = s.line.AddZone(b.ID, zoneName, []string{sec.ID})
	s.provision(b)
	s.journal.Append(audit.KindMaintenance, b.ID, "belt provisioned")
	writeJSON(w, http.StatusCreated, s.beltRow(b))
}

func nsKindFromString(kind string) ns.SectionKind {
	switch kind {
	case "transfer":
		return ns.KindTransfer
	case "extended":
		return ns.KindExtended
	}
	return ns.KindMain
}

func (s *Server) handleBeltStatus(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	b, ok := s.line.Belt(beltID)
	if !ok {
		writeError(w, http.StatusNotFound, errBeltMissing)
		return
	}
	row := s.beltRow(b)
	row["tagged"] = s.xlock.Tagged(beltID)
	route, err := s.line.PatrolRoute(beltID)
	if err == nil {
		row["patrol_route"] = route
	}
	var sections []string
	if zoneID := s.line.ZoneID(beltID); zoneID != "" {
		if z, ok := s.reg.Zone(zoneID); ok {
			for _, sid := range z.SectionIDs {
				if sec, ok := s.reg.Section(sid); ok {
					sections = append(sections, sec.Name)
				}
			}
		}
	}
	row["sections"] = sections
	writeJSON(w, http.StatusOK, row)
}

func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	if err := s.ctrl.StartOneKey(beltID); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.journal.Append(audit.KindStart, beltID, "one-key start")
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "started": true})
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	if err := s.ctrl.Stop(beltID); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.journal.Append(audit.KindStop, beltID, "manual stop")
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "stopped": true})
}

func (s *Server) handleExtend(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Section string `json:"section"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.line.Extend(beltID, body.Section); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	route, _ := s.line.PatrolRoute(beltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "patrol_route": route})
}

func (s *Server) handleSetSpeed(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Speed float64 `json:"speed"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.line.SetSpeed(beltID, body.Speed); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "speed": body.Speed})
}

func (s *Server) handleSetLoad(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Load float64 `json:"load"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.line.SetLoad(beltID, body.Load); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "load": body.Load})
}

func (s *Server) handleSetScale(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Instant float64 `json:"instant"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sc, ok := s.line.Scale(beltID)
	if !ok {
		writeError(w, http.StatusNotFound, errBeltMissing)
		return
	}
	sc.SetInstant(body.Instant)
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "scale_inst": body.Instant})
}

func (s *Server) handleBeltRollerSample(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Index int     `json:"index"`
		Temp  float64 `json:"temp"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	rw, ok := s.monitor.RollerByBelt(beltID, body.Index)
	if !ok {
		writeError(w, http.StatusNotFound, errRollerMissing)
		return
	}
	filtered, err := s.monitor.Sample(rw.ID, body.Temp)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	tripped := s.ctrl.Protect(beltID, body.Temp)
	if tripped {
		_ = s.line.SetStatus(beltID, belt.Fault)
		s.journal.Append(audit.KindInterlock, beltID, "roller temperature protection tripped")
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"roller_id": rw.ID,
		"filtered":  filtered,
		"tripped":   tripped,
	})
}

func (s *Server) handleRollerSample(w http.ResponseWriter, r *http.Request) {
	rollerID := chi.URLParam(r, "rollerID")
	var body struct {
		Temp float64 `json:"temp"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	filtered, err := s.monitor.Sample(rollerID, body.Temp)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	rw, _ := s.monitor.Roller(rollerID)
	tripped := false
	if rw != nil {
		tripped = s.ctrl.Protect(rw.BeltID, body.Temp)
		if tripped {
			_ = s.line.SetStatus(rw.BeltID, belt.Fault)
			s.journal.Append(audit.KindInterlock, rw.BeltID, "roller temperature protection tripped")
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"roller_id": rollerID,
		"filtered":  filtered,
		"tripped":   tripped,
	})
}

func (s *Server) handleRollerStatus(w http.ResponseWriter, r *http.Request) {
	rollerID := chi.URLParam(r, "rollerID")
	rw, ok := s.monitor.Roller(rollerID)
	if !ok {
		writeError(w, http.StatusNotFound, errRollerMissing)
		return
	}
	filtered, has := s.monitor.FilterValue(rw.BeltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"roller_id":  rw.ID,
		"belt_id":    rw.BeltID,
		"index":      rw.Index,
		"temp":       rw.Temp,
		"filtered":   filtered,
		"has_filter": has,
		"limit":      s.monitor.Limit(),
	})
}

func (s *Server) handleListRollers(w http.ResponseWriter, r *http.Request) {
	beltID := r.URL.Query().Get("belt_id")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"rollers": s.monitor.Rollers(beltID),
	})
}

func (s *Server) handleResetRollerFilter(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BeltID string `json:"belt_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.monitor.ResetFilter(body.BeltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": body.BeltID, "filter_reset": true})
}

func (s *Server) handleFlowTick(w http.ResponseWriter, r *http.Request) {
	s.transit.Tick()
	beltID := s.transit.BeltID()
	_ = s.feed.RecordProduction(beltID, s.transit.Inst())
	s.journal.Append(audit.KindQuota, beltID, "flow tick recorded")
	remaining, _ := s.ledger.Remaining(beltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"belt_id":   beltID,
		"inst":      s.transit.Inst(),
		"total":     s.transit.Total(),
		"remaining": remaining,
	})
}

func (s *Server) handleFlowStatus(w http.ResponseWriter, r *http.Request) {
	beltID := s.transit.BeltID()
	sc, _ := s.line.Scale(beltID)
	inst := 0.0
	if sc != nil {
		inst = sc.Read()
	}
	remaining, _ := s.ledger.Remaining(beltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"belt_id":         beltID,
		"scale_inst":      inst,
		"flow_inst":       s.transit.Inst(),
		"flow_total":      s.transit.Total(),
		"flow_peak":       s.transit.Peak(),
		"flow_average":    s.transit.Average(),
		"quota_remaining": remaining,
	})
}

func (s *Server) handleFlowReset(w http.ResponseWriter, r *http.Request) {
	s.transit.ResetTotal()
	s.journal.Append(audit.KindQuota, s.transit.BeltID(), "flow totals reset")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"belt_id": s.transit.BeltID(),
		"total":   s.transit.Total(),
	})
}

func (s *Server) handleCascade(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BeltIDs []string `json:"belt_ids"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	errs := s.xlock.Cascade(body.BeltIDs)
	for _, id := range body.BeltIDs {
		s.journal.Append(audit.KindStop, id, "cascade stop")
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"errors": errs,
		"count":  len(body.BeltIDs),
	})
}

func (s *Server) handleCalibrateGas(w http.ResponseWriter, r *http.Request) {
	face := chi.URLParam(r, "face")
	var body struct {
		Threshold float64 `json:"threshold"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sensor, ok := s.gasm.SensorByFace(face)
	if !ok {
		writeError(w, http.StatusNotFound, errSensorMissing)
		return
	}
	if err := s.gasm.Calibrate(sensor.ID, body.Threshold); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"face":      face,
		"threshold": body.Threshold,
		"version":   sensor.Version,
	})
}

func (s *Server) handleGasSample(w http.ResponseWriter, r *http.Request) {
	face := chi.URLParam(r, "face")
	var body struct {
		Concentration float64 `json:"concentration"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sensor, ok := s.gasm.SensorByFace(face)
	if !ok {
		writeError(w, http.StatusNotFound, errSensorMissing)
		return
	}
	if err := s.gasm.SetConcentration(sensor.ID, body.Concentration); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"face":          face,
		"concentration": body.Concentration,
	})
}

func (s *Server) handleGasVerdict(w http.ResponseWriter, r *http.Request) {
	face := chi.URLParam(r, "face")
	sensor, ok := s.gasm.SensorByFace(face)
	if !ok {
		writeError(w, http.StatusNotFound, errSensorMissing)
		return
	}
	verdict := s.xlock.GasVerdict(sensor.ID)
	if verdict {
		s.gasm.RecordAlarm(face, sensor.Concentration, sensor.Threshold)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"face":    face,
		"verdict": verdict,
	})
}

func (s *Server) handleGasAlarms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":  s.gasm.AlarmCount(),
		"alarms": s.gasm.Alarms(20),
	})
}

func (s *Server) handleTag(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	_ = s.xlock.Tag(beltID)
	s.journal.Append(audit.KindMaintenance, beltID, "tag placed")
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "tagged": true})
}

func (s *Server) handleClearTag(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	_ = s.xlock.ClearTag(beltID)
	s.journal.Append(audit.KindMaintenance, beltID, "tag removed")
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "tagged": false})
}

func (s *Server) handleRelease(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	if err := s.xlock.Release(beltID); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.journal.Append(audit.KindMaintenance, beltID, "maintenance release")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"belt_id": beltID,
		"powered": s.ctrl.Powered(beltID),
	})
}

func (s *Server) handleReplaceMotor(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		RatingKW float64 `json:"rating_kw"`
		Baseline float64 `json:"baseline"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	m, err := s.ctrl.ReplaceMotor(beltID, body.RatingKW, body.Baseline)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	s.journal.Append(audit.KindMaintenance, beltID, "motor replaced")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"motor_id":  m.ID,
		"belt_id":   beltID,
		"rating_kw": m.RatingKW,
		"baseline":  m.Baseline,
	})
}

func (s *Server) handleMotorStatus(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	m, ok := s.ctrl.MotorForBelt(beltID)
	if !ok {
		writeError(w, http.StatusNotFound, errMotorMissing)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"motor_id":  m.ID,
		"belt_id":   m.BeltID,
		"rating_kw": m.RatingKW,
		"baseline":  m.Baseline,
		"powered":   m.Powered,
		"temp_ok":   !s.ctrl.TempVerdict(beltID, m.Baseline),
	})
}

func (s *Server) handleBrakeTrip(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Cause string `json:"cause"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.brakes.Trip(beltID, body.Cause); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	s.journal.Append(audit.KindInterlock, beltID, "brake overload trip")
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "tripped": true})
}

func (s *Server) handleBrakeReset(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	if err := s.xlock.Reset(beltID); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.journal.Append(audit.KindInterlock, beltID, "brake latch reset")
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "reset": true})
}

func (s *Server) handleBrakeStatus(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	st, err := s.brakes.Status(beltID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	br, _ := s.brakes.Brake(beltID)
	id := ""
	if br != nil {
		id = br.ID
	}
	pressure, _ := s.brakes.Pressure(beltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"brake_id":        id,
		"belt_id":         st.BeltID,
		"engaged":         st.Engaged,
		"overloaded":      st.Overloaded,
		"cause":           st.Cause,
		"latch":           st.Latch,
		"pressure":        pressure.Bar,
		"pressure_ready":  pressure.Ready,
		"release_log":     s.brakes.ReleaseLog(beltID),
		"release_version": s.brakes.ReleaseVersion(beltID),
	})
}

func (s *Server) handleSetPressure(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Bar float64 `json:"bar"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	p, err := s.brakes.SetPressure(beltID, body.Bar)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"belt_id": beltID,
		"bar":     p.Bar,
		"ready":   p.Ready,
	})
}

func (s *Server) handleQuota(w http.ResponseWriter, r *http.Request) {
	beltID := r.URL.Query().Get("belt_id")
	q, ok := s.ledger.Quota(beltID)
	if !ok {
		writeError(w, http.StatusNotFound, errQuotaMissing)
		return
	}
	remaining, _ := s.ledger.Remaining(beltID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"belt_id":   q.BeltID,
		"limit":     q.Limit,
		"used":      q.Used,
		"remaining": remaining,
		"period":    q.Period,
		"exhausted": s.ledger.Exhausted(beltID),
	})
}

func (s *Server) handleQuotaReset(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	if err := s.ledger.ResetPeriod(beltID); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "reset": true})
}

func (s *Server) handleQuotaLimit(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Limit float64 `json:"limit"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.ledger.SetLimit(beltID, body.Limit); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "limit": body.Limit})
}

func (s *Server) handleQuotaRollover(w http.ResponseWriter, r *http.Request) {
	beltID := chi.URLParam(r, "beltID")
	var body struct {
		Period string `json:"period"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.ledger.Rollover(beltID, body.Period); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"belt_id": beltID, "period": body.Period})
}

func (s *Server) handleQuotaUsage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"quotas": s.ledger.Usage(),
	})
}

func (s *Server) handleQuotaSnapshotSave(w http.ResponseWriter, r *http.Request) {
	st := s.storeRef()
	if st == nil {
		writeError(w, http.StatusServiceUnavailable, errStoreMissing)
		return
	}
	rows := make([]store.QuotaSnapshot, 0, len(s.ledger.Usage()))
	for _, q := range s.ledger.Usage() {
		rows = append(rows, store.FromQuota(&q))
	}
	if err := st.SaveQuotaSnapshots(rows); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"saved": len(rows)})
}

func (s *Server) handleQuotaSnapshotLoad(w http.ResponseWriter, r *http.Request) {
	st := s.storeRef()
	if st == nil {
		writeError(w, http.StatusServiceUnavailable, errStoreMissing)
		return
	}
	rows, err := st.LoadQuotaSnapshots()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	for _, row := range rows {
		if q, ok := s.ledger.Quota(row.BeltID); ok {
			row.Apply(q)
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"loaded": len(rows)})
}

func (s *Server) handleTransferList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"points": s.transfer.All(),
	})
}

func (s *Server) handleTransferBlock(w http.ResponseWriter, r *http.Request) {
	pointID := chi.URLParam(r, "pointID")
	var body struct {
		Blocked bool `json:"blocked"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.transfer.SetBlocked(pointID, body.Blocked); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	point, _ := s.transfer.Point(pointID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"point_id": pointID,
		"name":     point.Name,
		"blocked":  point.Blocked,
	})
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query().Get("kind")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"counts": s.journal.Counts(),
		"events": s.journal.Recent(kind, 50),
	})
}

func (s *Server) handleSnapshotSave(w http.ResponseWriter, r *http.Request) {
	st := s.storeRef()
	if st == nil {
		writeError(w, http.StatusServiceUnavailable, errStoreMissing)
		return
	}
	rows := make([]store.BeltSnapshot, 0, len(s.line.All()))
	for _, b := range s.line.All() {
		rows = append(rows, store.FromBelt(b))
	}
	if err := st.SaveBeltSnapshots(rows); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"saved": len(rows)})
}

func (s *Server) handleSnapshotLoad(w http.ResponseWriter, r *http.Request) {
	st := s.storeRef()
	if st == nil {
		writeError(w, http.StatusServiceUnavailable, errStoreMissing)
		return
	}
	rows, err := st.LoadBeltSnapshots()
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	for _, row := range rows {
		if b, ok := s.line.Belt(row.ID); ok {
			row.Apply(b)
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"loaded": len(rows)})
}

func (s *Server) storeRef() *store.Store {
	return s.st
}
