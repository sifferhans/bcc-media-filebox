package api

import (
	"context"
	"strconv"
	"strings"

	"filebox/internal/auth"
	db "filebox/internal/db/gen"
)

// AllowedTargets is the effective access a caller has across upload targets.
//
//   - All=true means the caller can see every target (admin or any grant with
//     all_targets=1); IDs is unused in that case.
//   - Admin=true means the caller has admin privileges (RBAC for /api/admin).
//
// When All=false, IDs is the union of target IDs across every grant that
// matched the caller (direct user, builtin-group wildcard, custom-group
// membership). An empty IDs set with All=false means "no upload access" —
// the user can sign in but every target picker should be empty.
type AllowedTargets struct {
	Admin bool
	All   bool
	IDs   map[int64]struct{}
}

// EffectiveTargetIDs resolves the caller's grants into an AllowedTargets value.
// The users.role column is treated as authoritative for admin status — it's a
// denormalised cache kept in sync by RecomputeUserRoles. Grants are still
// consulted for non-admin target access.
func EffectiveTargetIDs(ctx context.Context, queries *db.Queries, caller *auth.Caller) (AllowedTargets, error) {
	out := AllowedTargets{IDs: map[int64]struct{}{}}
	if caller == nil {
		return out, nil
	}
	if caller.Role == "admin" {
		out.Admin = true
		out.All = true
		return out, nil
	}
	if caller.Email == "" {
		return out, nil
	}
	rows, err := queries.GrantsForUser(ctx, db.GrantsForUserParams{
		Email:    caller.Email,
		Provider: caller.Provider,
	})
	if err != nil {
		return out, err
	}
	for _, g := range rows {
		if g.Admin != 0 {
			out.Admin = true
			out.All = true
		}
		if g.AllTargets != 0 {
			out.All = true
		}
		for _, idStr := range splitTargetIDs(g.TargetIds) {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				out.IDs[id] = struct{}{}
			}
		}
	}
	return out, nil
}

func splitTargetIDs(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

// Role recomputation lives in the auth package — see auth.RecomputeUserRoles
// and auth.RecomputeRoleForUser. Re-exported here so the admin handlers can
// stay self-contained.
var (
	RecomputeUserRoles   = auth.RecomputeUserRoles
	RecomputeRoleForUser = auth.RecomputeRoleForUser
)
