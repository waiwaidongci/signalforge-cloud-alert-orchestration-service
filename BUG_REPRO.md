# Bug reproduction

- Bug: permanent migration failures are retried and lose their original cause.
- Trigger: run a migration whose driver returns a permanent failure sentinel.
- Error: the returned error is not identifiable as permanent and retry classification is incorrect.
