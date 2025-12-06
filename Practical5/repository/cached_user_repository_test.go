package repository

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCacheHitMiss(t *testing.T) {
	ctx := context.Background()
	repo := NewCachedUserRepository(cachedTestDB, cachedTestRedis)

	// Clear cache before test
	cachedTestRedis.FlushAll(ctx)

	// First call - cache miss
	user1, err := repo.GetByIDCached(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Second call - should hit cache
	user2, err := repo.GetByIDCached(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get cached user: %v", err)
	}

	// Verify both calls return same data
	if user1.Email != user2.Email {
		t.Errorf("Cache returned different data")
	}

	// Verify cache was populated
	cacheKey := "user:1"
	exists, err := cachedTestRedis.Exists(ctx, cacheKey).Result()
	if err != nil {
		t.Fatalf("Failed to check cache: %v", err)
	}
	if exists == 0 {
		t.Error("Expected cache to be populated")
	}
}

func TestCacheInvalidation(t *testing.T) {
	ctx := context.Background()
	repo := NewCachedUserRepository(cachedTestDB, cachedTestRedis)

	// Create and cache a user
	user, err := repo.CreateCached(ctx, "cache-test@example.com", "Cache Test")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer repo.DeleteCached(ctx, user.ID)

	// Get user to populate cache
	_, err = repo.GetByIDCached(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Update user - should invalidate cache
	err = repo.UpdateCached(ctx, user.ID, "updated@example.com", "Updated Name")
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify cache was invalidated
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	exists, err := cachedTestRedis.Exists(ctx, cacheKey).Result()
	if err != nil {
		t.Fatalf("Failed to check cache: %v", err)
	}
	if exists != 0 {
		t.Error("Expected cache to be invalidated after update")
	}
}

func TestTTLVerification(t *testing.T) {
	ctx := context.Background()
	repo := NewCachedUserRepository(cachedTestDB, cachedTestRedis)

	// Clear cache
	cachedTestRedis.FlushAll(ctx)

	// Get user to populate cache
	_, err := repo.GetByIDCached(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Check TTL is set
	cacheKey := "user:1"
	ttl, err := cachedTestRedis.TTL(ctx, cacheKey).Result()
	if err != nil {
		t.Fatalf("Failed to get TTL: %v", err)
	}

	if ttl <= 0 {
		t.Error("Expected positive TTL on cached item")
	}

	// TTL should be around 5 minutes (300 seconds)
	if ttl > 6*time.Minute || ttl < 4*time.Minute {
		t.Errorf("Expected TTL around 5 minutes, got: %v", ttl)
	}
}
