# insights

the general flow is to get all tasks from the server and filter on the
client side.

* total tasks open (relative to level)
* average task closure time
* range information
  * open vs closed in the last 7 days


given all the dates, we need to track per last 7 days, so we first need last
7 days start timestamp, then compare to 