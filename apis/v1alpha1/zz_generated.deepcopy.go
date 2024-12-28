//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApacheHttpd) DeepCopyInto(out *ApacheHttpd) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Attrs != nil {
		in, out := &in.Attrs, &out.Attrs
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApacheHttpd.
func (in *ApacheHttpd) DeepCopy() *ApacheHttpd {
	if in == nil {
		return nil
	}
	out := new(ApacheHttpd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Defaults) DeepCopyInto(out *Defaults) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Defaults.
func (in *Defaults) DeepCopy() *Defaults {
	if in == nil {
		return nil
	}
	out := new(Defaults)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DotNet) DeepCopyInto(out *DotNet) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DotNet.
func (in *DotNet) DeepCopy() *DotNet {
	if in == nil {
		return nil
	}
	out := new(DotNet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Exporter) DeepCopyInto(out *Exporter) {
	*out = *in
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(TLS)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Exporter.
func (in *Exporter) DeepCopy() *Exporter {
	if in == nil {
		return nil
	}
	out := new(Exporter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Extensions) DeepCopyInto(out *Extensions) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extensions.
func (in *Extensions) DeepCopy() *Extensions {
	if in == nil {
		return nil
	}
	out := new(Extensions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Go) DeepCopyInto(out *Go) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Go.
func (in *Go) DeepCopy() *Go {
	if in == nil {
		return nil
	}
	out := new(Go)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Instrumentation) DeepCopyInto(out *Instrumentation) {
	*out = *in
	out.Status = in.Status
	out.TypeMeta = in.TypeMeta
	in.Spec.DeepCopyInto(&out.Spec)
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Instrumentation.
func (in *Instrumentation) DeepCopy() *Instrumentation {
	if in == nil {
		return nil
	}
	out := new(Instrumentation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Instrumentation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationList) DeepCopyInto(out *InstrumentationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Instrumentation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationList.
func (in *InstrumentationList) DeepCopy() *InstrumentationList {
	if in == nil {
		return nil
	}
	out := new(InstrumentationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstrumentationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationSpec) DeepCopyInto(out *InstrumentationSpec) {
	*out = *in
	in.Exporter.DeepCopyInto(&out.Exporter)
	in.Resource.DeepCopyInto(&out.Resource)
	if in.Propagators != nil {
		in, out := &in.Propagators, &out.Propagators
		*out = make([]Propagator, len(*in))
		copy(*out, *in)
	}
	out.Sampler = in.Sampler
	out.Defaults = in.Defaults
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Java.DeepCopyInto(&out.Java)
	in.NodeJS.DeepCopyInto(&out.NodeJS)
	in.Python.DeepCopyInto(&out.Python)
	in.DotNet.DeepCopyInto(&out.DotNet)
	in.Go.DeepCopyInto(&out.Go)
	in.ApacheHttpd.DeepCopyInto(&out.ApacheHttpd)
	in.Nginx.DeepCopyInto(&out.Nginx)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationSpec.
func (in *InstrumentationSpec) DeepCopy() *InstrumentationSpec {
	if in == nil {
		return nil
	}
	out := new(InstrumentationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationStatus) DeepCopyInto(out *InstrumentationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationStatus.
func (in *InstrumentationStatus) DeepCopy() *InstrumentationStatus {
	if in == nil {
		return nil
	}
	out := new(InstrumentationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Java) DeepCopyInto(out *Java) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Extensions != nil {
		in, out := &in.Extensions, &out.Extensions
		*out = make([]Extensions, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Java.
func (in *Java) DeepCopy() *Java {
	if in == nil {
		return nil
	}
	out := new(Java)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nginx) DeepCopyInto(out *Nginx) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Attrs != nil {
		in, out := &in.Attrs, &out.Attrs
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nginx.
func (in *Nginx) DeepCopy() *Nginx {
	if in == nil {
		return nil
	}
	out := new(Nginx)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeJS) DeepCopyInto(out *NodeJS) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeJS.
func (in *NodeJS) DeepCopy() *NodeJS {
	if in == nil {
		return nil
	}
	out := new(NodeJS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Python) DeepCopyInto(out *Python) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.VolumeSizeLimit != nil {
		in, out := &in.VolumeSizeLimit, &out.VolumeSizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Python.
func (in *Python) DeepCopy() *Python {
	if in == nil {
		return nil
	}
	out := new(Python)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Sampler) DeepCopyInto(out *Sampler) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sampler.
func (in *Sampler) DeepCopy() *Sampler {
	if in == nil {
		return nil
	}
	out := new(Sampler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLS) DeepCopyInto(out *TLS) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLS.
func (in *TLS) DeepCopy() *TLS {
	if in == nil {
		return nil
	}
	out := new(TLS)
	in.DeepCopyInto(out)
	return out
}
