# task

tasks will not be limited to 140, but should have some limit close to 140,
the reasoning is we want tasks to be short, if they're too long, the implication
is to break them down into multiple tasks.

## usernames

1.  all lowercase
1.  characters a-z, 0-9, hyphen, underscore, dot
1.  must start with a-z
1.  length cannot exceed 15 chars

multiple user mentions means that the task is cloned.  however, changes to any
of the clones does not change the original task.

need to accommodate for lists/groups for usernames.

## hashtags

1.  all lowercase
1.  characters: a-z, 0-9
1.  must start with a-z
1.  length cannot exceed 32 chars

## retweet/retask

retask, this is a cloning action.

## embedded urls

need some shortener service

## delegate task

* user must be in your contact list and accepted before you can delegate tasks
to this person
* if user is not in your contact list then the request is rejected
* multiple users can be delegates, this constitutes a clone of the task
* there are several timing issues to be considered:
* the first is if multiple delegates are listed, and not all are in the contacts.
* the other issue is when a contact exists during delegation, but before the task is completed, the delegate is removed.  what happens in these cases.
