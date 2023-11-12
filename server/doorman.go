package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"golang.org/x/exp/slog"
)

type Doorman struct {
	*pb.UnimplementedDoormanServer

	conn *pgxpool.Pool

	processing chan bool

	sets    db.Sets
	changes db.Changes
	objects db.Objects
	roles   db.Roles
	tuples  db.Tuples
}

func (d *Doorman) ProcessAllChanges() error {
	for {
		if err := d.ProcessChange(); err != nil {
			if err == pgx.ErrNoRows {
				return nil
			}
			return fmt.Errorf("processing change failed: %w", err)
		}
	}
}

func (d *Doorman) Changes(ctx context.Context, request *pb.ChangesRequest) (*pb.ChangesResponse, error) {
	changes, err := d.changes.List(ctx, db.ChangeFilter{PaginationToken: request.PaginationToken})
	if err != nil {
		return nil, fmt.Errorf("db failed: %w", err)
	}
	res := &pb.ChangesResponse{
		Items: make([]*pb.Change, len(changes)),
	}
	for i, v := range changes {
		res.Items[i] = mapChangeToPb(v)
		if i == len(changes)-1 {
			res.PaginationToken = &v.ID
		}
	}
	return res, nil
}

func (d *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	success, err := d.sets.Contains(ctx, doorman.Set{
		Object: doorman.Object(request.Object),
		Verb:   doorman.Verb(request.Verb),
	}, doorman.Object(request.Subject))
	if err != nil {
		return &pb.CheckResponse{}, fmt.Errorf("check failed: %w", err)
	}
	return &pb.CheckResponse{Success: success}, nil
}

func (d *Doorman) Grant(ctx context.Context, request *pb.GrantRequest) (*pb.GrantResponse, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	res, err := d.grantWithTx(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	d.processChangesImmediately()

	return res, nil
}

func (d *Doorman) processChangesImmediately() {
	go func() {
		d.processing <- true
	}()
}

func (d *Doorman) ListObjects(ctx context.Context, request *pb.ListObjectsRequest) (*pb.ListObjectsResponse, error) {
	sub := doorman.Object(request.Subject)
	parents, err := d.sets.ListParents(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("sets.ListParents failed: %w", err)
	}

	items := make([]*pb.Relation, len(parents))
	for i, parent := range parents {
		items[i] = &pb.Relation{
			Subject: string(sub),
			Verb:    string(parent.Verb),
			Object:  string(parent.Object),
		}
	}

	return &pb.ListObjectsResponse{
		Items: items,
	}, nil
}

func (d *Doorman) ListRoles(ctx context.Context, request *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	roles, err := d.roles.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list failed: %w", err)
	}

	items := make([]*pb.Role, len(roles))
	for i, r := range roles {
		items[i] = mapRoleToPb(r)
	}

	return &pb.ListRolesResponse{
		Items: items,
	}, nil
}

func (d *Doorman) ProcessChange() error {
	timeout := time.Second * 5
	stalePeriod := time.Hour

	// This timeout should be higher than the "timeout", otherwise the tx.Commit will fail
	ctx, cancel := context.WithTimeout(context.Background(), timeout+time.Second*2)
	defer cancel()

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a tx: %w", err)
	}

	var c doorman.Change
	err = tx.QueryRow(ctx, `
		update changes
		set status='processed'
		where id in
		(
		  select id
		  from changes
			where status='pending'
		  order by random()
		  for update skip locked
		  limit 1
		)
		returning id, type, payload, created_at
	`).Scan(&c.ID, &c.Type, &c.Payload, &c.CreatedAt)

	// No rows = no tasks
	if err == pgx.ErrNoRows {
		slog.Debug("no tasks, sleeping")

		select {
		case <-d.processing:
		case <-time.After(timeout):
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("tx failed to commit: %w", err)
		}
		return pgx.ErrNoRows
	}

	// Failed to execute query, probably a bad query/schema
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("failed to rollback: %w", err)
		}
		return fmt.Errorf("failed to query/scan: %w", err)
	}

	// Process task
	if err := d.processChange(ctx, c); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("failed to rollback: %w", err)
		}

		time.Sleep(timeout)

		// Tasks older than stalePeriod get logged
		if time.Since(c.CreatedAt) > stalePeriod {
			// TODO: probably log this somewhere else
			slog.Info("stale change", "change", c)
		}

		return fmt.Errorf("failed to process change %s: %w", c.Type, err)
	}

	// No errors, so task can be deleted
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx failed to commit: %w", err)
	}
	return nil
}

func (d *Doorman) RebuildCache(ctx context.Context, request *pb.RebuildCacheRequest) (*pb.RebuildCacheResponse, error) {
	if err := d.changes.SetStatusOfAll(ctx, "pending"); err != nil {
		return nil, fmt.Errorf("marking all as pending failed: %w", err)
	}
	return &pb.RebuildCacheResponse{}, nil
}

