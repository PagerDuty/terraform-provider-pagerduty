package pagerduty

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// This file keeps the pagerduty_team_membership read path scalable for large teams.
// Steady-state reads share one team snapshot per client/team within a single run,
// while post-write verification always reads through to the API.

const enableTeamMembershipReadCacheEnv = "PAGERDUTY_ENABLE_TEAM_MEMBERSHIP_READ_CACHE"

type teamMembershipSnapshotFetcher func(context.Context, string) (*teamMembershipSnapshot, error)
type teamMembershipMemberFetcher func(context.Context, string, string) (teamMembershipCacheMember, bool, error)

type teamMembershipReadMode int

const (
	// teamMembershipReadSteadyState is the normal refresh/read path. It may use
	// the shared in-memory snapshot cache for repeated lookups within the same run.
	teamMembershipReadSteadyState teamMembershipReadMode = iota
	// teamMembershipReadPostWriteVerification is used immediately after create/update
	// to confirm the remote role state. It always bypasses the shared read cache.
	teamMembershipReadPostWriteVerification
)

type teamMembershipCacheMember struct {
	Role string
}

type teamMembershipSnapshot struct {
	generation   uint64
	fetchedAt    time.Time
	membersByID  map[string]teamMembershipCacheMember
	memberCount  int
	pagesFetched int
}

type teamMembershipReadCacheEntry struct {
	snapshot    *teamMembershipSnapshot
	waitCh      chan struct{}
	invalidated bool
}

type teamMembershipSnapshotAccess struct {
	snapshot *teamMembershipSnapshot
	waitCh   chan struct{}
}

type teamMembershipReadCache struct {
	mu          sync.Mutex
	entries     map[string]*teamMembershipReadCacheEntry
	generations map[string]uint64
	disabled    bool
}

type teamMembershipReadCacheRegistry struct {
	mu     sync.Mutex
	caches map[*pagerduty.Client]*teamMembershipReadCache
}

// Cache instances are scoped per PagerDuty client so provider aliases do not
// share team snapshots with each other.
var globalTeamMembershipReadCaches = newTeamMembershipReadCacheRegistry()

func newTeamMembershipReadCacheRegistry() *teamMembershipReadCacheRegistry {
	return &teamMembershipReadCacheRegistry{
		caches: map[*pagerduty.Client]*teamMembershipReadCache{},
	}
}

func newTeamMembershipReadCache() *teamMembershipReadCache {
	_, enabled := os.LookupEnv(enableTeamMembershipReadCacheEnv)
	if enabled {
		log.Printf("[INFO] team-members-cache enabled via %s", enableTeamMembershipReadCacheEnv)
	}
	return newTeamMembershipReadCacheWithDisabled(!enabled)
}

func newTeamMembershipReadCacheWithDisabled(disabled bool) *teamMembershipReadCache {
	return &teamMembershipReadCache{
		entries:     map[string]*teamMembershipReadCacheEntry{},
		generations: map[string]uint64{},
		disabled:    disabled,
	}
}

