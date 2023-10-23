<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## Unreleased

- (chore) [fse-792] Handle accounts with zero staking balance on rewards endpoint
- (chore) [fse-710] Bundle all price cron API calls into a single API call for all tokens
- (chore) [fse-710] Fetch evmos 24h price change and return it on the ERC20ModuleBalance endpoint

### Improvements

- (chore) [fse-536] Adding dependabot
- (chore) [fse-546] Upgrading docker version
- (ci) [fse-673] Deleting codeball
- (ci) [fse-478] Updating github actions and versions
- (ci) [#5](https://github.com/evmos/backend/pull/5) RPC server refactor
- (ci) [#107](https://github.com/evmos/backend/pull/107) Add golangci linter.
