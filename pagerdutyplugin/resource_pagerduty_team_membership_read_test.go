package pagerduty

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	pathpkg "path"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestTeamMembershipReadCacheLookupReusesSnapshot(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(false)
	var fetchCount atomic.Int32

	fetch := testCountingSnapshotFetcher(&fetchCount, testSnapshot(1,
		testMemberRole("user-a", "manager"),
		testMemberRole("user-b", "observer"),
	))

	snapA, err := cache.getSnapshot(context.Background(), "team-1", fetch)
	if err != nil {
		t.Fatalf("getSnapshot failed: %v", err)
	}

	snapB, err := cache.getSnapshot(context.Background(), "team-1", fetch)
	if err != nil {
		t.Fatalf("getSnapshot failed: %v", err)
	}

	if got := fetchCount.Load(); got != 1 {
		t.Fatalf("expected exactly one fetch, got %d", got)
	}
	if snapA.generation != 1 || snapB.generation != 1 {
		t.Fatalf("expected cached generation 1, got %d and %d", snapA.generation, snapB.generation)
	}
}

func TestTeamMembershipReadCacheLookupReturnsMemberFromSnapshot(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(false)
	var fetchCount atomic.Int32

	fetch := testCountingSnapshotFetcher(&fetchCount, testSnapshot(1,
		testMemberRole("user-a", "manager"),
		testMemberRole("user-b", "observer"),
	))

	memberA, foundA, err := cache.lookup(context.Background(), "team-1", "user-a", fetch)
	if err != nil {
		t.Fatalf("lookup user-a failed: %v", err)
	}
	if !foundA || memberA.Role != "manager" {
		t.Fatalf("unexpected lookup result for user-a: found=%t role=%q", foundA, memberA.Role)
	}

	memberB, foundB, err := cache.lookup(context.Background(), "team-1", "user-b", fetch)
	if err != nil {
		t.Fatalf("lookup user-b failed: %v", err)
	}
	if !foundB || memberB.Role != "observer" {
		t.Fatalf("unexpected lookup result for user-b: found=%t role=%q", foundB, memberB.Role)
	}

	if got := fetchCount.Load(); got != 1 {
		t.Fatalf("expected exactly one fetch, got %d", got)
	}
}

func TestTeamMembershipReadCacheInvalidateForcesRefetch(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(false)
	var fetchCount atomic.Int32

	fetch := func(_ context.Context, teamID string) (*teamMembershipSnapshot, error) {
		count := fetchCount.Add(1)
		return testSnapshot(int(count),
			testMemberRole("user-a", "manager"),
		), nil
	}

	snapA, err := cache.getSnapshot(context.Background(), "team-1", fetch)
	if err != nil {
		t.Fatalf("first getSnapshot failed: %v", err)
	}

	cache.invalidate("team-1", "test")

	snapB, err := cache.getSnapshot(context.Background(), "team-1", fetch)
	if err != nil {
		t.Fatalf("second getSnapshot failed: %v", err)
	}

	if got := fetchCount.Load(); got != 2 {
		t.Fatalf("expected refetch after invalidation, got %d fetches", got)
	}
	if snapA.generation == snapB.generation {
		t.Fatalf("expected a new generation after invalidation, got %d and %d", snapA.generation, snapB.generation)
	}
}

func TestTeamMembershipReadCacheConcurrentLookupDeduplicatesFetch(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(false)
	var fetchCount atomic.Int32
	started := make(chan struct{})
	var startOnce sync.Once

	fetch := func(_ context.Context, teamID string) (*teamMembershipSnapshot, error) {
		fetchCount.Add(1)
		startOnce.Do(func() { close(started) })
		time.Sleep(25 * time.Millisecond)
		return testSnapshot(1, testMemberRole("user-a", "manager")), nil
	}

	const workers = 8
	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, found, err := cache.lookup(context.Background(), "team-1", "user-a", fetch)
			if err != nil {
				errs <- err
				return
			}
			if !found {
				errs <- errLookupNotFound
			}
		}()
	}

	<-started
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent lookup failed: %v", err)
		}
	}

	if got := fetchCount.Load(); got != 1 {
		t.Fatalf("expected a single fetch for concurrent lookups, got %d", got)
	}
}

