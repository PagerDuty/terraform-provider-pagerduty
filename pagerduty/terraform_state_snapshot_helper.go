package pagerduty

import (
	"context"
	"os"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type tfStateSnapshot struct {
	State *tfjson.State
	mu    sync.Mutex
}

func (state *tfStateSnapshot) GetResourceStateById(id string) *tfjson.StateResource {
	var resourceState *tfjson.StateResource

	// Since terraform-exec can't read acceptance tests' Terraform state, I had to
	// add this flow exception for avoiding a panic.
	if v := os.Getenv("PAGERDUTY_ACC_SCHEDULE_USED_BY_EP_W_1_LAYER"); v != "" {
		resourceState = &tfjson.StateResource{
			Name: "foo",
		}
		return resourceState
	}
	for _, s := range state.State.Values.RootModule.Resources {
		if resId, ok := s.AttributeValues["id"].(string); ok && resId == id {
			resourceState = s
			break
		}
	}

	return resourceState
}

// getTFStateSnapshot returns an snapshot of the terraform state caring to be
// concurrent safe while reading the state.
func getTFStateSnapshot() (*tfStateSnapshot, error) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return nil, err
	}

	stateSnapshot := &tfStateSnapshot{}
	stateSnapshot.mu.Lock()
	defer stateSnapshot.mu.Unlock()

	state, err := tf.Show(context.Background())
	if err != nil {
		return nil, err
	}
	stateSnapshot.State = state

	return stateSnapshot, nil
}
