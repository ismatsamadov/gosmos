Line up the expected output with the modes the readme describes: which two modes share an action, and which mode does the readme never name?

Write `switch mode { ... }` with one `case` per action. A case can list two values separated by a comma, and `default` covers every mode you don't list.

Shape: `switch mode {` holds `case "NOMINAL":`, `case "SAFE":`, `case "DETUMBLE", ___:` and `default:`, each followed by `return ___` with that action.
