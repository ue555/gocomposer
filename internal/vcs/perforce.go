package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Util/Perforce.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gocomposer/internal/phpjson"
)

// SecurityError is Composer\Exception\SecurityException.
type SecurityError struct{ Msg string }

func (e *SecurityError) Error() string { return e.Msg }

// PlainError is a plain PHP \Exception (not a \RuntimeException).
type PlainError struct{ Msg string }

func (e *PlainError) Error() string { return e.Msg }

var (
	rePerforceRshPort = regexp.MustCompile(`(?i)^\s*(?:rsh|jsh)\s*:`)
	rePerforceValid   = regexp.MustCompile(`(?i)^(?:(?:tcp|ssl)(?:4|6|46|64)?:)?(?:\[[0-9a-f:.]+\]|[a-z0-9._][a-z0-9._-]*)(?::[a-z0-9._][a-z0-9._-]*)?$`)
)

// PerforceIsValidPort is Perforce::isValidPort.
func PerforceIsValidPort(u string) bool {
	if rePerforceRshPort.MatchString(u) {
		return false
	}
	return rePerforceValid.MatchString(u)
}

// PerforceCheckServerExists is Perforce::checkServerExists.
func PerforceCheckServerExists(u string, p *Process) bool {
	if !PerforceIsValidPort(u) {
		return false
	}
	code, _, _ := p.Execute([]string{"p4", "-p", u, "info", "-s"}, "")
	return code == 0
}

var (
	p4ExecutableOnce sync.Once
	p4Executable     string
)

// perforceExecutable is Perforce::getP4Executable (ExecutableFinder).
func perforceExecutable() string {
	p4ExecutableOnce.Do(func() {
		p4Executable = "p4"
		if p, err := exec.LookPath("p4"); err == nil {
			p4Executable = p
		}
	})
	return p4Executable
}

// optStr is a nullable PHP string.
type optStr struct {
	v  string
	ok bool
}

func some(s string) optStr { return optStr{s, true} }

// Perforce is Composer\Util\Perforce.
type Perforce struct {
	env              *Env
	path             string
	p4Depot          optStr
	p4Client         optStr
	p4User           optStr
	p4Password       optStr
	p4Port           string
	p4Stream         optStr
	p4DepotType      optStr
	p4Branch         optStr
	uniqueClientName string
	commandResult    string
	// now is the clock (tests)
	now func() time.Time
}

// NewPerforce is Perforce::create (repoConfig nil: no repository config).
func NewPerforce(repoConfig *phpjson.Map, port, path string, env *Env) (*Perforce, error) {
	if !PerforceIsValidPort(port) {
		return nil, &SecurityError{Msg: "Invalid Perforce port (" + port + "), it must be of the form [tcp|ssl:][host:]port"}
	}
	p := &Perforce{env: env, p4Port: port, now: time.Now}
	if err := p.InitializePath(path); err != nil {
		return nil, err
	}
	if err := p.initialize(repoConfig); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Perforce) initialize(repoConfig *phpjson.Map) error {
	host, _ := os.Hostname()
	p.uniqueClientName = host + "_" + strconv.FormatInt(p.now().Unix(), 10)
	if repoConfig == nil || len(repoConfig.Keys()) == 0 {
		return nil
	}
	str := func(k string) (string, bool) {
		v, ok := repoConfig.Get(k)
		if !ok || v == nil {
			return "", false
		}
		s, _ := phpjson.ScalarString(v)
		return s, true
	}
	if s, ok := str("unique_perforce_client_name"); ok {
		p.uniqueClientName = s
	}
	if s, ok := str("depot"); ok {
		p.p4Depot = some(s)
	}
	if s, ok := str("branch"); ok {
		p.p4Branch = some(s)
	}
	if s, ok := str("p4user"); ok {
		p.p4User = some(s)
	} else {
		v, err := p.p4Variable("P4USER")
		if err != nil {
			return err
		}
		p.p4User = v
	}
	if s, ok := str("p4password"); ok {
		p.p4Password = some(s)
	}
	return nil
}