func TestTeamMembershipReadCacheDisabledBypassesStorage(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(true)
	var fetchCount atomic.Int32

	fetch := testCountingSnapshotFetcher(&fetchCount, testSnapshot(1,
		testMemberRole("user-a", "manager"),
	))

	for i := 0; i < 2; i++ {
		_, found, err := cache.lookup(context.Background(), "team-1", "user-a", fetch)
		if err != nil {
			t.Fatalf("lookup %d failed: %v", i, err)
		}
		if !found {
			t.Fatalf("lookup %d did not find expected user", i)
		}
	}

	if got := fetchCount.Load(); got != 2 {
		t.Fatalf("expected cache-disabled path to fetch every time, got %d", got)
	}
}

func TestTeamMembershipReadCacheFetchErrorDoesNotPoisonCache(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(false)
	var fetchCount atomic.Int32

	fetch := func(_ context.Context, teamID string) (*teamMembershipSnapshot, error) {
		count := fetchCount.Add(1)
		if count == 1 {
			return nil, errors.New("transient fetch error")
		}
		return testSnapshot(1, testMemberRole("user-a", "manager")), nil
	}

	_, found, err := cache.lookup(context.Background(), "team-1", "user-a", fetch)
	if err == nil {
		t.Fatal("expected first lookup to fail")
	}
	if found {
		t.Fatal("expected first lookup to not report found")
	}

	member, found, err := cache.lookup(context.Background(), "team-1", "user-a", fetch)
	if err != nil {
		t.Fatalf("second lookup failed: %v", err)
	}
	if !found || member.Role != "manager" {
		t.Fatalf("unexpected second lookup result: found=%t role=%q", found, member.Role)
	}
	if got := fetchCount.Load(); got != 2 {
		t.Fatalf("expected fetch error path to refetch, got %d fetches", got)
	}
}

func TestFetchTeamMembershipSnapshotPaginates(t *testing.T) {
	t.Parallel()

	client, requests := newTestTeamMembershipClient(t, map[string][]testTeamMemberPage{
		"team-1": {
			{
				Members: []testTeamMember{
					{UserID: "user-a", Role: "manager"},
					{UserID: "user-b", Role: "observer"},
				},
				More:   true,
				Limit:  2,
				Offset: 0,
			},
			{
				Members: []testTeamMember{
					{UserID: "user-c", Role: "responder"},
				},
				More:   false,
				Limit:  2,
				Offset: 2,
			},
		},
	})

	snapshot, err := fetchTeamMembershipSnapshot(context.Background(), client, "team-1")
	if err != nil {
		t.Fatalf("fetchTeamMembershipSnapshot failed: %v", err)
	}

	if snapshot.memberCount != 3 {
		t.Fatalf("expected 3 members, got %d", snapshot.memberCount)
	}
	if snapshot.pagesFetched != 2 {
		t.Fatalf("expected 2 pages fetched, got %d", snapshot.pagesFetched)
	}
	if snapshot.membersByID["user-c"].Role != "responder" {
		t.Fatalf("expected user-c role responder, got %q", snapshot.membersByID["user-c"].Role)
	}
	if got := requests.Load(); got != 2 {
		t.Fatalf("expected 2 API requests, got %d", got)
	}
}

func TestFetchTeamMembershipSnapshotHandlesEmptyTeam(t *testing.T) {
	t.Parallel()

	client, requests := newTestTeamMembershipClient(t, map[string][]testTeamMemberPage{
		"team-1": {
			{
				Members: nil,
				More:    false,
				Limit:   100,
				Offset:  0,
			},
		},
	})

	snapshot, err := fetchTeamMembershipSnapshot(context.Background(), client, "team-1")
	if err != nil {
		t.Fatalf("fetchTeamMembershipSnapshot failed: %v", err)
	}

	if snapshot.memberCount != 0 {
		t.Fatalf("expected 0 members, got %d", snapshot.memberCount)
	}
	if snapshot.pagesFetched != 1 {
		t.Fatalf("expected 1 page fetched for empty team, got %d", snapshot.pagesFetched)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected 1 API request, got %d", got)
	}
}

