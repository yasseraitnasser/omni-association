# Omni-Association TODO list

## Mandatory
- [x] User authentication (sending JWT as a cookie after login)
- [x] split the handlers in main to different files
- [x] add logout
- [x] serve dashboard for authenticated users only
- [x] in the dashboard page, in the greeting, use the username instead of default it to yasser
- [x] prompt the user to retry if the email or password are invalid
- [x] serve projects for authenticated users only
- [ ] add sections for missing data for creating a new project
- [ ] display the ongoing projects at the bottom of the projects page
- [ ] fill the served pages with actual data from DB instead of dummy data
- [ ] add a page that enables the user to see his contributions, donations ..etc
- [ ] forward to the page the user requested instead of dashboard???

## Bonus
- [ ] login instead of logout for non authenticated users
- [ ] move admin.go file to another directory (it has nothing to do with auth)
- [ ] use a dummy hash for password comparaison