func (registry *teamMembershipReadCacheRegistry) cacheForClient(client *pagerduty.Client) *teamMembershipReadCache {
	if client == nil {
		return nil
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if cache, ok := registry.caches[client]; ok {
		return cache
	}

	cache := newTeamMembershipReadCache()
	registry.caches[client] = cache
	return cache
}

func invalidateTeamMembershipReadCache(client *pagerduty.Client, teamID, reason string) {
	cache := globalTeamMembershipReadCaches.cacheForClient(client)
	if cache == nil {
		return
	}
	cache.invalidate(teamID, reason)
}

func (cache *teamMembershipReadCache) lookup(ctx context.Context, teamID, userID string, fetch teamMembershipSnapshotFetcher) (teamMembershipCacheMember, bool, error) {
	snapshot, err := cache.getSnapshot(ctx, teamID, fetch)
	if err != nil {
		return teamMembershipCacheMember{}, false, err
	}

	member, found := snapshot.membersByID[userID]
	return member, found, nil
}

func (cache *teamMembershipReadCache) getSnapshot(ctx context.Context, teamID string, fetch teamMembershipSnapshotFetcher) (*teamMembershipSnapshot, error) {
	if cache == nil {
		return fetch(ctx, teamID)
	}
	if cache.disabled {
		log.Printf("[DEBUG] team-members-cache bypass team=%s reason=disabled", teamID)
		return fetch(ctx, teamID)
	}

	for {
		// Same-team lookups share one in-flight fetch so refresh can stay parallel
		// across teams without re-reading the same team pages for each resource.
		access := cache.getCachedSnapshotOrWait(teamID)
		if access.snapshot != nil {
			return access.snapshot, nil
		}
		if access.waitCh != nil {
			if err := waitForTeamMembershipSnapshot(ctx, access.waitCh); err != nil {
				return nil, err
			}
			continue
		}

		waitCh := cache.beginSnapshotFetch(teamID)

		log.Printf("[DEBUG] team-members-cache miss team=%s reason=not_present", teamID)
		snapshot, err := fetch(ctx, teamID)

		storedSnapshot, needsRetry, storeErr := cache.storeFetchedSnapshot(teamID, waitCh, snapshot, err)
		if storeErr != nil {
			return nil, storeErr
		}
		if needsRetry {
			continue
		}
		return storedSnapshot, nil
	}
}

func (cache *teamMembershipReadCache) getCachedSnapshotOrWait(teamID string) teamMembershipSnapshotAccess {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, ok := cache.entries[teamID]
	if !ok {
		return teamMembershipSnapshotAccess{}
	}
	if entry.snapshot != nil {
		return teamMembershipSnapshotAccess{snapshot: entry.snapshot}
	}

	return teamMembershipSnapshotAccess{waitCh: entry.waitCh}
}

func (cache *teamMembershipReadCache) beginSnapshotFetch(teamID string) chan struct{} {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	waitCh := make(chan struct{})
	cache.entries[teamID] = &teamMembershipReadCacheEntry{waitCh: waitCh}
	return waitCh
}

func waitForTeamMembershipSnapshot(ctx context.Context, waitCh chan struct{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitCh:
		return nil
	}
}

func (cache *teamMembershipReadCache) storeFetchedSnapshot(teamID string, waitCh chan struct{}, snapshot *teamMembershipSnapshot, fetchErr error) (*teamMembershipSnapshot, bool, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	current := cache.entries[teamID]
	entryChanged := current == nil || current.waitCh != waitCh || current.snapshot != nil
	if entryChanged {
		if fetchErr != nil {
			return nil, false, fetchErr
		}
		return snapshot, false, nil
	}

	if current.invalidated {
		// An in-flight fetch cannot be cancelled safely, so invalidate marks its
		// result as disposable and forces the caller to refetch on completion.
		delete(cache.entries, teamID)
		close(waitCh)
		if fetchErr != nil {
			return nil, false, fetchErr
		}
		return nil, true, nil
	}

	if fetchErr != nil {
		delete(cache.entries, teamID)
		close(waitCh)
		return nil, false, fetchErr
	}

	generation := cache.generations[teamID] + 1
	cache.generations[teamID] = generation
	snapshot.generation = generation
	current.snapshot = snapshot
	close(waitCh)

	logTeamMembershipCacheStore(teamID, snapshot)
	return snapshot, false, nil
}

func (cache *teamMembershipReadCache) invalidate(teamID, reason string) {
	if cache == nil || cache.disabled {
		return
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, ok := cache.entries[teamID]
	if !ok {
		return
	}

	var generation uint64
	if entry.snapshot != nil {
		generation = entry.snapshot.generation
	}

	if entry.snapshot == nil && entry.waitCh != nil {
		// Writes can race with an in-flight read. Marking the entry invalid lets
		// the fetch finish, but prevents that stale snapshot from being reused.
		entry.invalidated = true
		logTeamMembershipCacheInvalidate(teamID, generation, reason)
		return
	}

	delete(cache.entries, teamID)
	logTeamMembershipCacheInvalidate(teamID, generation, reason)
}

func requestGetTeamMembership(ctx context.Context, client *pagerduty.Client, id string, neededRole *string, mode teamMembershipReadMode, retryNotFound bool, diags *diag.Diagnostics) (resourceTeamMembershipModel, error) {
	if mode == teamMembershipReadPostWriteVerification {
		return requestGetTeamMembershipUncached(ctx, client, id, neededRole, retryNotFound, diags)
	}

	cache := globalTeamMembershipReadCaches.cacheForClient(client)
	return requestGetTeamMembershipWithFetcher(ctx, id, neededRole, retryNotFound, diags, func(ctx context.Context, teamID, userID string) (teamMembershipCacheMember, bool, error) {
		return cache.lookup(ctx, teamID, userID, func(ctx context.Context, teamID string) (*teamMembershipSnapshot, error) {
			return fetchTeamMembershipSnapshot(ctx, client, teamID)
		})
	})
}

func logTeamMembershipCacheStore(teamID string, snapshot *teamMembershipSnapshot) {
	log.Printf("[DEBUG] team-members-cache store team=%s generation=%d members=%d pages=%d", teamID, snapshot.generation, snapshot.memberCount, snapshot.pagesFetched)
}

func logTeamMembershipCacheInvalidate(teamID string, generation uint64, reason string) {
	log.Printf("[DEBUG] team-members-cache invalidate team=%s generation=%d reason=%s", teamID, generation, reason)
}

func requestGetTeamMembershipUncached(ctx context.Context, client *pagerduty.Client, id string, neededRole *string, retryNotFound bool, diags *diag.Diagnostics) (resourceTeamMembershipModel, error) {
	// Post-write verification reads fresh remote state so the cache cannot mask
	// eventual consistency or role propagation delays.
	return requestGetTeamMembershipWithFetcher(ctx, id, neededRole, retryNotFound, diags, func(ctx context.Context, teamID, userID string) (teamMembershipCacheMember, bool, error) {
		snapshot, err := fetchTeamMembershipSnapshot(ctx, client, teamID)
		if err != nil {
			return teamMembershipCacheMember{}, false, err
		}

		member, found := snapshot.membersByID[userID]
		return member, found, nil
	})
}

func requestGetTeamMembershipWithFetcher(ctx context.Context, id string, neededRole *string, retryNotFound bool, diags *diag.Diagnostics, fetchMember teamMembershipMemberFetcher) (resourceTeamMembershipModel, error) {
	var model resourceTeamMembershipModel

	userID, teamID, err := util.ResourcePagerDutyParseColonCompoundID(id)
	if err != nil {
		diags.AddError(fmt.Sprintf("Invalid Team Membership ID %s", id), err.Error())
		return model, nil
	}

	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		member, found, err := fetchMember(ctx, teamID, userID)
		if err != nil {
			return teamMembershipRetryError(err, retryNotFound)
		}

		if !found {
			return teamMembershipNotFoundRetryError(retryNotFound)
		}

		if neededRole != nil && member.Role != *neededRole {
			err = fmt.Errorf("Role %q fetched is different from configuration %q", member.Role, *neededRole)
			return retry.RetryableError(err)
		}

		model = flattenTeamMembership(userID, teamID, member.Role)
		return nil
	})

	return model, err
}

