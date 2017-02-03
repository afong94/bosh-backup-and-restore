package instance

import (
	"fmt"
	"path/filepath"
	"strings"
)

type BackupAndRestoreScripts []Script

const (
	backupScriptName           = "p-backup"
	restoreScriptName          = "p-restore"
	preBackupLockScriptName    = "p-pre-backup-lock"
	postBackupUnlockScriptName = "p-post-backup-unlock"

	jobBaseDirectory              = "/var/vcap/jobs/"
	jobDirectoryMatcher           = jobBaseDirectory + "*/bin/"
	backupScriptMatcher           = jobDirectoryMatcher + backupScriptName
	restoreScriptMatcher          = jobDirectoryMatcher + restoreScriptName
	preBackupLockScriptMatcher    = jobDirectoryMatcher + preBackupLockScriptName
	postBackupUnlockScriptMatcher = jobDirectoryMatcher + postBackupUnlockScriptName
)

type Script string

func (s Script) isBackup() bool {
	match, _ := filepath.Match(backupScriptMatcher, string(s))
	return match
}

func (s Script) isRestore() bool {
	match, _ := filepath.Match(restoreScriptMatcher, string(s))
	return match
}

func (s Script) isPreBackupUnlock() bool {
	match, _ := filepath.Match(preBackupLockScriptMatcher, string(s))
	return match
}

func (s Script) isPostBackupUnlock() bool {
	match, _ := filepath.Match(postBackupUnlockScriptMatcher, string(s))
	return match
}

func (s Script) isPlatformScript() bool {
	return s.isBackup() || s.isRestore() || s.isPreBackupUnlock() || s.isPostBackupUnlock()
}

func (s Script) JobName() (string, error) {
	if !strings.HasPrefix(string(s), jobBaseDirectory) {
		return "", fmt.Errorf("script %s is not a Job script", string(s))
	}

	strippedPrefix := strings.TrimPrefix(string(s), jobBaseDirectory)
	splitFirstElement := strings.SplitN(strippedPrefix, "/", 2)
	return splitFirstElement[0], nil
}

func NewBackupAndRestoreScripts(files []string) BackupAndRestoreScripts {
	bandrScripts := []Script{}
	for _, s := range files {
		s := Script(s)
		if s.isPlatformScript() {
			bandrScripts = append(bandrScripts, s)
		}
	}
	return bandrScripts
}
func (s BackupAndRestoreScripts) firstOrBlank() Script {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func (s BackupAndRestoreScripts) HasBackup() bool {
	return len(s.BackupOnly()) > 0
}

func (s BackupAndRestoreScripts) BackupOnly() BackupAndRestoreScripts {
	scripts := BackupAndRestoreScripts{}
	for _, script := range s {
		if script.isBackup() {
			scripts = append(scripts, script)
		}
	}
	return scripts
}

func (s BackupAndRestoreScripts) RestoreOnly() BackupAndRestoreScripts {
	scripts := BackupAndRestoreScripts{}
	for _, script := range s {
		if script.isRestore() {
			scripts = append(scripts, script)
		}
	}
	return scripts
}

func (s BackupAndRestoreScripts) PreBackupLockOnly() BackupAndRestoreScripts {
	scripts := BackupAndRestoreScripts{}
	for _, script := range s {
		if script.isPreBackupUnlock() {
			scripts = append(scripts, script)
		}
	}
	return scripts
}

func (s BackupAndRestoreScripts) PostBackupUnlockOnly() BackupAndRestoreScripts {
	scripts := BackupAndRestoreScripts{}
	for _, script := range s {
		if script.isPostBackupUnlock() {
			scripts = append(scripts, script)
		}
	}
	return scripts
}