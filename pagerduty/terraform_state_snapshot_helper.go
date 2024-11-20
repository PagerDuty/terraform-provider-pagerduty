package pagerduty

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
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
	ctx := context.Background()

	installer := install.NewInstaller()
	defer installer.Remove(ctx)

	tfVersionConstrain := ">= 1.1.0"
	log.Printf("[pagerduty] Ensuring terraform %q is installed...", tfVersionConstrain)
	execPath, err := installer.Ensure(ctx, []src.Source{
		&fs.Version{
			Product:     product.Terraform,
			Constraints: version.MustConstraints(version.NewConstraint(tfVersionConstrain)),
		},
	})

	if err != nil {
		isTFVersionUnavailableError := strings.Contains(err.Error(), "terraform: executable file not found in $PATH")
		if !isTFVersionUnavailableError {
			return nil, err
		}
		installVersionString := "1.9.8" // latest at the time of writing
		installTargetVersion := version.Must(version.NewVersion(installVersionString))
		log.Printf("[pagerduty] Unable to locate terraform binary matching %q in $PATH, installing %q", tfVersionConstrain, installVersionString)
		execPathInstalled, installError := installer.Ensure(context.Background(), []src.Source{
			&fs.ExactVersion{
				Product: product.Terraform,
				Version: installTargetVersion,
			},
			&releases.ExactVersion{
				Product: product.Terraform,
				Version: installTargetVersion,
			},
		})
		if installError != nil {
			return nil, fmt.Errorf("[pagerduty] Failed to install terraform %q", installVersionString)
		}
		log.Printf("[pagerduty] Successfully installed to %q", execPathInstalled)
		execPath = execPathInstalled
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
