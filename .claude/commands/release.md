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

7. **Push and trigger release**: After the user confirms, push to main with `git push`. Then wait for the auto-tag workflow to create the tag (poll with `gh run list --workflow=auto-tag.yml --limit 1` until it completes). Once the tag exists, delete it locally and remotely, then re-create and push it from your local machine. This is necessary because tags created by GitHub Actions' `GITHUB_TOKEN` do not trigger other workflows (like the release build). Steps:
   ```bash
   git fetch --tags
   git tag -d vX.Y.Z
   git push origin :refs/tags/vX.Y.Z
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin vX.Y.Z
   ```
   This ensures the release workflow is triggered by a user-pushed tag.

8. **Summarize**: Show what was done and confirm that the release workflow has been triggered.