// executeCommand is Perforce::executeCommand.
func (p *Perforce) executeCommand(cmd []string) (int, error) {
	code, out, err := p.env.Process.Execute(cmd, "")
	p.commandResult = out
	return code, err
}

// executeShell runs a string command through the shell like
// ProcessExecutor::execute with a string.
func (p *Perforce) executeShell(cmd string) (int, error) {
	code, out, err := p.env.Process.ExecuteShell(cmd, "")
	p.commandResult = out
	return code, err
}

// ExecuteShell is ProcessExecutor::execute with a command line string (run
// by /bin/sh).
func (p *Process) ExecuteShell(cmd, cwd string) (int, string, error) {
	argv := []string{"/bin/sh", "-c", cmd}
	if p.mock != nil {
		code, out, errOut, err := p.mock(argv, cwd)
		p.mu.Lock()
		p.errOut = errOut
		p.mu.Unlock()
		return code, out, err
	}
	if p.IO != nil && p.IO.IsDebug() {
		safe := rePasswordArg.ReplaceAllString(Sanitize(cmd), "--password '***' ")
		where := cwd
		if where == "" {
			where = "CWD"
		}
		p.IO.WriteError("Executing command (" + where + "): " + safe)
	}
	code, out, errOut, err := p.run(argv, cwd, nil)
	p.mu.Lock()
	p.errOut = errOut
	p.mu.Unlock()
	return code, out, err
}

// executeIn is ProcessExecutor::execute without a cwd while PHP's own
// working directory was changed to dir (chdir).
func (p *Perforce) executeIn(cmd []string, dir string) (int, error) {
	pr := p.env.Process
	if pr.mock != nil {
		code, out, errOut, err := pr.mock(cmd, dir)
		pr.mu.Lock()
		pr.errOut = errOut
		pr.mu.Unlock()
		p.commandResult = out
		return code, err
	}
	pr.outputCommandRun(cmd, "")
	code, out, errOut, err := pr.run(cmd, dir, nil)
	pr.mu.Lock()
	pr.errOut = errOut
	pr.mu.Unlock()
	p.commandResult = out
	return code, err
}

// runWithInput is new Symfony Process($command, null, null, $input)->run():
// no ProcessExecutor logging, the input is written to stdin.
func (p *Perforce) runWithInput(cmd []string, input string) int {
	pr := p.env.Process
	if pr.mock != nil {
		// not a ProcessExecutor call: ProcessExecutorMock never sees it
		return 0
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout, c.Stderr = &out, &out
	if err := c.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode()
		}
		return 127
	}
	return 0
}

// CleanupClientSpec is Perforce::cleanupClientSpec.
func (p *Perforce) CleanupClientSpec() error {
	client := p.Client()
	if _, err := p.executeCommand(p.generateP4Command([]string{"client", "-d", client}, false)); err != nil {
		return err
	}
	_ = os.Remove(p.ClientSpec())
	return nil
}

// Client is Perforce::getClient.
func (p *Perforce) Client() string {
	if !p.p4Client.ok {
		clean := strings.NewReplacer("//", "", "/", "_", "@", "").Replace(p.Stream())
		p.p4Client = some("composer_perforce_" + p.uniqueClientName + "_" + clean)
	}
	return p.p4Client.v
}

// InitializePath is Perforce::initializePath.
func (p *Perforce) InitializePath(path string) error {
	p.path = path
	return ensureDirectoryExists(path)
}

// SetStream is Perforce::setStream.
func (p *Perforce) SetStream(stream string) {
	p.p4Stream = some(stream)
	if idx := strings.LastIndex(stream, "/"); idx > 2 {
		p.p4DepotType = some("stream")
	}
}

// IsStream is Perforce::isStream.
func (p *Perforce) IsStream() bool { return p.p4DepotType.ok && p.p4DepotType.v == "stream" }

