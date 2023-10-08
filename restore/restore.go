package restore

type OptionFunc func(opts *options)

type options struct {
	hosts   []string
	paths   []string
	tags    []string
	exclude []string
	include []string
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

func WithIncludes(includes ...string) OptionFunc {
	return func(opts *options) {
		opts.include = append(opts.include, includes...)
	}
}

func WithExcludes(excludes ...string) OptionFunc {
	return func(opts *options) {
		opts.exclude = append(opts.exclude, excludes...)
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

	for _, exclude := range opts.exclude {
		args = append(args, "--exclude", exclude)
	}

	for _, include := range opts.include {
		args = append(args, "--include", include)
	}

	return args
}
