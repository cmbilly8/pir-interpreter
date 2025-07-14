yar west be
    f():
        ahoy("Hello from west").
        port t.
        gives "west".
    .
.

yar east be
    f():
        ahoy("Hello from east").
        port t.
        gives "east".
    .
.

$ idea: add port t blockade. Where blockade prevents others from porting.
$ The behavior of the other function would be as if it never ported

yar result1 be west().
ahoy("west() returned: " + result1)

yar result2 be east().
ahoy("east() returned: " + result2)