func teamMembershipRetryError(err error, retryNotFound bool) *retry.RetryError {
	if util.IsBadRequestError(err) {
		return retry.NonRetryableError(err)
	}
	if !retryNotFound && util.IsNotFoundError(err) {
		return retry.NonRetryableError(err)
	}
	return retry.RetryableError(err)
}

func teamMembershipNotFoundRetryError(retryNotFound bool) *retry.RetryError {
	err := pagerduty.APIError{StatusCode: http.StatusNotFound}
	if retryNotFound {
		return retry.RetryableError(err)
	}
	return retry.NonRetryableError(err)
}

func fetchTeamMembershipSnapshot(ctx context.Context, client *pagerduty.Client, teamID string) (*teamMembershipSnapshot, error) {
	snapshot := &teamMembershipSnapshot{
		fetchedAt:   time.Now(),
		membersByID: map[string]teamMembershipCacheMember{},
	}

	offset := uint(0)
	more := true

	for more {
		resp, err := client.ListTeamMembers(ctx, teamID, pagerduty.ListTeamMembersOptions{
			Limit:  100,
			Offset: offset,
		})
		if err != nil {
			return nil, err
		}

		snapshot.pagesFetched++
		for _, m := range resp.Members {
			snapshot.membersByID[m.User.ID] = teamMembershipCacheMember{Role: m.Role}
		}

		more = resp.More
		offset += resp.Limit
	}

	snapshot.memberCount = len(snapshot.membersByID)
	return snapshot, nil
}
