# Bug reproduction

- Bug: an unknown source API key is classified as an internal failure.
- Trigger: resolve a key that is not present in the source key store.
- Error: the returned error cannot be recognized as a not-found error.
