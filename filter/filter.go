package filter

import "fmt"

type OptionFunc func(opts *options)

type options struct {
	hosts  []string
	paths  []string
	tags   []string
	latest uint
}

func Args(opts ...OptionFunc) []string {
	var options options
	for _, opt := range opts {
		opt(&options)
	}

	return options.args()
}

func WithTags(tags ...string) OptionFunc {
	return func(opts *options) {
		opts.tags = append(opts.tags, tags...)
	}
}

func WithHosts(hosts ...string) OptionFunc {
	return func(opts *options) {
		opts.hosts = append(opts.hosts, hosts...)
	}
}

func WithPaths(paths ...string) OptionFunc {
	return func(opts *options) {
		opts.paths = append(opts.paths, paths...)
	}
}

func WithLatest(no uint) OptionFunc {
	return func(opts *options) {
		opts.latest = no
	}
}

func (opts options) args() []string {
	args := make([]string, 0)

	for _, h := range opts.hosts {
		args = append(args, "--host", h)
	}

	for _, p := range opts.paths {
		args = append(args, "--path", p)
	}

	for _, t := range opts.tags {
		args = append(args, "--tag", t)
	}

	if opts.latest > 0 {
		args = append(args, "--latest", fmt.Sprintf("%d", opts.latest))
	}

	return args
}
