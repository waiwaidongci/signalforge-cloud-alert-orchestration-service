# Bug reproduction

- Bug: scheduler dispatch closes its result channel before all workers finish.
- Trigger: concurrently dispatch batches containing both successful and failed jobs.
- Error: results are missing, a send-on-closed-channel panic occurs, or the race detector reports a race.
