package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	projectports "ai_hub.com/app/core/ports/projectports"
	"ai_hub.com/app/infra/config"

	"github.com/redis/go-redis/v9"
)

var _ projectports.ProjectCachePort = (*ProjectCacheAdapter)(nil)

//----------------Adapter configuration ----------------

// ProjectCacheTTLs sets TTLs:
// -EditTTL — for the edit counter project:edit:count:<id>
// -CacheTTL — for JSON snapshot project:cache:<id>.
type ProjectCacheTTLs struct {
	EditTTL  time.Duration
	CacheTTL time.Duration
}

// ProjectCacheAdapter (LEGACY-COMPATIBLE)
// Behavior as in the JS version:
// -IncrementEditCount: INCR + EXPIRE of key project:edit:count:<id>
// -IsInTop100: checks rank in ZSET "project:edit:count" (populates external worker)
// -CacheProject /GetCachedProject: JSON in key project:cache:<id>
// -DeleteFromCache: DEL project:cache:<id>
// -DeleteFromTop100: ZREM with "project:edit:count" + DEL per-project counter key
type ProjectCacheAdapter struct {
	redis *redis.Client
	ttls  ProjectCacheTTLs
}

func NewProjectCacheAdapter(r *redis.Client, ttls ProjectCacheTTLs) *ProjectCacheAdapter {
	return &ProjectCacheAdapter{redis: r, ttls: ttls}
}

// NewProjectCacheAdapterFromConfig -Constructor for DI from config.Env.
// Env expects seconds in a row (default "86400").
func NewProjectCacheAdapterFromConfig(r *redis.Client) *ProjectCacheAdapter {
	editSec, _ := strconv.ParseInt(config.Env.RedisProjectEditTTL, 10, 64)
	cacheSec, _ := strconv.ParseInt(config.Env.RedisProjectCacheTTL, 10, 64)

	return &ProjectCacheAdapter{
		redis: r,
		ttls: ProjectCacheTTLs{
			EditTTL:  time.Duration(editSec) * time.Second,
			CacheTTL: time.Duration(cacheSec) * time.Second,
		},
	}
}

// key builders ---------------------------------------------------------------

func keyEditCounter(projectID string) string {
	return fmt.Sprintf("project:edit:count:%s", projectID)
}

func keyProjectCache(projectID string) string {
	return fmt.Sprintf("project:cache:%s", projectID)
}

const zsetTopKey = "project:edit:count"

// IncrementEditCount increments the project's edit counter and updates the TTL.
// Key: project:edit:count:<projectId>
func (a *ProjectCacheAdapter) IncrementEditCount(ctx context.Context, projectID string) error {
	key := keyEditCounter(projectID)

	//increment the individual project counter
	newScore, err := a.redis.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	//update TTL
	if a.ttls.EditTTL > 0 {
		_ = a.redis.Expire(ctx, key, a.ttls.EditTTL).Err()
	}

	//ADD to the general top
	_ = a.redis.ZAdd(ctx, zsetTopKey, redis.Z{
		Member: projectID,
		Score:  float64(newScore),
	}).Err()

	return nil
}

// IsInTop100 checks if the project is in the top 100 by number of edits.
// ZSET: "project:edit:count" (maintained by external process).
func (a *ProjectCacheAdapter) IsInTop100(ctx context.Context, projectID string) (bool, error) {
	rank, err := a.redis.ZRevRank(ctx, zsetTopKey, projectID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return rank >= 0 && rank < 100, nil
}

// CacheProject stores the project snapshot as JSON with TTL.
// Key: project:cache:<projectId>
func (a *ProjectCacheAdapter) CacheProject(ctx context.Context, projectID string, projectData any) error {
	key := keyProjectCache(projectID)
	payload, err := json.Marshal(projectData)
	if err != nil {
		return err
	}
	return a.redis.Set(ctx, key, payload, a.ttls.CacheTTL).Err()
}

// GetCachedProject reads and returns raw JSON as map[string]any (or nil).
func (a *ProjectCacheAdapter) GetCachedProject(ctx context.Context, projectID string) (any, error) {
	key := keyProjectCache(projectID)
	raw, err := a.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteFromCache deletes the JSON snapshot of the project.
func (a *ProjectCacheAdapter) DeleteFromCache(ctx context.Context, projectID string) error {
	key := keyProjectCache(projectID)
	if err := a.redis.Del(ctx, key).Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// DeleteFromTop100 removes the project from ZSET "project:edit:count" and cleans
// per-project counter key. Both operations are best-effort and idempotent.
func (a *ProjectCacheAdapter) DeleteFromTop100(ctx context.Context, projectID string) error {
	//ZREM from the rating list
	if err := a.redis.ZRem(ctx, zsetTopKey, projectID).Err(); err != nil && err != redis.Nil {
		return err
	}
	//We delete the local edit counter so as not to leave garbage
	if err := a.redis.Del(ctx, keyEditCounter(projectID)).Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}
