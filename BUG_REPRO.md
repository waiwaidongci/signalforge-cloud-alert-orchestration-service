# Bug reproduction

- Bug: a successfully retried incident remains in the failed state or is projected inconsistently.
- Trigger: advance an incident from failed through a successful retry and query its history.
- Error: the terminal state or completed history projection is missing.
