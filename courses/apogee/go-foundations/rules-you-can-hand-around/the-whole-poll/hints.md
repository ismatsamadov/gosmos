`rule` is a function value on each trip through the loop. Its result tells you whether this check accepts the pass.

A refusal must return from `poll` immediately. Reaching the end means every rule accepted, including when there were no rules.
