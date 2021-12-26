package coinbasepro_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/moonr-app/go-coinbasepro/v2"
)

func TestGetProfiles(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	profiles, err := client.GetProfiles(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range profiles {
		if p.ID == "" {
			t.Fatal("profile id missing")
		}
	}
}

func TestGetProfile(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()

	profiles, err := client.GetProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}

	profile, err := client.GetProfile(ctx, profiles[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if profile.ID == "" {
		t.Fatal("profile id missing")
	}
}

func TestCreateProfileTransfer(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()

	profiles, err := client.GetProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var fromProfile, toProfile coinbasepro.Profile
	for _, profile := range profiles {
		if profile.IsDefault {
			fromProfile = profile
			continue
		}

		if profile.Active {
			toProfile = profile
			break
		}
	}

	if toProfile.ID == "" {
		t.Skip(fmt.Sprintf("needed at least two active profiles for this test"))
	}

	// Send from first profile to second profile
	newTransfer := coinbasepro.ProfileTransfer{
		From:     fromProfile.ID,
		To:       toProfile.ID,
		Currency: "USD",
		Amount:   "9.99",
	}

	err = client.CreateProfileTransfer(ctx, &newTransfer)
	if err != nil {
		t.Fatal(err)
	}
}