// Stream is Perforce::getStream.
func (p *Perforce) Stream() string {
	if !p.p4Stream.ok {
		if p.IsStream() {
			p.p4Stream = some("//" + p.p4Depot.v + "/" + p.p4Branch.v)
		} else {
			p.p4Stream = some("//" + p.p4Depot.v)
		}
	}
	return p.p4Stream.v
}

func streamWithoutLabel(stream string) string {
	if i := strings.Index(stream, "@"); i >= 0 {
		return stream[:i]
	}
	return stream
}

// ClientSpec is Perforce::getP4ClientSpec.
func (p *Perforce) ClientSpec() string { return p.path + "/" + p.Client() + ".p4.spec" }

// User is Perforce::getUser (ok false: null).
func (p *Perforce) User() (string, bool) { return p.p4User.v, p.p4User.ok }

func (p *Perforce) queryP4User() error {
	if p.p4User.v != "" {
		return nil
	}
	v, err := p.p4Variable("P4USER")
	if err != nil {
		return err
	}
	p.p4User = v
	if p.p4User.v != "" {
		return nil
	}
	ans, err := p.env.IO.Ask("Enter P4 User:", nil)
	if err != nil {
		return err
	}
	if s, ok := ans.(string); ok {
		p.p4User = some(s)
	} else {
		p.p4User = optStr{}
	}
	_, err = p.executeShell("export P4USER=" + Escape(p.p4User.v))
	return err
}

// p4Variable is Perforce::getP4variable on non-Windows platforms.
func (p *Perforce) p4Variable(name string) (optStr, error) {
	if _, err := p.executeShell("echo $" + name); err != nil {
		return optStr{}, err
	}
	return some(strings.TrimSpace(p.commandResult)), nil
}

func (p *Perforce) queryP4Password() (optStr, error) {
	if p.p4Password.ok {
		return p.p4Password, nil
	}
	pw, err := p.p4Variable("P4PASSWD")
	if err != nil {
		return optStr{}, err
	}
	if pw.v == "" {
		u, _ := p.User()
		ans, err := p.env.IO.AskAndHideAnswer("Enter password for Perforce user " + u + ": ")
		if err != nil {
			return optStr{}, err
		}
		if s, ok := ans.(string); ok {
			pw = some(s)
		} else {
			pw = optStr{}
		}
	}
	p.p4Password = pw
	return pw, nil
}

// generateP4Command is Perforce::generateP4Command.
func (p *Perforce) generateP4Command(args []string, useClient bool) []string {
	cmd := []string{perforceExecutable()}
	if u, ok := p.User(); ok {
		cmd = append(cmd, "-u", u)
	}
	if useClient {
		cmd = append(cmd, "-c", p.Client())
	}
	cmd = append(cmd, "-p", p.p4Port)
	return append(cmd, args...)
}

// isLoggedIn is Perforce::isLoggedIn.
func (p *Perforce) isLoggedIn() (bool, error) {
	code, err := p.executeCommand(p.generateP4Command([]string{"login", "-s"}, false))
	if err != nil {
		return false, err
	}
	if code != 0 {
		errOut := p.env.Process.ErrorOutput()
		u, _ := p.User()
		if !strings.Contains(errOut, u) {
			if !strings.Contains(errOut, "p4") {
				return false, nil
			}
			return false, &PlainError{Msg: "p4 command not found in path: " + errOut}
		}
		return false, &PlainError{Msg: "Invalid user name: " + u}
	}
	return true, nil
}

// ConnectClient is Perforce::connectClient.
func (p *Perforce) ConnectClient() {
	spec, _ := os.ReadFile(p.ClientSpec())
	p.runWithInput(p.generateP4Command([]string{"client", "-i"}, true), string(spec))
}

// SyncCodeBase is Perforce::syncCodeBase.
func (p *Perforce) SyncCodeBase(label *string) error {
	cmd := p.generateP4Command([]string{"sync", "-f"}, true)
	if label != nil {
		cmd = append(cmd, "@"+*label)
	}
	_, err := p.executeIn(cmd, p.path)
	return err
}

