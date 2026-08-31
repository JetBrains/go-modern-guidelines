# Contributing

Thank you for contributing to Modern Go Guidelines. Every contribution is welcome, from a typo fix to a new guideline.

We are a small team, and we manually verify every change because the accuracy and quality of these guidelines directly affect the code that agents produce. A thoughtful, well-reviewed contribution helps us use that limited review capacity effectively.

## Discuss substantial changes first

Before making a substantial guideline change, open an issue describing the proposal and how you plan to evaluate it. This allows maintainers to confirm the direction before significant work begins. Typo fixes and other obvious corrections do not require prior discussion.

## Keep changes focused

Each pull request should introduce one specific feature or fix one specific problem. Do not combine unrelated guideline changes, refactoring, cleanup, or formatting. Split independent changes into separate pull requests so each can be evaluated and reviewed on its own.

## Evaluate changes to the guidelines

A guideline change should demonstrably help coding agents. Passing unit tests and sounding reasonable are necessary, but they are not enough on their own.

Before submitting a new guideline, an improvement, or a bug fix:

1. Choose representative coding tasks that expose the behavior you want to change.
2. Run the tasks with the current guidelines to establish a baseline.
3. Run the same tasks with your proposed change. Keep the agent, model, prompt, tools, and target Go version the same so the comparison is meaningful.
4. Repeat runs when results may be nondeterministic, and review the generated code for correctness, modern Go style, and unintended regressions.
5. Include the setup, tasks, results, and relevant output in the pull request. Explain how the evidence shows that the change made agent behavior better.

Use `make dev-install` and set `GO_MODERN_GUIDELINES_DEV=1` before launching your agent to test the guidelines from your checkout. Re-run `make dev-install` after each local change. See [Local development](README.md#local-development) for details.

For a bug fix, include a case that demonstrates both the incorrect behavior before the change and the corrected behavior afterward. For a new or revised guideline, also check that it does not encourage invalid changes on unsupported Go versions or in cases where the older form should remain.

## Review your work before opening a pull request

You are responsible for the changes you submit, whether you wrote them by hand or with an AI tool. Before opening a pull request:

- Read and understand the complete diff. Do not submit unreviewed agent output.
- Remove unrelated changes and unnecessary generated or formatting noise.
- Check technical claims and Go version requirements against authoritative sources.
- Make sure examples are valid, focused, and consistent with the guideline.
- After changing `internal/guidelines/guidelines.json`, run `make generate-features`, review the generated `FEATURES.md`, and include it in the pull request. Do not edit `FEATURES.md` directly.
- Run `make test`.

## Write a clear pull request description

Write a clear, self-contained description that is easy to scan. It should explain:

- the problem or opportunity;
- the proposed change and why it is appropriate;
- how the change was evaluated, when it affects guidelines;
- which checks you ran; and
- any limitations or follow-up work.

Link a relevant issue when one exists. A small pull request with a precise description and convincing evidence is much easier to review well.
