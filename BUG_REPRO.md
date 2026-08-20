# Bug reproduction

- Bug: a cancelled source probe continues with a background context.
- Trigger: cancel one probe request and immediately issue a second request with a fresh context.
- Error: the cancelled request succeeds, or the second request inherits the first request lifecycle.
