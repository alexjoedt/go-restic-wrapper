package restic

import "fmt"

// BackupOption configures a backup operation.
type BackupOption func(*backupOptions)

type backupOptions struct {
	host    string
	tags    []string
	exclude []string
	include []string
}

// WithHost sets the hostname for the backup.
func WithHost(host string) BackupOption {
	return func(opts *backupOptions) {
		opts.host = host
	}
}

// WithTags adds tags to the backup snapshot.
func WithTags(tags ...string) BackupOption {
	return func(opts *backupOptions) {
		opts.tags = append(opts.tags, tags...)
	}
}

// WithExclude adds patterns to exclude from the backup.
func WithExclude(patterns ...string) BackupOption {
	return func(opts *backupOptions) {
		opts.exclude = append(opts.exclude, patterns...)
	}
}

// WithInclude adds patterns to explicitly include in the backup.
func WithInclude(patterns ...string) BackupOption {
	return func(opts *backupOptions) {
		opts.include = append(opts.include, patterns...)
	}
}

func (opts backupOptions) args() []string {
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

	for _, include := range opts.include {
		args = append(args, "--include", include)
	}

	return args
}

// FilterOption configures snapshot filtering operations.
type FilterOption func(*filterOptions)

type filterOptions struct {
	hosts  []string
	paths  []string
	tags   []string
	latest uint
}

// FilterByHost filters snapshots by hostname.
func FilterByHost(hosts ...string) FilterOption {
	return func(opts *filterOptions) {
		opts.hosts = append(opts.hosts, hosts...)
	}
}

// FilterByPath filters snapshots by path.
func FilterByPath(paths ...string) FilterOption {
	return func(opts *filterOptions) {
		opts.paths = append(opts.paths, paths...)
	}
}

// FilterByTag filters snapshots by tag.
func FilterByTag(tags ...string) FilterOption {
	return func(opts *filterOptions) {
		opts.tags = append(opts.tags, tags...)
	}
}

// FilterLatest limits results to the latest N snapshots.
func FilterLatest(n uint) FilterOption {
	return func(opts *filterOptions) {
		opts.latest = n
	}
}

func (opts filterOptions) args() []string {
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

// ForgetOption configures a forget operation.
type ForgetOption func(*forgetOptions)

type forgetOptions struct {
	id       string
	hosts    []string
	paths    []string
	tags     []string
	prune    bool
	keepLast uint
}

// ForgetSnapshot specifies a snapshot ID to forget.
func ForgetSnapshot(id string) ForgetOption {
	return func(opts *forgetOptions) {
		opts.id = id
	}
}

// ForgetWithPrune enables automatic pruning after forget.
func ForgetWithPrune() ForgetOption {
	return func(opts *forgetOptions) {
		opts.prune = true
	}
}

// ForgetByHost filters snapshots to forget by hostname.
func ForgetByHost(hosts ...string) ForgetOption {
	return func(opts *forgetOptions) {
		opts.hosts = append(opts.hosts, hosts...)
	}
}

// ForgetByPath filters snapshots to forget by path.
func ForgetByPath(paths ...string) ForgetOption {
	return func(opts *forgetOptions) {
		opts.paths = append(opts.paths, paths...)
	}
}

// ForgetByTag filters snapshots to forget by tag.
func ForgetByTag(tags ...string) ForgetOption {
	return func(opts *forgetOptions) {
		opts.tags = append(opts.tags, tags...)
	}
}

// ForgetKeepLast keeps the last N snapshots.
func ForgetKeepLast(n uint) ForgetOption {
	return func(opts *forgetOptions) {
		opts.keepLast = n
	}
}

func (opts forgetOptions) args() []string {
	args := make([]string, 0)

	// Snapshot ID must be the first argument after "forget"
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

// RestoreOption configures a restore operation.
type RestoreOption func(*restoreOptions)

type restoreOptions struct {
	hosts   []string
	paths   []string
	tags    []string
	exclude []string
	include []string
}

// RestoreByHost filters restore by hostname.
func RestoreByHost(hosts ...string) RestoreOption {
	return func(opts *restoreOptions) {
		opts.hosts = append(opts.hosts, hosts...)
	}
}

// RestoreByPath filters restore by path.
func RestoreByPath(paths ...string) RestoreOption {
	return func(opts *restoreOptions) {
		opts.paths = append(opts.paths, paths...)
	}
}

// RestoreByTag filters restore by tag.
func RestoreByTag(tags ...string) RestoreOption {
	return func(opts *restoreOptions) {
		opts.tags = append(opts.tags, tags...)
	}
}

// RestoreExclude excludes patterns from restore.
func RestoreExclude(patterns ...string) RestoreOption {
	return func(opts *restoreOptions) {
		opts.exclude = append(opts.exclude, patterns...)
	}
}

// RestoreInclude includes patterns in restore.
func RestoreInclude(patterns ...string) RestoreOption {
	return func(opts *restoreOptions) {
		opts.include = append(opts.include, patterns...)
	}
}

func (opts restoreOptions) args() []string {
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
