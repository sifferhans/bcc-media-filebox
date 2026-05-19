package config

import (
	"context"
	"fmt"
	"log"
	"os"

	db "filebox/internal/db/gen"
)

// Target is the in-memory shape of an upload destination. It is no longer the
// runtime source of truth — targets live in the `targets` table now — but the
// type stays around for the env-var bootstrap path and any callers that want
// to pass a target snapshot through configuration.
type Target struct {
	Name string
	Dir  string
}

// LoadTargetsFromEnv reads numbered TARGET_N_NAME / TARGET_N_DIR environment
// variable pairs (starting at 1). Unlike the historical LoadTargets, it does
// NOT fail when zero targets are configured — that's a valid state for a
// fresh deploy where the operator plans to create targets from the admin UI.
// It still rejects malformed pairs and missing directories so that bootstrap
// either succeeds completely or fails loudly.
func LoadTargetsFromEnv() ([]Target, error) {
	var targets []Target

	for i := 1; ; i++ {
		name := os.Getenv(fmt.Sprintf("TARGET_%d_NAME", i))
		dir := os.Getenv(fmt.Sprintf("TARGET_%d_DIR", i))

		if name == "" && dir == "" {
			break
		}
		if name == "" {
			return nil, fmt.Errorf("TARGET_%d_DIR is set but TARGET_%d_NAME is missing", i, i)
		}
		if dir == "" {
			return nil, fmt.Errorf("TARGET_%d_NAME is set but TARGET_%d_DIR is missing", i, i)
		}

		info, err := os.Stat(dir)
		if err != nil {
			return nil, fmt.Errorf("target %q directory %q: %w", name, dir, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("target %q: %q is not a directory", name, dir)
		}

		targets = append(targets, Target{Name: name, Dir: dir})
	}

	return targets, nil
}

// BootstrapTargets seeds the targets table from env vars on first start. If
// the table already has rows, the env vars are ignored — the DB is the
// source of truth from that point on. This is a one-way migration: edits in
// the UI never get written back to TARGET_N_* env vars.
func BootstrapTargets(ctx context.Context, queries *db.Queries, envTargets []Target) error {
	existing, err := queries.CountTargets(ctx)
	if err != nil {
		return fmt.Errorf("count targets: %w", err)
	}
	if existing > 0 {
		if len(envTargets) > 0 {
			log.Printf("targets table already populated (%d rows); ignoring TARGET_N_* env vars", existing)
		}
		return nil
	}
	if len(envTargets) == 0 {
		log.Println("no targets configured (DB empty, no TARGET_N_* env vars) — create one from /admin once signed in as admin")
		return nil
	}
	for _, t := range envTargets {
		if _, err := queries.CreateTarget(ctx, db.CreateTargetParams{Name: t.Name, Path: t.Dir}); err != nil {
			return fmt.Errorf("seed target %q: %w", t.Name, err)
		}
		log.Printf("seeded target %q -> %s", t.Name, t.Dir)
	}
	return nil
}

// TargetDirFromDB resolves a target name to its filesystem directory via the
// DB. Returns ("", false) if the target doesn't exist. Validation that the
// directory exists and is writable is left to the admin API on create/update;
// the TUS hook re-validates lazily and refuses to write if the directory has
// since vanished.
func TargetDirFromDB(ctx context.Context, queries *db.Queries, name string) (string, bool) {
	t, err := queries.GetTargetByName(ctx, name)
	if err != nil {
		return "", false
	}
	return t.Path, true
}

// TargetDir is kept for backwards compatibility with code paths that still
// hand a snapshot slice around. Prefer TargetDirFromDB for runtime lookups.
func TargetDir(targets []Target, name string) (string, bool) {
	for _, t := range targets {
		if t.Name == name {
			return t.Dir, true
		}
	}
	return "", false
}
