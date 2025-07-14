yar greeting be f(x):
	gives 'Hello ' + x + "!".
.
ahoy(greeting('world')).

yar adder be 
	f(x, y): 
		gives x + y.
	.

ahoy("Call 1: " + adder(1, 2)).

yar map be {"key1": 1, "key2": 2}.
map["key3"] be adder.
ahoy("Call 2: " + map["key3"](2, 2)).
empty(map)
ahoy(map)
yar x be 1 > 2.

if x:
	ahoy("this shouldn't happen").
lsif !x or x:
	ahoy("aye matey").
lsif 1 + 1 > 1 and ay:
	ahoy("A previous condition prevents us from getting here").
ls:
	ahoy("same with here").
.

if ay:
	ahoy("ay is the logical true, nay is the logical false").
	ahoy(1 = 1).
	ahoy(1 <> 1).
.
