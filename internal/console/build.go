package console

import (
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

func Build(st *store.Store) *Server {
	reg := ns.NewRegistry()
	mainSec := reg.AddSection("主运一段", ns.KindMain)
	zone := reg.AddZone("主运区", mainSec.ID)
	line := belt.NewLine(reg)
	monitor := roller.NewMonitor(90, 3)
	ctrl := motor.NewController(line, monitor, st)
	brakes := brake.NewManager()
	ledger := quota.NewLedger()
	feed := coal.NewFeederManager(ledger)
	gasm := gas.NewManager()
	journal := audit.NewJournal(500, st)
	transfer := belt.NewTransferState()

	b1 := line.AddBelt("主运皮带", "main", "主运一段")
	_ = line.AddZone(b1.ID, zone.Name, []string{mainSec.ID})
	_ = ctrl.AddMotor(b1.ID, 250, 90)
	_ = brakes.AddBrake(b1.ID)
	_ = ledger.AddQuota(b1.ID, 5000, "shift")
	_ = feed.AddFeeder(b1.ID, 120)
	_ = monitor.AddRoller(b1.ID, 1)
	_ = monitor.AddRoller(b1.ID, 2)
	_ = transfer.AddPoint("tp1", "主运机头转载点", b1.ID, "洗选车间")
	g1 := gasm.AddSensor("综采面", 1.0)
	_ = gasm.SetConcentration(g1.ID, 0.4)

	brakes.SetReadyPort(ctrl)
	ctrl.SetFeederPort(feed)
	ctrl.SetBrakePort(brakes)
	feed.SetBeltReadyGate(ctrl)
	xlock := interlock.New(gasm, ctrl, ctrl, brakes)
	gasm.SetCalibrateListener(func(string) { xlock.SyncThresholds() })
	sc, _ := line.Scale(b1.ID)
	transit := coal.NewTransit(b1.ID, sc, 1)

	s := &Server{
		st:       st,
		reg:      reg,
		line:     line,
		ctrl:     ctrl,
		brakes:   brakes,
		feed:     feed,
		transit:  transit,
		transfer: transfer,
		gasm:     gasm,
		xlock:    xlock,
		monitor:  monitor,
		ledger:   ledger,
		journal:  journal,
	}
	s.restore(st)
	return s
}

func (s *Server) restore(st *store.Store) {
	if st == nil {
		return
	}
	quotas, err := st.LoadQuotaSnapshots()
	if err == nil {
		for _, row := range quotas {
			if q, ok := s.ledger.Quota(row.BeltID); ok {
				row.Apply(q)
			}
		}
	}
	rows, err := st.LoadBeltSnapshots()
	if err == nil {
		for _, row := range rows {
			if b, ok := s.line.Belt(row.ID); ok {
				row.Apply(b)
			}
		}
	}
	var events []audit.Event
	if err := st.LoadJSON("audit-events.json", &events); err == nil {
		s.journal.Load(events)
	}
}

func (s *Server) provision(b *belt.Belt) {
	_ = s.ctrl.AddMotor(b.ID, 250, 90)
	_ = s.brakes.AddBrake(b.ID)
	_ = s.ledger.AddQuota(b.ID, 5000, "shift")
	_ = s.feed.AddFeeder(b.ID, 100)
	_ = s.monitor.AddRoller(b.ID, 1)
	_ = s.monitor.AddRoller(b.ID, 2)
}
