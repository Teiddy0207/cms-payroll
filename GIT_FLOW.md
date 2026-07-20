# Git Flow for cal-salary

## Branches

- `main`: stable, production-ready code.
- `develop`: integration branch for ongoing work.
- `feature/*`: branch for new features, created from `develop`.
- `release/*`: branch for release preparation, created from `develop`.
- `hotfix/*`: branch for urgent fixes, created from `main`.

## Suggested workflow

1. Create `develop` from `main` once:

```bash
git checkout -b develop
git push -u origin develop
```

2. Start a feature:

```bash
git checkout develop
git checkout -b feature/my-feature
```

3. Finish a feature:

```bash
git checkout develop
git merge --no-ff feature/my-feature
git branch -d feature/my-feature
git push origin develop
```

4. Prepare a release:

```bash
git checkout -b release/1.0.0 develop
```

5. Hotfix from `main`:

```bash
git checkout -b hotfix/critical-fix main
```

## Naming rules

- Feature: `feature/<short-name>`
- Release: `release/<version>`
- Hotfix: `hotfix/<short-name>`

## Recommended convention

- Keep `main` protected.
- Merge only through pull request.
- Use `develop` for day-to-day integration.
- Delete feature branches after merge.
