# go-multirequest

This example shows how to call three long running functions concurrently using channels, wait groups, and mutexes. It assumes that all three functions must complete successfully.

Each function takes 500ms to run. Calling them individually would take 1500ms. The pattern in this example allows all three to return in ~500ms. 

This code might be used when calling three different APIs to gather data required to do something. The use case would be the function should return an error if any one API fails. If partial data from the APIs is ok, don't structure it like this example.  