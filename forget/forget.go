package forget

import "fmt"

type OptionFunc func(opts *options)

type options struct {
	id       string
	hosts    []string
	paths    []string
	tags     []string
	prune    bool
	keepLast uint
}

func Args(opts ...OptionFunc) []string {
	var options options
	for _, opt := range opts {
		opt(&options)
	}

	return options.args()
}

func WithSnapshotID(id string) OptionFunc {
	return func(opts *options) {
		opts.id = id
	}
}

func WithPrune() OptionFunc {
	return func(opts *options) {
		opts.prune = true
	}
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

func WithKeepLast(no uint) OptionFunc {
	return func(opts *options) {
		opts.keepLast = no
	}
}

func (opts options) args() []string {
	args := make([]string, 0)

	// id must be the first arg after forget
	if opts.id != "" {
		args = append(args, opts.id)
	}

	for _, h := range opts.hosts {
		args = append(args, "--host", h)
	}

	for _, p := range opts.paths {
		args = append(args, "--path", p)
	}

	for _, t := range opts.tags {
		args = append(args, "--tag", t)
	}

	if opts.keepLast > 0 {
		args = append(args, "--keep-last", fmt.Sprintf("%d", opts.keepLast))
	}

	if opts.prune {
		args = append(args, "--prune")
	}

	return args
}
