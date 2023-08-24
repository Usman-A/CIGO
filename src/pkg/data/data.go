package data

type Target struct {
	// List of target dependencies
	DependsOn []DependsTarget
	// List of commands that this target can run
	Cmds []string
	// List of path(s) to the generated artifact(s), could be a directory or file
	Artifacts []string
	// Map containing environment variables
	Env map[string]string
}
type DependsTarget struct {
	Project string
	Target  string
}

func (t *Target) AddDependsOn(proj string, target string) {
	t.DependsOn = append(t.DependsOn, DependsTarget{
		Project: proj,
		Target:  target,
	})
}

func (t *Target) AddCmd(cmd string) {
	t.Cmds = append(t.Cmds, cmd)
}

func (t *Target) AddArtifact(artifact string) {
	t.Artifacts = append(t.Artifacts, artifact)
}

func (t *Target) AddEnv(key string, value string) {
	t.Env[key] = value
}

type TargetBuilder struct {
	// target to be built
	target Target
}

func (t *TargetBuilder) Build() Target {
	return t.target
}

func (t *TargetBuilder) SetDependsOn(dependsOn []DependsTarget) {
	t.target.DependsOn = dependsOn
}

func (t *TargetBuilder) SetCmds(cmds []string) {
	t.target.Cmds = cmds
}

func (t *TargetBuilder) SetArtifacts(artifacts []string) {
	t.target.Artifacts = artifacts
}

func (t *TargetBuilder) SetEnv(env map[string]string) {
	t.target.Env = env
}

type ProjectDefinition struct {
	// The main language of the project
	MainLanguage string
	// The language version or standard
	LangVersion string
	// The project Name
	Name string
	// The list of targets. The key is the target name
	Targets map[string]Target
	// Project version
	Version string
	// List  of project owners/maintainers
	Owners []string
	// List  of projects that this project depends on
	DependsOn []string
	// Custom metadata for the project
	Metadata map[string]string
	/* List  of tags that this project affects
	"For example, you can have a common client library that affects all 'client' tags,
	any project that have this tag without having to explicitly list the dependencies."
	- Omar A
	*/
	AffectsTags []string
	// List  of tags that this project is affected by
	AffectedByTags []string
}

func (p *ProjectDefinition) AddTarget(targetName string, target Target) {
	p.Targets[targetName] = target
}

func (p *ProjectDefinition) AddOwner(owner string) {
	p.Owners = append(p.Owners, owner)
}

func (p *ProjectDefinition) AddDependsOn(dependsOn string) {
	p.DependsOn = append(p.DependsOn, dependsOn)
}

func (p *ProjectDefinition) AddMetadata(key string, value string) {
	p.Metadata[key] = value
}

func (p *ProjectDefinition) AddAffectsTag(affectsTag string) {
	p.AffectsTags = append(p.AffectsTags, affectsTag)
}

func (p *ProjectDefinition) AddAffectedByTag(affectedByTag string) {
	p.AffectedByTags = append(p.AffectedByTags, affectedByTag)
}

type ProjectDefinitionBuilder struct {
	// ProjectDefinition to be built
	projectDefinition ProjectDefinition
}

func (p *ProjectDefinitionBuilder) Build() ProjectDefinition {
	return p.projectDefinition
}

func (p *ProjectDefinitionBuilder) SetMainLanguage(mainLang string) {
	p.projectDefinition.MainLanguage = mainLang
}

func (p *ProjectDefinitionBuilder) SetLangVersion(langVers string) {
	p.projectDefinition.LangVersion = langVers
}

func (p *ProjectDefinitionBuilder) SetName(name string) {
	p.projectDefinition.Name = name
}

func (p *ProjectDefinitionBuilder) SetTargets(targets map[string]Target) {
	p.projectDefinition.Targets = targets
}

func (p *ProjectDefinitionBuilder) SetVersion(version string) {
	p.projectDefinition.Version = version
}

func (p *ProjectDefinitionBuilder) SetOwners(owners []string) {
	p.projectDefinition.Owners = owners
}

func (p *ProjectDefinitionBuilder) SetDependsOn(dependsOn []string) {
	p.projectDefinition.DependsOn = dependsOn
}

func (p *ProjectDefinitionBuilder) SetMetadata(metadata map[string]string) {
	p.projectDefinition.Metadata = metadata
}

func (p *ProjectDefinitionBuilder) SetAffectsTags(affectsTags []string) {
	p.projectDefinition.AffectsTags = affectsTags
}

func (p *ProjectDefinitionBuilder) SetAffectedByTags(affectedByTags []string) {
	p.projectDefinition.AffectedByTags = affectedByTags
}

type Workspace struct {
	// List of repository maintainers/owners
	Owners []string
	// The program version associated with this file
	AppVer string
	// List containing paths to the projects
	Projects map[string]string
	// List of available tags
	Tags []string
	// List of required targets to be defined
	RequiredTargets []string
	// Link to where the repository is hosted
	RemoteUrl string
}

func (w *Workspace) AddOwner(owner string) {
	w.Owners = append(w.Owners, owner)
}

func (w *Workspace) AddProject(name string, path string) {
	if w.Projects == nil {
		w.Projects = make(map[string]string)
	}
	w.Projects[name] = path
}

func (w *Workspace) AddTag(tag string) {
	w.Tags = append(w.Tags, tag)
}

func (w *Workspace) AddRequiredTarget(requiredTarget string) {
	w.RequiredTargets = append(w.RequiredTargets, requiredTarget)
}

type WorkspaceBuilder struct {
	// Workspace to be built
	workspace Workspace
}

func (w *WorkspaceBuilder) Build() Workspace {
	return w.workspace
}

func (w *WorkspaceBuilder) SetOwners(owners []string) {
	w.workspace.Owners = owners
}

func (w *WorkspaceBuilder) SetAppVer(appVer string) {
	w.workspace.AppVer = appVer
}

func (w *WorkspaceBuilder) SetProjects(projects map[string]string) {
	w.workspace.Projects = projects
}

func (w *WorkspaceBuilder) SetTags(tags []string) {
	w.workspace.Tags = tags
}

func (w *WorkspaceBuilder) SetRequiredTargets(requiredTargets []string) {
	w.workspace.RequiredTargets = requiredTargets
}

func (w *WorkspaceBuilder) SetRemoteUrl(remoteUrl string) {
	w.workspace.RemoteUrl = remoteUrl
}
