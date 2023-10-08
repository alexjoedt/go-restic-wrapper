package backup

type OptionFunc func(opts *options)

type options struct {
	host    string
	path    string
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

func WithHost(host string) OptionFunc {
	return func(opts *options) {
		opts.host = host
	}
}

func WithPath(path string) OptionFunc {
	return func(opts *options) {
		opts.path = path
	}
}

func (opts options) args() []string {
	args := make([]string, 0)

	if opts.host != "" {
		args = append(args, "--host", opts.host)
	}

	for _, t := range opts.tags {
		args = append(args, "--tag", t)
	}

	for _, exclude := range opts.exclude {
		args = append(args, "--exclude", exclude)
	}

	return args
}