// WriteP4ClientSpec is Perforce::writeP4ClientSpec.
func (p *Perforce) WriteP4ClientSpec() error {
	var b strings.Builder
	u, _ := p.User()
	date := p.now().UTC().Format("2006/01/02 15:04:05")
	b.WriteString("Client: " + p.Client() + "\n\n")
	b.WriteString("Update: " + date + "\n\n")
	b.WriteString("Access: " + date + "\n")
	b.WriteString("Owner:  " + u + "\n\n")
	b.WriteString("Description:\n")
	b.WriteString("  Created by " + u + " from composer.\n\n")
	b.WriteString("Root: " + p.path + "\n\n")
	b.WriteString("Options:  noallwrite noclobber nocompress unlocked modtime rmdir\n\n")
	b.WriteString("SubmitOptions:  revertunchanged\n\n")
	b.WriteString("LineEnd:  local\n\n")
	if p.IsStream() {
		b.WriteString("Stream:\n")
		b.WriteString("  " + streamWithoutLabel(p.p4Stream.v) + "\n")
	} else {
		b.WriteString("View:  " + p.Stream() + "/...  //" + p.Client() + "/... \n")
	}
	if err := os.WriteFile(p.ClientSpec(), []byte(b.String()), 0o666); err != nil {
		return &RuntimeError{Msg: "fopen(" + p.ClientSpec() + "): Failed to open stream: " + err.Error()}
	}
	return nil
}

// P4Login is Perforce::p4Login.
func (p *Perforce) P4Login() error {
	if err := p.queryP4User(); err != nil {
		return err
	}
	ok, err := p.isLoggedIn()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	pw, err := p.queryP4Password()
	if err != nil {
		return err
	}
	if code := p.runWithInput(p.generateP4Command([]string{"login", "-a"}, false), pw.v); code != 0 {
		return &PlainError{Msg: "Error logging in:" + p.env.Process.ErrorOutput()}
	}
	return nil
}

// ComposerInformation is Perforce::getComposerInformation (json_decode:
// nil for invalid or empty content).
func (p *Perforce) ComposerInformation(identifier string) (phpjson.Value, error) {
	content, ok, err := p.FileContent("composer.json", identifier)
	if err != nil || !ok || content == "0" {
		return nil, err
	}
	if !json.Valid([]byte(content)) {
		return nil, nil
	}
	v, err := phpjson.Decode([]byte(content))
	if err != nil {
		return nil, nil
	}
	return v, nil
}

// FileContent is Perforce::getFileContent.
func (p *Perforce) FileContent(file, identifier string) (string, bool, error) {
	path, ok, err := p.FilePath(file, identifier)
	if err != nil || !ok {
		return "", false, err
	}
	if _, err := p.executeCommand(p.generateP4Command([]string{"print", path}, true)); err != nil {
		return "", false, err
	}
	res := p.commandResult
	if t := strings.TrimSpace(res); t == "" || t == "0" {
		return "", false, nil
	}
	return res, true, nil
}

// FilePath is Perforce::getFilePath.
func (p *Perforce) FilePath(file, identifier string) (string, bool, error) {
	idx := strings.Index(identifier, "@")
	if idx < 0 {
		return identifier + "/" + file, true, nil
	}
	path := identifier[:idx] + "/" + file + identifier[idx:]
	if _, err := p.executeCommand(p.generateP4Command([]string{"files", path}, false)); err != nil {
		return "", false, err
	}
	res := p.commandResult
	if !strings.Contains(res, "no such file(s).") {
		if i3 := strings.Index(res, "change"); i3 >= 0 {
			fields := strings.Split(strings.TrimSpace(res[i3:]), " ")
			f1 := ""
			if len(fields) > 1 {
				f1 = fields[1]
			}
			return identifier[:idx] + "/" + file + "@" + f1, true, nil
		}
	}
	return "", false, nil
}