func TestRequestGetTeamMembershipUsesCacheForSteadyStateRead(t *testing.T) {
	client, requests := newTestTeamMembershipClient(t, map[string][]testTeamMemberPage{
		"team-1": {
			{
				Members: []testTeamMember{
					{UserID: "user-a", Role: "manager"},
					{UserID: "user-b", Role: "observer"},
				},
				More:   false,
				Limit:  100,
				Offset: 0,
			},
		},
	})
	enableTeamMembershipReadCacheForClient(client)
	var diags diag.Diagnostics

	modelA, err := requestGetTeamMembership(context.Background(), client, "user-a:team-1", nil, teamMembershipReadSteadyState, false, &diags)
	if err != nil {
		t.Fatalf("first requestGetTeamMembership failed: %v", err)
	}
	modelB, err := requestGetTeamMembership(context.Background(), client, "user-b:team-1", nil, teamMembershipReadSteadyState, false, &diags)
	if err != nil {
		t.Fatalf("second requestGetTeamMembership failed: %v", err)
	}

	if modelA.Role.ValueString() != "manager" || modelB.Role.ValueString() != "observer" {
		t.Fatalf("unexpected roles: %q %q", modelA.Role.ValueString(), modelB.Role.ValueString())
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected steady-state reads to reuse cached snapshot, got %d requests", got)
	}
}

func TestRequestGetTeamMembershipBypassesCacheForPostWriteVerification(t *testing.T) {
	client, requests := newTestTeamMembershipClient(t, map[string][]testTeamMemberPage{
		"team-1": {
			{
				Members: []testTeamMember{
					{UserID: "user-a", Role: "manager"},
				},
				More:   false,
				Limit:  100,
				Offset: 0,
			},
		},
	})
	enableTeamMembershipReadCacheForClient(client)
	var diags diag.Diagnostics

	_, err := requestGetTeamMembership(context.Background(), client, "user-a:team-1", nil, teamMembershipReadSteadyState, false, &diags)
	if err != nil {
		t.Fatalf("cached requestGetTeamMembership failed: %v", err)
	}

	role := "manager"
	_, err = requestGetTeamMembership(context.Background(), client, "user-a:team-1", &role, teamMembershipReadPostWriteVerification, true, &diags)
	if err != nil {
		t.Fatalf("post-write verification requestGetTeamMembership failed: %v", err)
	}

	if got := requests.Load(); got != 2 {
		t.Fatalf("expected post-write verification path to bypass cache, got %d requests", got)
	}
}

func TestRequestGetTeamMembershipReturnsNotFoundForMissingMember(t *testing.T) {
	client, requests := newTestTeamMembershipClient(t, map[string][]testTeamMemberPage{
		"team-1": {
			{
				Members: []testTeamMember{
					{UserID: "user-a", Role: "manager"},
				},
				More:   false,
				Limit:  100,
				Offset: 0,
			},
		},
	})
	enableTeamMembershipReadCacheForClient(client)
	var diags diag.Diagnostics

	_, err := requestGetTeamMembership(context.Background(), client, "user-missing:team-1", nil, teamMembershipReadSteadyState, false, &diags)
	if err == nil {
		t.Fatal("expected missing member lookup to fail")
	}
	if !util.IsNotFoundError(err) {
		t.Fatalf("expected not found error, got %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected one API request for missing member lookup, got %d", got)
	}
}

func TestRequestGetTeamMembershipReturnsNotFoundForMissingTeam(t *testing.T) {
	t.Parallel()

	client, requests := newTestTeamMembershipClient(t, map[string][]testTeamMemberPage{})
	var diags diag.Diagnostics

	_, err := requestGetTeamMembership(context.Background(), client, "user-a:team-missing", nil, teamMembershipReadSteadyState, false, &diags)
	if err == nil {
		t.Fatal("expected missing team lookup to fail")
	}
	if !util.IsNotFoundError(err) {
		t.Fatalf("expected not found error, got %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected one API request for missing team lookup, got %d", got)
	}
}

