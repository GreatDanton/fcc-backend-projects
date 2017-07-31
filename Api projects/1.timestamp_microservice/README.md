# Timestamp microservice

Solved the following set of problems:

* User Story: I can pass a string as a parameter, and it will check to see whether that string contains either a unix timestamp or a natural language date (example: January 1, 2016).
* User Story: If it does, it returns both the Unix timestamp and the natural language form of that date.
* User Story: If it does not contain a date or Unix timestamp, it returns null for those properties.


example usage:
```
url/December%2015,%202015
url/1450137600
```

example output:
```
{ "unix": 1450137600, "natural": "December 15, 2015" }
```