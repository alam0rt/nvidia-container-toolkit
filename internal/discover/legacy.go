/**
# Copyright (c) 2022, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package discover

import (
	"github.com/NVIDIA/nvidia-container-toolkit/internal/lookup"
	"github.com/container-orchestrated-devices/container-device-interface/pkg/cdi"
	"github.com/sirupsen/logrus"
)

// NewLegacyDiscoverer creates a discoverer for the experimental runtime
func NewLegacyDiscoverer(logger *logrus.Logger, cfg *Config) (Discover, error) {
	d := legacy{
		logger: logger,
		lookup: lookup.NewExecutableLocator(logger, cfg.Root),
	}

	return &d, nil
}

type legacy struct {
	None
	logger *logrus.Logger
	lookup lookup.Locator
}

var _ Discover = (*legacy)(nil)

const (
	nvidiaContainerRuntimeHookExecutable      = "nvidia-container-runtime-hook"
	nvidiaContainerRuntimeHookDefaultFilePath = "/usr/bin/nvidia-container-runtime-hook"
)

// Hooks returns the "legacy" NVIDIA Container Runtime hook. This mirrors the behaviour of the stable
// modifier.
func (d legacy) Hooks() ([]Hook, error) {
	hookPath := nvidiaContainerRuntimeHookDefaultFilePath
	targets, err := d.lookup.Locate(nvidiaContainerRuntimeHookExecutable)
	if err != nil {
		d.logger.Warnf("Failed to locate %v: %v", nvidiaContainerRuntimeHookExecutable, err)
	} else if len(targets) == 0 {
		d.logger.Warnf("%v not found", nvidiaContainerRuntimeHookExecutable)
	} else {
		d.logger.Debugf("Found %v candidates: %v", nvidiaContainerRuntimeHookExecutable, targets)
		hookPath = targets[0]
	}
	d.logger.Debugf("Using NVIDIA Container Runtime Hook path %v", hookPath)

	args := []string{hookPath, "prestart"}
	h := Hook{
		Lifecycle: cdi.PrestartHook,
		Path:      hookPath,
		Args:      args,
	}

	return []Hook{h}, nil
}