# Cherry-picking to release branches

When a fix lands on `main` and needs to ride a patch release, follow this
sequence.

## 1. Identify the target branch

Patch releases live on `release-MAJOR.MINOR` branches (for example
`release-0.5`). The branch must already exist; if not, ask a maintainer to
create it from the previous tag.

## 2. Cherry-pick

```bash
git fetch origin
git checkout release-0.5
git cherry-pick <commit-sha>
```

## 3. Push

```bash
git push origin release-0.5
```

## 4. Tag

When the patch branch is ready, the maintainer cuts a tag (`v0.5.1`) which
triggers the Release workflow.
