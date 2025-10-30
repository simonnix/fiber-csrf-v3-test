Fiber CSRF V3 Test
==================

Test project to reproduce a bug in Go Fiber CSRF v3-branch.

I was not able to reproduce this bug at runtime (after a build), only in tests.

The "main" branch uses Ginkgo while the "testify" used only Go "testing" package and "testify".

This bug happens intermittently.
