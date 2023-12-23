## Principle
### Maintain a Simple Interface
In the context of the command line interface, the interface consists of arguments, options, and option arguments. Users prefer not to deal with complex commands. Developers should always consider whether it is possible to reduce the number of arguments or options.  
  

### Users: Humans and Computers
There are two types of users: humans and computers.  

When developers introduce interactive interfaces for humans, they should also provide means for computers to achieve the same objectives. Specifically, design the system to allow achieving goals by executing commands from shell scripts or GitHub Actions.


### Don't Stay Silent
In the book "The UNIX Philosophy," there's a section titled "Avoid unnecessary output." The concept of "Avoid unnecessary output" suggests that unnecessary output should not be produced when a command succeeds. 
  
While I appreciate this philosophy, I am not fond of applications that produce no logs whatsoever during execution. It makes me uneasy. In situations where the application remains silent for an extended period, especially when the Rainbow project's commands may alter the infrastructure, having logs during the changes is preferable and more user-friendly.

### Prompt for User Confirmation
Changes to the infrastructure can be irreversible, so it's important to be cautious. Therefore, consider the following points:
  
- When making irreversible changes, obtain user consent through keyboard input.
- Provide the option to skip confirmation using the --force (-f) option (computers may find it challenging to interactively provide consent).
- Allow users to perform a pre-check of the operation with the --dry-run (-n) option.