var reNonAlnumSpace = regexp.MustCompile(`[^A-Za-z0-9 ]`)

// Branches is Perforce::getBranches.
func (p *Perforce) Branches() ([]Ref, error) {
	var possible []Ref
	if !p.IsStream() {
		possible = setRef(possible, p.p4Branch.v, p.Stream())
	} else {
		if _, err := p.executeCommand(p.generateP4Command([]string{"streams", "//" + p.p4Depot.v + "/..."}, true)); err != nil {
			return nil, err
		}
		for _, line := range strings.Split(p.commandResult, "\n") {
			bits := strings.Split(line, " ")
			if len(bits) > 4 {
				possible = setRef(possible, reNonAlnumSpace.ReplaceAllString(bits[4], ""), bits[1])
			}
		}
	}
	if _, err := p.executeCommand(p.generateP4Command([]string{"changes", p.Stream() + "/..."}, false)); err != nil {
		return nil, err
	}
	first := strings.Split(p.commandResult, "\n")[0]
	parts := strings.Split(first, " ")
	if len(parts) < 2 {
		return nil, &RuntimeError{Msg: "Undefined array key 1"}
	}
	for _, b := range possible {
		if b.Name == p.p4Branch.v {
			return []Ref{{"master", b.ID + "@" + parts[1]}}, nil
		}
	}
	return nil, &RuntimeError{Msg: `Undefined array key "` + p.p4Branch.v + `"`}
}

// Tags is Perforce::getTags.
func (p *Perforce) Tags() ([]Ref, error) {
	if _, err := p.executeCommand(p.generateP4Command([]string{"labels"}, true)); err != nil {
		return nil, err
	}
	var tags []Ref
	for _, line := range strings.Split(p.commandResult, "\n") {
		if strings.Contains(line, "Label") {
			fields := strings.Split(line, " ")
			if len(fields) > 1 {
				tags = setRef(tags, fields[1], p.Stream()+"@"+fields[1])
			}
		}
	}
	return tags, nil
}

// CheckStream is Perforce::checkStream.
func (p *Perforce) CheckStream() (bool, error) {
	if _, err := p.executeCommand(p.generateP4Command([]string{"depots"}, false)); err != nil {
		return false, err
	}
	for _, line := range strings.Split(p.commandResult, "\n") {
		if strings.Contains(line, "Depot") {
			fields := strings.Split(line, " ")
			if len(fields) > 1 && p.p4Depot.v == fields[1] {
				if len(fields) > 3 {
					p.p4DepotType = some(fields[3])
				} else {
					p.p4DepotType = optStr{}
				}
				return p.IsStream(), nil
			}
		}
	}
	return false, nil
}

// changeList is Perforce::getChangeList (ok false: null).
func (p *Perforce) changeList(reference string) (string, bool, error) {
	idx := strings.Index(reference, "@")
	if idx < 0 {
		return "", false, nil
	}
	if _, err := p.executeCommand(p.generateP4Command([]string{"changes", "-m1", reference[idx:]}, true)); err != nil {
		return "", false, err
	}
	changes := p.commandResult
	if !strings.HasPrefix(changes, "Change") {
		return "", false, nil
	}
	fields := strings.Split(changes, " ")
	if len(fields) < 2 {
		return "", false, nil
	}
	return fields[1], true, nil
}

// CommitLogs is Perforce::getCommitLogs ("" for null).
func (p *Perforce) CommitLogs(from, to string) (string, error) {
	fromCL, ok, err := p.changeList(from)
	if err != nil || !ok {
		return "", err
	}
	toCL, ok, err := p.changeList(to)
	if err != nil || !ok {
		return "", err
	}
	main := from[:strings.Index(from, "@")] + "/..."
	if _, err := p.executeCommand(p.generateP4Command([]string{"filelog", main + "@" + fromCL + "," + toCL}, true)); err != nil {
		return "", err
	}
	return p.commandResult, nil
}
