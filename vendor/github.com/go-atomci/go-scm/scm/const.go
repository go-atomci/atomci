// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"encoding/json"
)

// State represents the commit state.
type State int

// State values.
const (
	StateUnknown State = iota
	StatePending
	StateRunning
	StateSuccess
	StateFailure
	StateCanceled
	StateError
)

// Action identifies webhook actions.
type Action int

// Action values.
const (
	ActionUnknown Action = iota
	ActionCreate
	ActionUpdate
	ActionDelete
	// issues
	ActionOpen
	ActionReopen
	ActionClose
	ActionLabel
	ActionUnlabel
	// pull requests
	ActionSync
	ActionMerge
	// issue comment
	ActionEdit
)

// String returns the string representation of Action.
func (a Action) String() (s string) {
	switch a {
	case ActionCreate:
		return "created"
	case ActionUpdate:
		return "updated"
	case ActionDelete:
		return "deleted"
	case ActionLabel:
		return "labeled"
	case ActionUnlabel:
		return "unlabeled"
	case ActionOpen:
		return "opened"
	case ActionReopen:
		return "reopened"
	case ActionClose:
		return "closed"
	case ActionSync:
		return "synchronized"
	case ActionMerge:
		return "merged"
	default:
		return
	}
}

// MarshalJSON returns the JSON-encoded Action.
func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON unmarshales the JSON-encoded Action.
func (a *Action) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "created":
		*a = ActionCreate
	case "updated":
		*a = ActionUpdate
	case "deleted":
		*a = ActionDelete
	case "labeled":
		*a = ActionLabel
	case "unlabeled":
		*a = ActionUnlabel
	case "opened":
		*a = ActionOpen
	case "reopened":
		*a = ActionReopen
	case "closed":
		*a = ActionClose
	case "synchronized":
		*a = ActionSync
	case "merged":
		*a = ActionMerge
	}
	return nil
}

// Driver identifies source code management driver.
type Driver int

// Driver values.
const (
	DriverUnknown Driver = iota
	DriverGithub
	DriverGitlab
	DriverGogs
	DriverGitea
	DriverBitbucket
	DriverStash
	DriverCoding
)

// String returns the string representation of Driver.
func (d Driver) String() (s string) {
	switch d {
	case DriverGithub:
		return "github"
	case DriverGitlab:
		return "gitlab"
	case DriverGogs:
		return "gogs"
	case DriverGitea:
		return "gitea"
	case DriverBitbucket:
		return "bitbucket"
	case DriverStash:
		return "stash"
	case DriverCoding:
		return "coding"
	default:
		return "unknown"
	}
}

// Role defines membership roles.
type Role int

// Role values.
const (
	RoleUndefined Role = iota
	RoleMember
	RoleAdmin
)

// String returns the string representation of Role.
func (r Role) String() (s string) {
	switch r {
	case RoleMember:
		return "member"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// ContentKind defines the kind of a content in a directory.
type ContentKind int

// ContentKind values.
const (
	ContentKindUnsupported ContentKind = iota
	ContentKindFile
	ContentKindDirectory
	ContentKindSymlink
	ContentKindGitlink
)

// String returns the string representation of ContentKind.
func (k ContentKind) String() string {
	switch k {
	case ContentKindFile:
		return "file"
	case ContentKindDirectory:
		return "directory"
	case ContentKindSymlink:
		return "symlink"
	case ContentKindGitlink:
		return "gitlink"
	default:
		return "unsupported"
	}
}

// UnmarshalJSON unmarshales the JSON-encoded ContentKind.
func (k *ContentKind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case ContentKindFile.String():
		*k = ContentKindFile
	case ContentKindDirectory.String():
		*k = ContentKindDirectory
	case ContentKindSymlink.String():
		*k = ContentKindSymlink
	case ContentKindGitlink.String():
		*k = ContentKindGitlink
	default:
		*k = ContentKindUnsupported
	}
	return nil
}

// Visibility defines repository visibility.
type Visibility int

// Role values.
const (
	VisibilityUndefined Visibility = iota
	VisibilityPublic
	VisibilityInternal
	VisibilityPrivate
)

// String returns the string representation of Role.
func (v Visibility) String() (s string) {
	switch v {
	case VisibilityPublic:
		return "public"
	case VisibilityInternal:
		return "internal"
	case VisibilityPrivate:
		return "private"
	default:
		return "unknown"
	}
}
