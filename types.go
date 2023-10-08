/*
BSD 2-Clause License

Copyright (c) 2014, Alexander Neumann <alexander@bumpern.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package restic

type BackupSummary struct {
	MessageType         string  `json:"message_type"`
	FilesNew            int     `json:"files_new"`
	FilesChanged        int     `json:"files_changed"`
	FilesUnmodified     int     `json:"files_unmodified"`
	DirsNew             int     `json:"dirs_new"`
	DirsChanged         int     `json:"dirs_changed"`
	DirsUnmodified      int     `json:"dirs_unmodified"`
	DataBlobs           int     `json:"data_blobs"`
	TreeBlobs           int     `json:"tree_blobs"`
	DataAdded           int     `json:"data_added"`
	TotalFilesProcessed int     `json:"total_files_processed"`
	TotalBytesProcessed int     `json:"total_bytes_processed"`
	TotalDuration       float64 `json:"total_duration"`
	SnapshotID          string  `json:"snapshot_id"`
}

type RestoreSummary struct {
	MessageType   string `json:"message_type"`
	TotalFiles    int    `json:"total_files"`
	FilesRestored int    `json:"files_restored"`
	TotalBytes    int    `json:"total_bytes"`
	BytesRestored int    `json:"bytes_restored"`
}

type ForgetSummary struct {
	Tags  []string `json:"tags"`
	Host  string   `json:"host"`
	Paths []string `json:"paths"`
	Keep  []struct {
		Time           string   `json:"time"`
		Parent         string   `json:"parent"`
		Tree           string   `json:"tree"`
		Paths          []string `json:"paths"`
		Hostname       string   `json:"hostname"`
		Username       string   `json:"username"`
		UID            int      `json:"uid"`
		GID            int      `json:"gid"`
		Tags           []string `json:"tags"`
		ProgramVersion string   `json:"program_version"`
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
	} `json:"keep"`
	Remove []struct {
		Time           string   `json:"time"`
		Parent         string   `json:"parent"`
		Tree           string   `json:"tree"`
		Paths          []string `json:"paths"`
		Hostname       string   `json:"hostname"`
		Username       string   `json:"username"`
		UID            int      `json:"uid"`
		GID            int      `json:"gid"`
		Tags           []string `json:"tags"`
		ProgramVersion string   `json:"program_version"`
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
	} `json:"remove"`
	Reasons []struct {
		Snapshot struct {
			Time           string   `json:"time"`
			Parent         string   `json:"parent"`
			Tree           string   `json:"tree"`
			Paths          []string `json:"paths"`
			Hostname       string   `json:"hostname"`
			Username       string   `json:"username"`
			UID            int      `json:"uid"`
			GID            int      `json:"gid"`
			Tags           []string `json:"tags"`
			ProgramVersion string   `json:"program_version"`
		} `json:"snapshot"`
		Matches  []string `json:"matches"`
		Counters struct {
			Last int `json:"last"`
		} `json:"counters"`
	} `json:"reasons"`
}
