package gounity

import (
	"github.com/sirupsen/logrus"
)

// Options defines the mapping of optional parameters.
type Options map[string]interface{}

// Option defines an optional parameter.
type Option func(*Options)

// NewOptions constructs the optional parameters.
func NewOptions(opts ...Option) *Options {
	res := &Options{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// WarnNotUsedOptions logs warning messages for not-used optional parameters.
func (o *Options) WarnNotUsedOptions() {
	for input, value := range *o {
		logrus.Warnf("argument (%s:%v) ignored", input, value)
	}
}

func (o *Options) push(key string, value interface{}) {
	logrus.Debugf("argument (%s:%v) input", key, value)
	(*o)[key] = value
}

func (o *Options) pop(key string) interface{} {
	if res, ok := (*o)[key]; ok {
		logrus.Debugf("argument (%s:%v) used", key, res)
		delete(*o, key)
		return res
	}
	return nil
}

// PushName adds optional parameter `name`.
func (o *Options) PushName(name string) {
	o.push("name", name)
}

// PopName retrieves optional parameter `name` and removes it from options.
func (o *Options) PopName() interface{} {
	return o.pop("name")
}

// NameOpt constructs optional parameter `name`.
func NameOpt(name string) Option {
	return func(o *Options) {
		o.PushName(name)
	}
}

// PushSize adds optional parameter `size`.
func (o *Options) PushSize(sizeBytes uint64) {
	o.push("size", sizeBytes)
}

// PopSize retrieves optional parameter `size` and removes it from options.
func (o *Options) PopSize() interface{} {
	return o.pop("size")
}

// SizeGBOpt constructs optional parameter `size` in GB.
func SizeGBOpt(sizeGB uint64) Option {
	return func(o *Options) {
		o.PushSize(gbToBytes(sizeGB))
	}
}

// PushHostAccess adds optional parameter `hostAccess`.
func (o *Options) PushHostAccess(host *Host, accessMask HostLUNAccessEnum) {
	allHostAccess := []interface{}{}
	existing := o.pop("hostAccess")
	if existing != nil {
		allHostAccess = existing.([]interface{})
	}

	allHostAccess = append(
		allHostAccess,
		map[string]interface{}{
			"host":       host.Repr(),
			"accessMask": accessMask,
		},
	)
	o.push("hostAccess", allHostAccess)
}

// PopHostAccess retrieves optional parameter `hostAccess` and removes it from options.
func (o *Options) PopHostAccess() interface{} {
	return o.pop("hostAccess")
}

// HostAccessOpt constructs optional parameter `hostAccess`.
func HostAccessOpt(host *Host, accessMask HostLUNAccessEnum) Option {
	return func(o *Options) {
		o.PushHostAccess(host, accessMask)
	}
}

// PushDefaultAccess adds optional parameter `defaultAccess`.
func (o *Options) PushDefaultAccess(da NFSShareDefaultAccessEnum) {
	o.push("defaultAccess", da)
}

// PopDefaultAccess retrieves optional parameter `defaultAccess` and removes it from options.
func (o *Options) PopDefaultAccess() interface{} {
	return o.pop("defaultAccess")
}

// DefaultAccessOpt constructs optional parameter `defaultAccess`.
func DefaultAccessOpt(da NFSShareDefaultAccessEnum) Option {
	return func(o *Options) {
		o.PushDefaultAccess(da)
	}
}

// NewLunParameters constructs `lunParameters` used in `CreateLun`, .etc.
func (o *Options) NewLunParameters(p *Pool) map[string]interface{} {

	lunParams := map[string]interface{}{"pool": p.Repr()}
	if size := o.PopSize(); size != nil {
		lunParams["size"] = size
	}
	return lunParams
}
