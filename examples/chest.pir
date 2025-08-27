$ Define a chest.
chest myChest|foo, bar|.

$ Instantiate it with positional or named args to set the chest's items.
yar instance be myChest|"fooVal", "barVal"|. 
yar anotherInstance be myChest|bar: "anotherBarVal", 
                               foo: "anotherFooVal"|.

ahoy(instance).

$ Access props with |
anotherInstance|foo be intance|foo.

ahoy(anotherInstance).
