thorf
=====

Interpreter for a massively simplified version of
[Forth](https://en.wikipedia.org/wiki/Forth_(programming_language)). I
originally wrote this as a solution to an [exercism.io](exercism.io) exercise,
but had so much fun that I wanted to take it a bit further.


Supported "words"
-----------------

Basic integer arithmetic operations:

    +: add last two items
    -: subtract last item from second to last item
    *: multiply last two items
    /: divide second to last item by last item

Stack manipulation operations:

    DUP: duplicate the last item
    DROP: remove the last item
    SWAP: swap the order of the last two items
    OVER: duplicate the second to last item

It also supports defining new words at runtime:

    : word-name definition ;
