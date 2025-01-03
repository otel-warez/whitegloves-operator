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

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	networkingv1 "k8s.io/api/networking/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	k8sapiflag "k8s.io/component-base/cli/flag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/otel-warez/whitegloves-operator/internal/config"
	"github.com/otel-warez/whitegloves-operator/internal/constants"
	"github.com/otel-warez/whitegloves-operator/internal/instrumentation"
	"github.com/otel-warez/whitegloves-operator/internal/version"
	"github.com/otel-warez/whitegloves-operator/internal/webhook/podmutation"
)

var (
	scheme   = k8sruntime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

type tlsConfig struct {
	minVersion   string
	cipherSuites []string
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(networkingv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

// stringFlagOrEnv defines a string flag which can be set by an environment variable.
// Precedence: flag > env var > default value.
func stringFlagOrEnv(p *string, name string, envName string, defaultValue string, usage string) {
	envValue := os.Getenv(envName)
	if envValue != "" {
		defaultValue = envValue
	}
	pflag.StringVar(p, name, defaultValue, usage)
}

func main() {
	// registers any flags that underlying libraries might use
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	v := version.Get()

	// add flags related to this operator
	var (
		metricsAddr                              string
		probeAddr                                string
		pprofAddr                                string
		enableLeaderElection                     bool
		createRBACPermissions                    bool
		createOpenShiftDashboard                 bool
		enableMultiInstrumentation               bool
		enableApacheHttpdInstrumentation         bool
		enableDotNetInstrumentation              bool
		enableGoInstrumentation                  bool
		enablePythonInstrumentation              bool
		enableNginxInstrumentation               bool
		enableNodeJSInstrumentation              bool
		enableJavaInstrumentation                bool
		enableCRMetrics                          bool
		autoInstrumentationJava                  string
		autoInstrumentationNodeJS                string
		autoInstrumentationPython                string
		autoInstrumentationDotNet                string
		autoInstrumentationApacheHttpd           string
		autoInstrumentationNginx                 string
		autoInstrumentationGo                    string
		autoInstrumentationOtelCollectorEndpoint string
		labelsFilter                             []string
		tmplabelsFilter                          []string
		annotationsFilter                        []string
		webhookPort                              int
		tlsOpt                                   tlsConfig
		encodeMessageKey                         string
		encodeLevelKey                           string
		encodeTimeKey                            string
		encodeLevelFormat                        string
		fipsDisabledComponents                   string
	)

	pflag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	pflag.StringVar(&probeAddr, "health-probe-addr", ":8081", "The address the probe endpoint binds to.")
	pflag.StringVar(&pprofAddr, "pprof-addr", "", "The address to expose the pprof server. Default is empty string which disables the pprof server.")
	pflag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	pflag.BoolVar(&createRBACPermissions, "create-rbac-permissions", false, "Automatically create RBAC permissions needed by the processors (deprecated)")
	pflag.BoolVar(&createOpenShiftDashboard, "openshift-create-dashboard", false, "Create an OpenShift dashboard for monitoring the OpenTelemetryCollector instances")
	pflag.BoolVar(&enableMultiInstrumentation, "enable-multi-instrumentation", true, "Controls whether the operator supports multi instrumentation")
	pflag.BoolVar(&enableApacheHttpdInstrumentation, constants.FlagApacheHttpd, true, "Controls whether the operator supports Apache HTTPD auto-instrumentation")
	pflag.BoolVar(&enableDotNetInstrumentation, constants.FlagDotNet, true, "Controls whether the operator supports dotnet auto-instrumentation")
	pflag.BoolVar(&enableGoInstrumentation, constants.FlagGo, false, "Controls whether the operator supports Go auto-instrumentation")
	pflag.BoolVar(&enablePythonInstrumentation, constants.FlagPython, true, "Controls whether the operator supports python auto-instrumentation")
	pflag.BoolVar(&enableNginxInstrumentation, constants.FlagNginx, false, "Controls whether the operator supports nginx auto-instrumentation")
	pflag.BoolVar(&enableNodeJSInstrumentation, constants.FlagNodeJS, true, "Controls whether the operator supports nodejs auto-instrumentation")
	pflag.BoolVar(&enableJavaInstrumentation, constants.FlagJava, true, "Controls whether the operator supports java auto-instrumentation")
	pflag.BoolVar(&enableCRMetrics, constants.FlagCRMetrics, false, "Controls whether exposing the CR metrics is enabled")

	stringFlagOrEnv(&autoInstrumentationJava, "auto-instrumentation-java-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_JAVA", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-java:%s", v.AutoInstrumentationJava), "The default OpenTelemetry Java instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationNodeJS, "auto-instrumentation-nodejs-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_NODEJS", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-nodejs:%s", v.AutoInstrumentationNodeJS), "The default OpenTelemetry NodeJS instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationPython, "auto-instrumentation-python-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_PYTHON", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-python:%s", v.AutoInstrumentationPython), "The default OpenTelemetry Python instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationDotNet, "auto-instrumentation-dotnet-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_DOTNET", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-dotnet:%s", v.AutoInstrumentationDotNet), "The default OpenTelemetry DotNet instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationGo, "auto-instrumentation-go-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_GO", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-go-instrumentation/autoinstrumentation-go:%s", v.AutoInstrumentationGo), "The default OpenTelemetry Go instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationApacheHttpd, "auto-instrumentation-apache-httpd-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_APACHE_HTTPD", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-apache-httpd:%s", v.AutoInstrumentationApacheHttpd), "The default OpenTelemetry Apache HTTPD instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationNginx, "auto-instrumentation-nginx-image", "RELATED_IMAGE_AUTO_INSTRUMENTATION_NGINX", fmt.Sprintf("ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-apache-httpd:%s", v.AutoInstrumentationNginx), "The default OpenTelemetry Nginx instrumentation image. This image is used when no image is specified in the CustomResource.")
	stringFlagOrEnv(&autoInstrumentationOtelCollectorEndpoint, "auto-instrumentation-otel-collector-endpoint", "RELATED_IMAGE_AUTO_INSTRUMENTATION_OTEL_COLLECTOR_ENDPOINT", "https://0.0.0.0:4318", "The OpenTelemetry Collector exporter URL.")
	pflag.StringArrayVar(&labelsFilter, "labels-filter", []string{}, "Labels to filter away from propagating onto deploys. It should be a string array containing patterns, which are literal strings optionally containing a * wildcard character. Example: --labels-filter=.*filter.out will filter out labels that looks like: label.filter.out: true")
	pflag.StringArrayVar(&tmplabelsFilter, "label", []string{}, "(Deprecated, please use the labels-filter flag) Labels to filter away from propagating onto deploys. It should be a string array containing patterns, which are literal strings optionally containing a * wildcard character. Example: --label=.*filter.out will filter out labels that looks like: label.filter.out: true")
	pflag.StringArrayVar(&annotationsFilter, "annotations-filter", []string{}, "Annotations to filter away from propagating onto deploys. It should be a string array containing patterns, which are literal strings optionally containing a * wildcard character. Example: --annotations-filter=.*filter.out will filter out annotations that looks like: annotation.filter.out: true")
	pflag.StringVar(&tlsOpt.minVersion, "tls-min-version", "VersionTLS12", "Minimum TLS version supported. Value must match version names from https://golang.org/pkg/crypto/tls/#pkg-constants.")
	pflag.StringSliceVar(&tlsOpt.cipherSuites, "tls-cipher-suites", nil, "Comma-separated list of cipher suites for the server. Values are from tls package constants (https://golang.org/pkg/crypto/tls/#pkg-constants). If omitted, the default Go cipher suites will be used")
	pflag.StringVar(&encodeMessageKey, "zap-message-key", "message", "The message key to be used in the customized Log Encoder")
	pflag.StringVar(&encodeLevelKey, "zap-level-key", "level", "The level key to be used in the customized Log Encoder")
	pflag.StringVar(&encodeTimeKey, "zap-time-key", "timestamp", "The time key to be used in the customized Log Encoder")
	pflag.StringVar(&encodeLevelFormat, "zap-level-format", "uppercase", "The level format to be used in the customized Log Encoder")
	pflag.StringVar(&fipsDisabledComponents, "fips-disabled-components", "uppercase", "Disabled collector components when operator runs on FIPS enabled platform. Example flag value =receiver.foo,receiver.bar,exporter.baz")
	pflag.IntVar(&webhookPort, "webhook-port", 9443, "The port the webhook endpoint binds to.")
	pflag.Parse()

	// Using labelfilters both from label and labels-filter flags, until label flag is removed
	labelsFilter = append(labelsFilter, tmplabelsFilter...)

	opts.EncoderConfigOptions = append(opts.EncoderConfigOptions, func(ec *zapcore.EncoderConfig) {
		ec.MessageKey = encodeMessageKey
		ec.LevelKey = encodeLevelKey
		ec.TimeKey = encodeTimeKey
		ec.EncodeLevel = config.WithEncodeLevelFormat(encodeLevelFormat)
	})

	logger := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)

	logger.Info("Starting the OpenTelemetry Operator",
		"opentelemetry-operator", v.Operator,
		"auto-instrumentation-java", autoInstrumentationJava,
		"auto-instrumentation-nodejs", autoInstrumentationNodeJS,
		"auto-instrumentation-python", autoInstrumentationPython,
		"auto-instrumentation-dotnet", autoInstrumentationDotNet,
		"auto-instrumentation-go", autoInstrumentationGo,
		"auto-instrumentation-apache-httpd", autoInstrumentationApacheHttpd,
		"auto-instrumentation-nginx", autoInstrumentationNginx,
		"build-date", v.BuildDate,
		"go-version", v.Go,
		"go-arch", runtime.GOARCH,
		"go-os", runtime.GOOS,
		"labels-filter", labelsFilter,
		"annotations-filter", annotationsFilter,
		"enable-multi-instrumentation", enableMultiInstrumentation,
		"enable-apache-httpd-instrumentation", enableApacheHttpdInstrumentation,
		"enable-dotnet-instrumentation", enableDotNetInstrumentation,
		"enable-go-instrumentation", enableGoInstrumentation,
		"enable-python-instrumentation", enablePythonInstrumentation,
		"enable-nginx-instrumentation", enableNginxInstrumentation,
		"enable-nodejs-instrumentation", enableNodeJSInstrumentation,
		"enable-java-instrumentation", enableJavaInstrumentation,
		"create-openshift-dashboard", createOpenShiftDashboard,
		"zap-message-key", encodeMessageKey,
		"zap-level-key", encodeLevelKey,
		"zap-time-key", encodeTimeKey,
		"zap-level-format", encodeLevelFormat,
	)

	restConfig := ctrl.GetConfigOrDie()

	var namespaces map[string]cache.Config
	watchNamespace, found := os.LookupEnv("WATCH_NAMESPACE")
	if found {
		setupLog.Info("watching namespace(s)", "namespaces", watchNamespace)
		namespaces = map[string]cache.Config{}
		for _, ns := range strings.Split(watchNamespace, ",") {
			namespaces[ns] = cache.Config{}
		}
	} else {
		setupLog.Info("the env var WATCH_NAMESPACE isn't set, watching all namespaces")
	}

	// see https://github.com/openshift/library-go/blob/4362aa519714a4b62b00ab8318197ba2bba51cb7/pkg/config/leaderelection/leaderelection.go#L104
	leaseDuration := time.Second * 137
	renewDeadline := time.Second * 107
	retryPeriod := time.Second * 26

	optionsTlSOptsFuncs := []func(*tls.Config){
		func(config *tls.Config) { tlsConfigSetting(config, tlsOpt) },
	}

	mgrOptions := ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress:        probeAddr,
		LeaderElection:                enableLeaderElection,
		LeaderElectionID:              "9f7554c3.opentelemetry.io",
		LeaderElectionReleaseOnCancel: true,
		LeaseDuration:                 &leaseDuration,
		RenewDeadline:                 &renewDeadline,
		RetryPeriod:                   &retryPeriod,
		PprofBindAddress:              pprofAddr,
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    webhookPort,
			TLSOpts: optionsTlSOptsFuncs,
		}),
		Cache: cache.Options{
			DefaultNamespaces: namespaces,
		},
	}

	mgr, err := ctrl.NewManager(restConfig, mgrOptions)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := ctrl.SetupSignalHandler()

	cfg := config.New(
		config.WithLogger(ctrl.Log.WithName("config")),
		config.WithVersion(v),
		config.WithEnableMultiInstrumentation(enableMultiInstrumentation),
		config.WithEnableApacheHttpdInstrumentation(enableApacheHttpdInstrumentation),
		config.WithEnableDotNetInstrumentation(enableDotNetInstrumentation),
		config.WithEnableGoInstrumentation(enableGoInstrumentation),
		config.WithEnableNginxInstrumentation(enableNginxInstrumentation),
		config.WithEnablePythonInstrumentation(enablePythonInstrumentation),
		config.WithEnableNodeJSInstrumentation(enableNodeJSInstrumentation),
		config.WithEnableJavaInstrumentation(enableJavaInstrumentation),
		config.WithAutoInstrumentationJavaImage(autoInstrumentationJava),
		config.WithAutoInstrumentationNodeJSImage(autoInstrumentationNodeJS),
		config.WithAutoInstrumentationPythonImage(autoInstrumentationPython),
		config.WithAutoInstrumentationDotNetImage(autoInstrumentationDotNet),
		config.WithAutoInstrumentationOtelCollectorEndpoint(autoInstrumentationOtelCollectorEndpoint),
		config.WithAutoInstrumentationGoImage(autoInstrumentationGo),
		config.WithAutoInstrumentationApacheHttpdImage(autoInstrumentationApacheHttpd),
		config.WithAutoInstrumentationNginxImage(autoInstrumentationNginx),
		config.WithLabelFilters(labelsFilter),
		config.WithAnnotationFilters(annotationsFilter),
	)
	if err != nil {
		setupLog.Error(err, "failed to autodetect config variables")
	}
	if cfg.AnnotationsFilter() != nil {
		for _, basePattern := range cfg.AnnotationsFilter() {
			_, compileErr := regexp.Compile(basePattern)
			if compileErr != nil {
				setupLog.Error(compileErr, "could not compile the regexp pattern for Annotations filter")
			}
		}
	}

	if cfg.LabelsFilter() != nil {
		for _, basePattern := range cfg.LabelsFilter() {
			_, compileErr := regexp.Compile(basePattern)
			if compileErr != nil {
				setupLog.Error(compileErr, "could not compile the regexp pattern for Labels filter")
			}
		}
	}

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		decoder := admission.NewDecoder(mgr.GetScheme())
		mgr.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{
			Handler: podmutation.NewWebhookHandler(cfg, ctrl.Log.WithName("pod-webhook"), decoder, mgr.GetClient(),
				[]podmutation.PodMutator{
					instrumentation.NewMutator(logger, mgr.GetClient(), mgr.GetEventRecorderFor("opentelemetry-operator"), cfg),
				}),
		})
	} else {
		ctrl.Log.Info("Webhooks are disabled, operator is running an unsupported mode", "ENABLE_WEBHOOKS", "false")
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	// NOTE: We enable LeaderElectionReleaseOnCancel, and to be safe we need to exit right after the manager does
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// This function get the option from command argument (tlsConfig), check the validity through k8sapiflag
// and set the config for webhook server.
// refer to https://pkg.go.dev/k8s.io/component-base/cli/flag
func tlsConfigSetting(cfg *tls.Config, tlsOpt tlsConfig) {
	// TLSVersion helper function returns the TLS Version ID for the version name passed.
	tlsVersion, err := k8sapiflag.TLSVersion(tlsOpt.minVersion)
	if err != nil {
		setupLog.Error(err, "TLS version invalid")
	}
	cfg.MinVersion = tlsVersion

	// TLSCipherSuites helper function returns a list of cipher suite IDs from the cipher suite names passed.
	cipherSuiteIDs, err := k8sapiflag.TLSCipherSuites(tlsOpt.cipherSuites)
	if err != nil {
		setupLog.Error(err, "Failed to convert TLS cipher suite name to ID")
	}
	cfg.CipherSuites = cipherSuiteIDs
}
