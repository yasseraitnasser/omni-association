# Omni-Association TODO list

## Mandatory
- [x] User authentication (sending JWT as a cookie after login)
- [x] split the handlers in main to different files
- [x] add logout
- [x] serve dashboard for authenticated users only
- [x] in the dashboard page, in the greeting, use the username instead of default it to yasser
- [x] prompt the user to retry if the email or password are invalid
- [x] serve projects for authenticated users only
- [x] add sections for missing data for creating a new project
- [x] display the projects at the bottom of the projects page
- [x] fill the served pages with actual data from DB instead of dummy data
- [x] add a page that enables the user to see his contributions, donations ..etc
- [x] add a button to assign committee to project
- [x] add a button to initiate a transaction
- [x] login instead of logout for non authenticated users
- [x] adjust registration protocol
- [x] add membership status in the profile page
- [x] add a button to register a membership payment
- [x] net balance in the project details page
- [x] remove "create a new project+" for non board members (for assigning committees, creating transactions ..)
- [ ] complete invite member handler
- [ ] respond with files when proofs are requested
- [ ] add a button to delete a project(in the project full description page)
- [ ] add a button to change project status (from planning to active, completed ..etc)
- [ ] check why on profile, the normal member's "Payment & Activity History" is broken
- [ ] adjust "next renewal" in the dashboard
- [ ] check for TODOs, FIXs, WARNs ..etc
- [ ] test colours for white theme

## Bonus
- [ ] set a max length for project name (in both backend and frontend)
- [ ] use a dummy hash for password comparaison
- [ ] forward to the page the user requested instead of dashboard
