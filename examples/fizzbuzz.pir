yar theMap be {3: "fizz", 5: "buzz"}.
yar modPrecedences be [3, 5].
yar y be 0.
yar line be "".
i be 1.
4 i < 100:
    line be "".
    y be 0.
    4 y < len(modPrecedences):
        yar divisor be modPrecedences[y].
        if i mod divisor = 0:
            line be line + theMap[divisor].
        .
        y be y + 1.
    .
    ahoy(i + ": " + line).
    i be i + 1.
.
