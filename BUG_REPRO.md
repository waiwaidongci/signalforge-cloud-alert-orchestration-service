# Bug reproduction

- Bug: concurrent dedup index registration and snapshot lookup expose inconsistent shared state.
- Trigger: start registration and snapshot collection at the same time for the same index.
- Error: the baseline reports a race or returns data changed by another caller.
