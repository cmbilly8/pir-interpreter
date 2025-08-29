$ Define a chest.
chest myChest|foo, bar|.

$ A function to store in the chest.
yar funco be f():
    gives "returns".
.

$ Instantiate it with positional or named args to set the chest's items.
yar instance be myChest|"fooVal", funco|.
yar anotherInstance be myChest|bar: "anotherBarVal", 
                               foo: "anotherFooVal"|.

ahoy(instance).

$ Access props with |
anotherInstance|foo be instance|foo.

ahoy(anotherInstance).

ahoy(instance|bar()).