func TestTeamMembershipReadCacheInvalidateDuringInflightFetchRefetches(t *testing.T) {
	t.Parallel()

	cache := newTeamMembershipReadCacheWithDisabled(false)
	var fetchCount atomic.Int32
	firstFetchStarted := make(chan struct{})
	releaseFirstFetch := make(chan struct{})

	fetch := func(_ context.Context, teamID string) (*teamMembershipSnapshot, error) {
		count := fetchCount.Add(1)
		if count == 1 {
			close(firstFetchStarted)
			<-releaseFirstFetch
		}

		role := "observer"
		if count == 1 {
			role = "manager"
		}

		return testSnapshot(1, testMemberRole("user-a", role)), nil
	}

	resultCh := make(chan teamMembershipCacheMember, 1)
	errCh := make(chan error, 1)

	go func() {
		member, found, err := cache.lookup(context.Background(), "team-1", "user-a", fetch)
		if err != nil {
			errCh <- err
			return
		}
		if !found {
			errCh <- errLookupNotFound
			return
		}
		resultCh <- member
	}()

	<-firstFetchStarted
	cache.invalidate("team-1", "test")
	close(releaseFirstFetch)

	select {
	case err := <-errCh:
		t.Fatalf("lookup failed: %v", err)
	case member := <-resultCh:
		if member.Role != "observer" {
			t.Fatalf("expected refetched role after invalidation, got %q", member.Role)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for lookup to complete after invalidation")
	}

	if got := fetchCount.Load(); got != 2 {
		t.Fatalf("expected invalidation during inflight fetch to trigger a second fetch, got %d", got)
	}
}

var errLookupNotFound = &lookupError{message: "lookup did not find expected user"}

type lookupError struct {
	message string
}

func (e *lookupError) Error() string {
	return e.message
}

func testSnapshot(pagesFetched int, members ...testTeamMembership) *teamMembershipSnapshot {
	snapshotMembers := make(map[string]teamMembershipCacheMember, len(members))
	for _, member := range members {
		snapshotMembers[member.userID] = teamMembershipCacheMember{Role: member.role}
	}

	return &teamMembershipSnapshot{
		fetchedAt:    time.Now(),
		membersByID:  snapshotMembers,
		memberCount:  len(members),
		pagesFetched: pagesFetched,
	}
}

func enableTeamMembershipReadCacheForClient(client *pagerduty.Client) {
	globalTeamMembershipReadCaches = newTeamMembershipReadCacheRegistry()
	globalTeamMembershipReadCaches.caches[client] = newTeamMembershipReadCacheWithDisabled(false)
}

func testCountingSnapshotFetcher(fetchCount *atomic.Int32, snapshot *teamMembershipSnapshot) teamMembershipSnapshotFetcher {
	return func(_ context.Context, teamID string) (*teamMembershipSnapshot, error) {
		fetchCount.Add(1)
		return snapshot, nil
	}
}

func testMemberRole(userID, role string) testTeamMembership {
	return testTeamMembership{
		userID: userID,
		role:   role,
	}
}

type testTeamMembership struct {
	userID string
	role   string
}

type testTeamMember struct {
	UserID string
	Role   string
}

type testTeamMemberPage struct {
	Members []testTeamMember
	More    bool
	Limit   uint
	Offset  uint
}

func newTestTeamMembershipClient(t *testing.T, pages map[string][]testTeamMemberPage) (*pagerduty.Client, *atomic.Int32) {
	t.Helper()

	var requests atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		requests.Add(1)
		teamID := pathpkg.Base(pathpkg.Dir(r.URL.Path))
		offset := uint(0)
		if v := r.URL.Query().Get("offset"); v != "" {
			i, err := strconv.Atoi(v)
			if err != nil {
				t.Fatalf("invalid offset %q: %v", v, err)
			}
			offset = uint(i)
		}

		var page testTeamMemberPage
		found := false
		for _, candidate := range pages[teamID] {
			if candidate.Offset == offset {
				page = candidate
				found = true
				break
			}
		}
		if !found {
			http.Error(w, `{"error":{"message":"not found","code":2100}}`, http.StatusNotFound)
			return
		}

		var members []map[string]any
		for _, member := range page.Members {
			members = append(members, map[string]any{
				"user": map[string]any{
					"id":   member.UserID,
					"type": "user_reference",
				},
				"role": member.Role,
			})
		}

		resp := map[string]any{
			"members": members,
			"limit":   page.Limit,
			"offset":  page.Offset,
			"more":    page.More,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := pagerduty.NewClient("test-token", WithHTTPClient(server.Client()), pagerduty.WithAPIEndpoint(server.URL))
	return client, &requests
}
