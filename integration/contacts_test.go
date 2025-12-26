//go:build integration

// Contacts Integration Tests Coverage:
//   - List, ListWithOptions, Get, ListAll, CRUD, ListGroups ✓
//
// All ContactsService methods are fully tested.
// Note: CRUD test covers Create, Update, and Delete in sequence with cleanup.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/contacts"
)

func TestContacts_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.Contacts.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		t.Logf("Found %d contacts", len(resp.Data))

		for _, contact := range resp.Data {
			if contact.ID == "" {
				t.Error("Contact ID should not be empty")
			}
		}
	})
}

func TestContacts_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		limit := 5
		resp, err := client.Contacts.List(ctx, grantID, &contacts.ListOptions{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) > limit {
			t.Errorf("List() returned %d contacts, want <= %d", len(resp.Data), limit)
		}

		t.Logf("Found %d contacts (limit %d)", len(resp.Data), limit)
	})
}

func TestContacts_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First list to get a contact ID
		listResp, err := client.Contacts.List(ctx, grantID, &contacts.ListOptions{
			Limit: intPtr(1),
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No contacts found for this provider")
		}

		contactID := listResp.Data[0].ID

		contact, err := client.Contacts.Get(ctx, grantID, contactID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", contactID, err)
		}

		if contact.ID != contactID {
			t.Errorf("Get() ID = %s, want %s", contact.ID, contactID)
		}

		t.Logf("Got contact: %s (id: %s)", contact.DisplayName, contact.ID)
	})
}

func TestContacts_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		iter := client.Contacts.ListAll(ctx, grantID, &contacts.ListOptions{
			Limit: intPtr(10), // Small page size to test pagination
		})

		// Collect up to 50 contacts to avoid long test times
		var all []*contacts.Contact
		count := 0
		for {
			contact, err := iter.Next()
			if err != nil {
				break
			}
			all = append(all, contact)
			count++
			if count >= 50 {
				break
			}
		}

		t.Logf("ListAll() found %d contacts", len(all))

		for _, contact := range all {
			if contact.ID == "" {
				t.Error("Contact ID should not be empty")
			}
		}
	})
}

func TestContacts_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Create a test contact
		createReq := &contacts.CreateRequest{
			GivenName:   "SDK",
			Surname:     "Test",
			DisplayName: "SDK Test Contact",
			CompanyName: "Nylas",
			JobTitle:    "Integration Test",
			Notes:       "Created by Nylas Go SDK integration tests",
			Emails: []contacts.Email{
				{Email: "sdk-test@test.nylas.com", Type: "work"},
			},
			PhoneNumbers: []contacts.PhoneNumber{
				{Number: "+1-555-0100", Type: "work"},
			},
		}

		created, err := client.Contacts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Skipf("Create() error = %v (provider may not support contact creation)", err)
		}

		// Register cleanup to delete the contact
		cleanup.Add(func() {
			_ = client.Contacts.Delete(ctx, grantID, created.ID)
		})

		if created.ID == "" {
			t.Fatal("Create() returned empty ID")
		}
		if created.GivenName != createReq.GivenName {
			t.Errorf("Create() GivenName = %s, want %s", created.GivenName, createReq.GivenName)
		}

		t.Logf("Created contact: %s (id: %s)", created.DisplayName, created.ID)

		// Update the contact
		newJobTitle := "Senior Integration Test"
		updated, err := client.Contacts.Update(ctx, grantID, created.ID, &contacts.UpdateRequest{
			JobTitle: &newJobTitle,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if updated.JobTitle != newJobTitle {
			t.Errorf("Update() JobTitle = %s, want %s", updated.JobTitle, newJobTitle)
		}

		t.Logf("Updated contact: %s -> %s", createReq.JobTitle, updated.JobTitle)

		// Delete the contact
		err = client.Contacts.Delete(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		t.Log("Deleted contact successfully")
	})
}

func TestContacts_ListGroups(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		groups, err := client.Contacts.ListGroups(ctx, grantID)
		if err != nil {
			t.Skipf("ListGroups() error = %v (provider may not support contact groups)", err)
		}

		t.Logf("Found %d contact groups", len(groups))

		for _, group := range groups {
			if group.ID == "" {
				t.Error("Group ID should not be empty")
			}
			t.Logf("  - %s (id: %s)", group.Name, group.ID)
		}
	})
}
