# Bug reproduction

- Bug: a failed item in a batch transaction is committed and its business error is lost.
- Trigger: process multiple transaction items with the second item returning an error.
- Error: the failing item is committed or the original error is not returned.
