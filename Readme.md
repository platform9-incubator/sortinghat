
Sorting Hat

Sorting Hat is an application which turns logs into useful alerts. It is primarily designed
to be used for analyzing OpenStack-logs.

It helps users by categorizing similarly looking logs into 'buckets' and let the user:
- Annotate which buckets are interesting and which ones are not
- Annotate the bucket with a  user friendly message

The above two features provide a way to fingerprint and publish 'Error definition' (Pending)
that can be consumed by other users of the same product.

The architecture is very simple pipeline


LogAggregation Service -- filter by keywords (like ERROR) --> SortingHat-logingest (parse and sort into buckets) --> Monog DB

SortingHat UI --> SortingHat UI backend --> Mongo DB

For more information please contact rparikh@platform9.com