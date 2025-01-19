yar x be 'the spot'.

yar y be 10.
$ This is a comment!
yar z be 10 * (20 / 2,3).

$ booleans are ay or nay
$ ? preforms a boolean comparison of the statement between the two preceding double quotes
yar z be "ay or nay"? $ z is ay.

$ Raw strings

yar message be avast'We can do whatever we want in here. /n "'.
chantey message. $ Std out: We can do whatever we want in here. /n "

yar arrrray be arr(1, 2, 3 ,4). $ array literals are a language function.
arr[0] $ You use them the same as any language.

$ Functions, gives returns the result of the following expression. 
f plunder(weapon, battleCry, seadog):
    chantey battleCry.
    attack(weapon).
    yar goldPlundered be loot(seadog).
    $ if gives ends in a end line (.) then it returns nothing.
    gives.

$ These bad boys are first class.
f plunderFac():
    gives plunder.

yar functionVar be plunderFac().
$ or
yar functionVar be f(x): gives x+x.

$ f(): denotes a function literal
$  Conditionals
if ay and nay or nay:
    doSomthing().
lsif ay:
    doOtherThing().
ls:
    cmd/RunShell('ls -a').
.



$  Make sense?
