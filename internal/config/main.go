// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config contains the operator's runtime configuration.
package config

import (
	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/otel-warez/whitegloves-operator/internal/version"
)

// Config holds the static configuration for this operator.
type Config struct {
	logger                                   logr.Logger
	autoInstrumentationPythonImage           string
	enableMultiInstrumentation               bool
	enableApacheHttpdInstrumentation         bool
	enableDotNetInstrumentation              bool
	enableGoInstrumentation                  bool
	enableNginxInstrumentation               bool
	enablePythonInstrumentation              bool
	enableNodeJSInstrumentation              bool
	enableJavaInstrumentation                bool
	autoInstrumentationDotNetImage           string
	autoInstrumentationGoImage               string
	autoInstrumentationApacheHttpdImage      string
	autoInstrumentationNginxImage            string
	targetAllocatorConfigMapEntry            string
	operatorOpAMPBridgeConfigMapEntry        string
	autoInstrumentationNodeJSImage           string
	autoInstrumentationJavaImage             string
	autoInstrumentationOtelCollectorEndpoint string

	labelsFilter      []string
	annotationsFilter []string
}

// New constructs a new configuration based on the given options.
func New(opts ...Option) Config {
	// initialize with the default values
	o := options{
		logger:                    logf.Log.WithName("config"),
		version:                   version.Get(),
		enableJavaInstrumentation: true,
		annotationsFilter:         []string{"kubectl.kubernetes.io/last-applied-configuration"},
	}

	for _, opt := range opts {
		opt(&o)
	}

	return Config{
		enableMultiInstrumentation:               o.enableMultiInstrumentation,
		enableApacheHttpdInstrumentation:         o.enableApacheHttpdInstrumentation,
		enableDotNetInstrumentation:              o.enableDotNetInstrumentation,
		enableGoInstrumentation:                  o.enableGoInstrumentation,
		enableNginxInstrumentation:               o.enableNginxInstrumentation,
		enablePythonInstrumentation:              o.enablePythonInstrumentation,
		enableNodeJSInstrumentation:              o.enableNodeJSInstrumentation,
		enableJavaInstrumentation:                o.enableJavaInstrumentation,
		logger:                                   o.logger,
		autoInstrumentationJavaImage:             o.autoInstrumentationJavaImage,
		autoInstrumentationNodeJSImage:           o.autoInstrumentationNodeJSImage,
		autoInstrumentationPythonImage:           o.autoInstrumentationPythonImage,
		autoInstrumentationDotNetImage:           o.autoInstrumentationDotNetImage,
		autoInstrumentationGoImage:               o.autoInstrumentationGoImage,
		autoInstrumentationApacheHttpdImage:      o.autoInstrumentationApacheHttpdImage,
		autoInstrumentationNginxImage:            o.autoInstrumentationNginxImage,
		autoInstrumentationOtelCollectorEndpoint: o.autoInstrumentationOtelCollectorEndpoint,
		labelsFilter:                             o.labelsFilter,
		annotationsFilter:                        o.annotationsFilter,
	}
}

// EnableMultiInstrumentation is true when the operator supports multi instrumentation.
func (c *Config) EnableMultiInstrumentation() bool {
	return c.enableMultiInstrumentation
}

// EnableApacheHttpdAutoInstrumentation is true when the operator supports ApacheHttpd auto instrumentation.
func (c *Config) EnableApacheHttpdAutoInstrumentation() bool {
	return c.enableApacheHttpdInstrumentation
}

// EnableDotNetAutoInstrumentation is true when the operator supports dotnet auto instrumentation.
func (c *Config) EnableDotNetAutoInstrumentation() bool {
	return c.enableDotNetInstrumentation
}

// EnableGoAutoInstrumentation is true when the operator supports Go auto instrumentation.
func (c *Config) EnableGoAutoInstrumentation() bool {
	return c.enableGoInstrumentation
}

// EnableNginxAutoInstrumentation is true when the operator supports nginx auto instrumentation.
func (c *Config) EnableNginxAutoInstrumentation() bool {
	return c.enableNginxInstrumentation
}

// EnableJavaAutoInstrumentation is true when the operator supports nginx auto instrumentation.
func (c *Config) EnableJavaAutoInstrumentation() bool {
	return c.enableJavaInstrumentation
}

// EnablePythonAutoInstrumentation is true when the operator supports dotnet auto instrumentation.
func (c *Config) EnablePythonAutoInstrumentation() bool {
	return c.enablePythonInstrumentation
}

// EnableNodeJSAutoInstrumentation is true when the operator supports dotnet auto instrumentation.
func (c *Config) EnableNodeJSAutoInstrumentation() bool {
	return c.enableNodeJSInstrumentation
}

// AutoInstrumentationJavaImage returns OpenTelemetry Java auto-instrumentation container image.
func (c *Config) AutoInstrumentationJavaImage() string {
	return c.autoInstrumentationJavaImage
}

func (c *Config) AutoInstrumentationOtelCollectorEndpoint() string {
	return c.autoInstrumentationOtelCollectorEndpoint
}

// AutoInstrumentationNodeJSImage returns OpenTelemetry NodeJS auto-instrumentation container image.
func (c *Config) AutoInstrumentationNodeJSImage() string {
	return c.autoInstrumentationNodeJSImage
}

// AutoInstrumentationPythonImage returns OpenTelemetry Python auto-instrumentation container image.
func (c *Config) AutoInstrumentationPythonImage() string {
	return c.autoInstrumentationPythonImage
}

// AutoInstrumentationDotNetImage returns OpenTelemetry DotNet auto-instrumentation container image.
func (c *Config) AutoInstrumentationDotNetImage() string {
	return c.autoInstrumentationDotNetImage
}

// AutoInstrumentationGoImage returns OpenTelemetry Go auto-instrumentation container image.
func (c *Config) AutoInstrumentationGoImage() string {
	return c.autoInstrumentationGoImage
}

// AutoInstrumentationApacheHttpdImage returns OpenTelemetry ApacheHttpd auto-instrumentation container image.
func (c *Config) AutoInstrumentationApacheHttpdImage() string {
	return c.autoInstrumentationApacheHttpdImage
}

// AutoInstrumentationNginxImage returns OpenTelemetry Nginx auto-instrumentation container image.
func (c *Config) AutoInstrumentationNginxImage() string {
	return c.autoInstrumentationNginxImage
}

// LabelsFilter Returns the filters converted to regex strings used to filter out unwanted labels from propagations.
func (c *Config) LabelsFilter() []string {
	return c.labelsFilter
}

// AnnotationsFilter Returns the filters converted to regex strings used to filter out unwanted labels from propagations.
func (c *Config) AnnotationsFilter() []string {
	return c.annotationsFilter
}
