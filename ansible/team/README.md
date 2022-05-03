# Ansible roles to install a current (2022) team machine

This installs
-
- 
- 
- 
- 
- 


     TODO


# Ansible role explanation:
### `browser`
Browser installs and configures browser homepage(s). 

TODO describe how to set homepages.

# `users`
Adds two users called `test` and `contest` which can be used to acces the contest.

### `clientpackages`
Installs the basic set of editors, debuggers, and ide's while removing packages that are unneeded.
`./clientpackages/defaults/main.yaml` contains all packages which will be handled by this role.


Also install submit.py, though this is installed from DOMjudge

    TODO verify the previous statement

### `compilers`

Installs all current compiler aliasses

### `vscodeextensions`

    TODO