func (d *Doorman) RemoveRole(ctx context.Context, request *pb.RemoveRoleRequest) (*pb.Role, error) {
	role, err := d.roles.Retrieve(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("db.Retrieve failed: %w", err)
	}

	var tuples []doorman.Tuple
	tuples, err = d.tuples.ListTuplesForRole(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("ListTuplesForRole failed: %w", err)
	}

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	for _, t := range tuples {
		_, err := d.revokeWithTx(ctx, tx, &pb.RevokeRequest{
			Subject: string(t.Subject),
			Role:    role.ID,
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("revoke failed for %s: %w, %w", t, err, tx.Rollback(ctx))
		}
	}

	if err := d.roles.WithTx(tx).Remove(ctx, request.Id); err != nil {
		return nil, fmt.Errorf("update failed: %w, %w", err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return &pb.Role{}, nil
}

func (d *Doorman) Revoke(ctx context.Context, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	res, err := d.revokeWithTx(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	d.processChangesImmediately()

	return res, nil
}

func (d *Doorman) UpsertRole(ctx context.Context, request *pb.UpsertRoleRequest) (*pb.Role, error) {
	role, err := d.roles.Retrieve(ctx, request.Id)
	if err == db.ErrInvalidRole {
		role = &doorman.Role{ID: request.Id}
	} else if err != nil {
		return nil, fmt.Errorf("db.Retrieve failed: %w", err)
	}

	tuples, err := d.tuples.ListTuplesForRole(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("ListTuplesForRole failed: %w", err)
	}

	// THIS SHOULD BE A TX. IF IT FAILS AT ANY POINT AFTER REVOKING AT LEAST 1 THEN WE ARE F*CKED
	// BECAUSE WE HAVE LOST SOURCE OF TRUF
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	for _, t := range tuples {
		_, err := d.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(t.Subject),
			Role:    role.ID,
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("revoke failed for %s: %w, %w", t, err, tx.Rollback(ctx))
		}
	}

	role.Verbs = []doorman.Verb{}
	for _, v := range request.Verbs {
		role.Verbs = append(role.Verbs, doorman.Verb(v))
	}

	// If this fails all grant requests will fail...
	if err := d.roles.Upsert(ctx, role); err != nil {
		return nil, fmt.Errorf("update failed: %w, %w", err, tx.Rollback(ctx))
	}

	for _, t := range tuples {
		_, err := d.grantWithTx(ctx, tx, &pb.GrantRequest{
			Subject: string(t.Subject),
			Role:    role.ID,
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("grant failed: %w, %w", err, tx.Rollback(ctx))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	d.processChangesImmediately()

	return mapRoleToPb(*role), nil
}

func (d *Doorman) grantWithTx(ctx context.Context, tx pgx.Tx, request *pb.GrantRequest) (*pb.GrantResponse, error) {
	if err := d.tuples.WithTx(tx).Lock(ctx); err != nil {
		return nil, fmt.Errorf("tuples.Lock failed: %w, %w", err, tx.Rollback(ctx))
	}

	tuple := doorman.Tuple{
		Subject: doorman.Object(request.Subject),
		Role:    request.Role,
		Object:  doorman.Object(request.Object),
	}
	if err := d.tuples.WithTx(tx).Add(ctx, tuple); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, fmt.Errorf("rollback failed after failing to add tuple: %w", err)
		}
		return nil, err
	}

	payload, err := json.Marshal(tuple)
	if err != nil {
		return nil, fmt.Errorf("json marshaling failed: %w", err)
	}

	change := doorman.Change{
		ID:        xid.New().String(),
		Type:      "GRANTED",
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := d.changes.WithTx(tx).Add(ctx, change); err != nil {
		return nil, fmt.Errorf("adding change failed: %w", err)
	}

	return &pb.GrantResponse{}, nil
}

func (d *Doorman) processChange(ctx context.Context, change doorman.Change) error {
	switch change.Type {
	case "GRANTED":
		return d.processChangeGrantedOrRevoked(ctx, change)
	case "REVOKED":
		return d.processChangeGrantedOrRevoked(ctx, change)
	default:
		slog.Warn("unhandled change", "type", change.Type)
	}

	return nil
}

func (d *Doorman) processChangeGrantedOrRevoked(ctx context.Context, change doorman.Change) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	// TODO: think if we can remove locking + tx here?
	if err := d.tuples.WithTx(tx).Lock(ctx); err != nil {
		return err
	}

	var tuple doorman.Tuple
	if err := json.Unmarshal(change.Payload, &tuple); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	var staleObjects []doorman.Path

	a := time.Now()
	if change.Type == "GRANTED" {
		staleObjects, err = d.tuples.WithTx(tx).ListConnected(ctx, tuple.Subject, false)
		if err != nil {
			return fmt.Errorf("listConnected failed: %w", err)
		}
	} else {
		staleObjects, err = d.tuples.WithTx(tx).ListConnected(ctx, tuple.Subject, false)
		if err != nil {
			return fmt.Errorf("listConnected failed: %w", err)
		}

		removedPath := doorman.Path{doorman.Connection{Role: tuple.Role, Object: tuple.Object}}
		staleObjects = append(staleObjects, removedPath)

		staleObjectsViaRemoved, err := d.tuples.WithTx(tx).ListConnected(ctx, tuple.Object, false)
		if err != nil {
			return fmt.Errorf("listConnected failed: %w", err)
		}

		for _, incompletePath := range staleObjectsViaRemoved {
			path := append(removedPath, incompletePath...)
			staleObjects = append(staleObjects, path)
		}
	}
	fmt.Println("fetch", time.Since(a))

	a = time.Now()

	if err := d.refreshParents(ctx, tx, tuple.Subject); err != nil {
		return err
	}

	fmt.Println("refreshParents", time.Since(a))

	a = time.Now()

	fmt.Println("stale", len(staleObjects))
	for _, path := range staleObjects {
		if tuple.Subject.Type() == "group" {
			p := path[len(path)-1]
			if err := d.refreshGroups(ctx, tx, p.Object, p.Role); err != nil {
				return err
			}
		}
	}

	fmt.Println("refreshGroups", time.Since(a))

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (d *Doorman) refreshGroups(ctx context.Context, tx pgx.Tx, obj doorman.Object, roleId string) error {
	a := time.Now()
	connectedSubjects, err := d.tuples.WithTx(tx).ListConnected(ctx, obj, true)
	if err != nil {
		return fmt.Errorf("db.ListParents failed: %w", err)
	}

	if false {
		fmt.Println("listconnected", time.Since(a))
	}

	sets := []doorman.Set{}
	for _, path := range connectedSubjects {
		p := path[len(path)-1]
		if path.GroupsOnly() {
			sets = append(sets, doorman.Set{Object: p.Object, Verb: "inherits"})
		}
	}

	role, err := d.roles.Retrieve(ctx, roleId)
	if err != nil {
		// WHY not error? Because this might run after a role is removed (from rebuild cache)
		return nil
	}

	for _, verb := range role.Verbs {
		if err := d.sets.UpdateSubsets(ctx, doorman.Set{Object: obj, Verb: verb}, sets); err != nil {
			return fmt.Errorf("updateSubsets failed: %w", err)
		}
	}

	return nil
}

func (d *Doorman) refreshParents(ctx context.Context, tx pgx.Tx, obj doorman.Object) error {
	parents, err := d.tuples.WithTx(tx).ListParents(ctx, obj)
	if err != nil {
		return fmt.Errorf("db.ListParents failed: %w", err)
	}

	sets := []doorman.Set{}
	for _, tuple := range parents {
		role, err := d.roles.Retrieve(ctx, tuple.Role)
		if err != nil {
			// WHY not error? Because this might run after a role is removed (from rebuild cache)
			return nil
		}
		for _, verb := range role.Verbs {
			sets = append(sets, doorman.Set{Object: tuple.Object, Verb: verb})
		}
	}

	return d.sets.UpdateParents(ctx, obj, sets)
}

func (d *Doorman) revokeWithTx(ctx context.Context, tx pgx.Tx, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	if err := d.tuples.WithTx(tx).Lock(ctx); err != nil {
		return nil, fmt.Errorf("tuples.Lock failed: %w, %w", err, tx.Rollback(ctx))
	}

	tuple := doorman.Tuple{
		Subject: doorman.Object(request.Subject),
		Role:    request.Role,
		Object:  doorman.Object(request.Object),
	}

	if err := d.tuples.WithTx(tx).Remove(ctx, tuple); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, fmt.Errorf("rollback failed after failing to add tuple: %w", err)
		}
		return nil, err
	}

	payload, err := json.Marshal(tuple)
	if err != nil {
		return nil, fmt.Errorf("json marshaling failed: %w", err)
	}

	change := doorman.Change{
		ID:        xid.New().String(),
		Type:      "REVOKED",
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := d.changes.WithTx(tx).Add(ctx, change); err != nil {
		return nil, fmt.Errorf("adding change failed: %w", err)
	}

	return &pb.RevokeResponse{}, nil
}

func NewDoorman(conn *pgxpool.Pool) *Doorman {
	return &Doorman{processing: make(chan bool, 1), conn: conn, changes: db.NewChanges(conn), sets: db.NewSets(conn), roles: db.NewRoles(conn), tuples: db.NewTuples(conn), objects: db.NewObjects(conn)}
}

func mapChangeToPb(c doorman.Change) *pb.Change {
	return &pb.Change{
		Type: c.Type,
	}
}

func mapRoleToPb(r doorman.Role) *pb.Role {
	verbs := make([]string, len(r.Verbs))
	for i, v := range r.Verbs {
		verbs[i] = string(v)
	}
	return &pb.Role{
		Id:    r.ID,
		Verbs: verbs,
	}
}
