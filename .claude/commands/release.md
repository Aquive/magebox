Prepare a release for MageBox. Follow these steps:

1. **Determine version number**: Read the current version from the `VERSION` file. Look at recent git history (`git log --oneline -20`) and the current CHANGELOG.md to understand what changed since the last release. Based on the changes, suggest an appropriate next version following semver:
   - PATCH: bug fixes, minor updates
   - MINOR: new features, non-breaking changes
   - MAJOR: breaking changes
   Present the suggested version and ask the user to confirm or provide a different one before proceeding.

2. **Update CHANGELOG.md**: Add a new release section at the top (below the header) following the Keep a Changelog format. Analyze all commits since the last release tag to build the changelog entry. Group changes under appropriate headers (Added, Changed, Fixed, Removed, etc.). Include PR references where available. Use today's date.

3. **Update VERSION file**: Write the confirmed version number (without `v` prefix) to the `VERSION` file.

4. **Update VitePress docs if needed**: Check if any new features or commands were added that should be reflected in the documentation under `vitepress/`. Update relevant docs pages if needed.

5. **Run quality checks**: Run `make lint` and `make test` to ensure everything passes.

6. **Create the release commit**: Stage all changed files and create a commit with message `Release vX.Y.Z`. Do NOT push - let the user review and push manually.

7. **Summarize**: Show what was done and remind the user that pushing to main will trigger the auto-tag workflow, which creates the git tag and triggers the release build.
