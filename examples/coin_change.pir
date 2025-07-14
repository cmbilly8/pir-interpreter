yar min be 
    f(x,y):
        if x < y:
            gives x.
        .
        gives y.
    .
.

yar coinChange be 
    f(coins, amount):
        yar solutions be [0].
        yar x be 0.
        4 x < amount + 1:
            push(solutions, -1).
            x be x + 1.
        .

        yar i be 1.
        yar j be 0.
        yar lastSolAmount be -1.

        4 i <= amount:
            j be 0.
            4 j < len(coins):
                lastSolAmount be i - coins[j].
                if lastSolAmount >= 0 and solutions[lastSolAmount] <> -1:
                    if solutions[i] = -1:
                        solutions[i] be 1 + solutions[lastSolAmount].
                    ls:
                        solutions[i] be min(solutions[i], 1 + solutions[lastSolAmount]).
                    .
                .
                j be j + 1.
            .
            i be i + 1.
        .
        gives solutions[amount].
    .
.

yar tests be [
    [[1,2,5], 11, 3],
    [[2], 3, -1],
    [[1], 0, 0]
].

yar t be 0.
4 t < len(tests):
    result be coinChange(tests[t][0], tests[t][1]).
    if result <> tests[t][2]:
       ahoy("Test " + t + "..." + "FAIL. Expected: " + tests[t][2] + " Got: " + result).
    ls:
        ahoy("Test " + t + "..." + "PASS").
    .
    t be t + 1.
.
