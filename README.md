# goproject
A Sample Golang project with Standard Layout.

It also provides a bootstrap code to begin with for any golang application.

Clone this project and replace the application and module name according to your use case.

### Package Details

| package name    | Details                                                                                                                                                                     |
|-----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| api/types       | Put Request and Response schema of your application here                                                                                                           |
| assets          | Other assets to go along with your repository (images, logos,i18n, migration files ..etc).                                                                                  |
| cmd             | applications of this project having main func(). The directory name for each application should match the name of the executable you want to have (e.g., /cmd/myapp)        |
| http/handler    | All HTTP handler will go here                                                                                                                                               |
| http/middleware | All HTTP middleware will ho here                                                                                                                                            |
| internal/app    | Your main Application code will go here. This is the code you don't want others importing in their applications or libraries. eg /internal/app/myapp                        |
| internal/pkg    | Internal shared code will go here. All code which are shared among interal/apps lies here.                                                                                  |
| pkg             | Library code that's ok to use by external applications. Other projects will import these libraries expecting them to work, so think twice before you put something here :-) |
