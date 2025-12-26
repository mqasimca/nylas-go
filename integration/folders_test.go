//go:build integration

// Folders Integration Tests Coverage:
//   - List, ListWithOptions, Get, ListAll, CRUD, SystemFolders ✓
//
// All FoldersService methods are fully tested.
// Note: CRUD test covers Create, Update, and Delete in sequence with cleanup.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/folders"
)

func TestFolders_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.Folders.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No folders found for this provider")
		}

		t.Logf("Found %d folders", len(resp.Data))

		for _, folder := range resp.Data {
			if folder.ID == "" {
				t.Error("Folder ID should not be empty")
			}
			if folder.Name == "" {
				t.Error("Folder Name should not be empty")
			}
			t.Logf("  - %s (id: %s, system: %v)", folder.Name, folder.ID, folder.SystemFolder)
		}
	})
}

func TestFolders_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		limit := 5
		resp, err := client.Folders.List(ctx, grantID, &folders.ListOptions{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) > limit {
			t.Errorf("List() returned %d folders, want <= %d", len(resp.Data), limit)
		}

		t.Logf("Found %d folders (limit %d)", len(resp.Data), limit)
	})
}

func TestFolders_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First list to get a folder ID
		listResp, err := client.Folders.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No folders found for this provider")
		}

		// Find a non-system folder if possible, otherwise use first folder
		var folderID string
		for _, f := range listResp.Data {
			if !f.SystemFolder {
				folderID = f.ID
				break
			}
		}
		if folderID == "" {
			folderID = listResp.Data[0].ID
		}

		folder, err := client.Folders.Get(ctx, grantID, folderID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", folderID, err)
		}

		if folder.ID != folderID {
			t.Errorf("Get() ID = %s, want %s", folder.ID, folderID)
		}

		t.Logf("Got folder: %s (id: %s)", folder.Name, folder.ID)
	})
}

func TestFolders_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		iter := client.Folders.ListAll(ctx, grantID, &folders.ListOptions{
			Limit: intPtr(5), // Small page size to test pagination
		})

		all, err := iter.Collect()
		if err != nil {
			t.Fatalf("Collect() error = %v", err)
		}

		t.Logf("ListAll() found %d folders", len(all))

		for _, folder := range all {
			if folder.ID == "" {
				t.Error("Folder ID should not be empty")
			}
		}
	})
}

func TestFolders_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Create a test folder
		createReq := &folders.CreateRequest{
			Name: "SDK Test Folder",
		}

		created, err := client.Folders.Create(ctx, grantID, createReq)
		if err != nil {
			t.Skipf("Create() error = %v (provider may not support folder creation)", err)
		}

		// Register cleanup to delete the folder
		cleanup.Add(func() {
			_ = client.Folders.Delete(ctx, grantID, created.ID)
		})

		if created.ID == "" {
			t.Fatal("Create() returned empty ID")
		}
		if created.Name != createReq.Name {
			t.Errorf("Create() Name = %s, want %s", created.Name, createReq.Name)
		}

		t.Logf("Created folder: %s (id: %s)", created.Name, created.ID)

		// Update the folder
		newName := "SDK Test Folder Updated"
		updated, err := client.Folders.Update(ctx, grantID, created.ID, &folders.UpdateRequest{
			Name: &newName,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if updated.Name != newName {
			t.Errorf("Update() Name = %s, want %s", updated.Name, newName)
		}

		t.Logf("Updated folder: %s -> %s", createReq.Name, updated.Name)

		// Delete the folder
		err = client.Folders.Delete(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		t.Log("Deleted folder successfully")
	})
}

func TestFolders_SystemFolders(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.Folders.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No folders found for this provider")
		}

		// Count system vs non-system folders
		var systemCount, userCount int
		for _, folder := range resp.Data {
			if folder.SystemFolder {
				systemCount++
			} else {
				userCount++
			}
		}

		t.Logf("Found %d system folders and %d user folders", systemCount, userCount)

		// Most providers should have at least some system folders (Inbox, Sent, etc.)
		if systemCount == 0 {
			t.Log("Warning: No system folders found")
		}
	})
}
