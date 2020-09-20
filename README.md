## Purpose

Use the formal hexagonal architecture definition of DDD presented 
https://www.youtube.com/watch?v=oL6JBUk6tj0
by
Kat Zien


## Features
* Makefile to build consistently in a local environment and remote environment
* Dockerfile for a generic image to build for 
* Go Mod (which you should to your project path change)
* VS Code environment
* Generic docker Push (depends on your environment)
* Find files from a directory

## DDD Design 

* Context: find things on the "file system"
* Language: files, directories, paths, file-system, storage, blocks, inodes
* Entities: FindResult
* Value Objects: FilterOption
* Aggregates: FindResults
* Events: Files, file does not exist, permission not valid,...
* Repository: Filesystem, Twitter Likes?



## DDD:

Establish your domain and domain logic
Define your bounded context(s), the model within each context and the ubiquitous
language

Categorizing the building block of your system


* Entity: an object that has an identity, identity could be a row in a database table or other type of identity that differs one entity and the other.
* Value object: an object which we are only interested in its attributes not the identity. Value objects do not have identities and creation/modification of these objects should not be a problem. Thus it is better to be implemented as immutable objects.
* Repository: Entities are usually stored in database, but we would not want to expose our domain logic to database logic so we could have a ‘virtual’ storage for our entities — Repositories.
* Service: Business logic sometimes would involve more than one different objects. When it doesn’t feel right to put a function in either object, we would use a service.


#TODO

* add context
* add dynamic credential ask based on logged in user
* add login system
* add webui
* refactor
