# unisync

Mackup inspired tool to sync application settings, powered by unison

## What it is

unisync helps you keep your application settings (plists, config folders, licenses) synchronized across machines. It does this by using `unison` to diff files, then copies them into a synced folder of your choice

## How is this different from mackup?

unisync does not use symlinks like mackup. So instead of symlinking things into Dropbox/iCloud, unisync uses `unison` for diffing and copies the actual files.

### Why?

I like unison and ran into issues where symlinks randomly broke without me noticing it. I wanted the file content to be identical no matter where it is, including the actual application settings folder.

The sync folder in the cloud acts as a backup and source for diffing, but the apps would continue to work as is if those were deleted.

## Install & Usage

Install `unison` and make sure it's in your path

```
go install github.com/dvcrn/unisync@latest
```

Create a config file in ~/.config/unisync/unisync.yaml:

```yaml
targetPath: ~/.config/appconfigsync
preferDirection: target
apps:
  - dash
  - raycast
```

- `preferDirection` specifies which direction should get picked on conflict. Set it to `target` to have the targetPath take precedence (aka, your dropbox folder). Set it to `local` to say that your local apps should take prededence and override the sync store. Defaults to `target`

### Usage

```
Commands:
  sync - run sync between all enabled apps
  init-from-target - run initial sync, targetPath -> appPath
  init-from-app - run initial sync, appPath -> targetPath
  list - list all available apps
```

The initial sync options **force** one way. For example if you want your local configuration to be overwritten on initial sync with whatever you have in your storage folder, run `init-from-target`

## How to add support for new apps

Check out `apps/` in this repository. Apps are simple yaml files that explain what has to be copied:

```yaml
name: Raycast
friendlyName: raycast
files:
  - basePath: ~/Library/Preferences/
    includedFiles:
      - com.raycast.macos.plist
    ignoredFiles:
      - Name somethingToIgnore.plist
```

- `basePath` is the path that contains the files to sync. It can be `~/Library/Preferences/` when the file in question is just a preferences file, but also soemthing like `~/Library/Application Support/Dash`
- `includedFiles` are all files **within** basePath to process
- `ignoredFiles` is the opposite of `includedFiles` - stuff you don't want to get processed. Has to be Name, Path, BelowPath or Regex

## Troubleshooting & Caveats

- Apps that are currently running will override the config files usually when they're quit
