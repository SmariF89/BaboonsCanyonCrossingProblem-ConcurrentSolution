# CADP - Project 02
## Baboons' canyon crossing.

A concurrent implementation of the famous Baboons' canyon crossing problem.

This is my solution to the problem. It was a project in the course Concurrent and Distributed Programming in Reykjavik University.
It took two weeks to perfect. 

It takes the following concurrency features into account:

### Fairness
- Westheading baboons will have as much of a chance to cross the canyon as the eastheading ones.

### Mutual exclusion
- Westheading baboons will never enter the rope when eastheading baboons are crossing and vice versa. We don't want the baboons
to crash into each other and fall to their deaths.

Additionally there is ###weight control
- The rope's weight will never exceed its limit and break.
