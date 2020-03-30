// Package main showcases how to use mitchellh/cli to build a command line utility
// ## Overview
//
// 1- create a main function. copy paste the following . Change `EntryPointName`
// to match name of the executable you are making
// 2- change values in version.go as you see fit
//
// To add subcommands, follow the following iterative process:
// 1- for every new subcommand , create a new file with the name of the subcommand
// under `command/` dir. follow the sample `command/version.go` file. in fact, copy-paste
// the content inside of it and modify it to fit your needs
// 2- after previous step , change `Commands` variable inside `init()` in commands.go file
//
// ## Simple Subcommand
//
// every subcommand is difined in a file that shares the same name as the subcommand
// for every subcommand , define a new struct with name `<command name>Command`. Make sure the name is in
// Pascal case so that it is public.
// the struct must have at least one field called `Ui` of type cli.Ui.
// you can add other fields in case the subcommand needs more information to be injected into
// it for the operation to move forward
//
// every subcommand must have three methods :
// - Run : executes the business logic of your application. This is where you execute functions
// that you import from other packages that actually do the task the subcommand expects to perfrom
// - Synopsis : a short synopsis of the subcommand
// - Help : detailed instructions of what the subcommand does. make sure to updare
// the part with `Usage:` in returnvalue and use the name main name of your executable
// because as you are copy pasting code, you may forget to change the name of the binary
//
// ## Flags
//
// My style of cli does not have positional arguments, it only accepts information input
// through flags and environment variables, why ? to simplofy operation for the end user
// and not confuse them and have them wondering whether they should use flags or positional
// arguments to input information
//
// I generally like to create a `flags` package per cli , inside the package, I would
// put all the flags subcommands take. You can also define flags per subcommand inside `command`
// package. It's up to you.
// I put all flags in a cntralized package to make modification and code review simpler because
// in a large library, finding where variables are defined can be challenging. Another reason is because
// there may be multiple subcommands that have the same flags. Having all those refer to the same package
// will make modification and refactoring simpler.
//
// In flags package, as I am declaring functions that extract flags, I would also have them read env variables.
// essentially , for every optional/essential information that needs to be fed into the cli,
// I would first read env variables, then check user input through flags. In case there is user input, it
// would override values passed in through env variables. this is to make sure I don't forget reading
// from env variables. I would also add one last failsafe which is the default values. In case nothing
// was passed in through env variables or flags, I would set the values to a defualt value ( if we can set a default value )
//
// In flags package, there is helper methods in `flag_slice_value.go` I have alse defined a function
// called AppendSliceValue which implements the flag.Value interface and allows multiple
// calls to the same variable to append a list.
//
// ## Complex Subcommands
//
// there may be sub-commands that need more method or configuration.
// In that case create a file with the name of the subcommand won't suffice
// In those case create a package with the same name under `command` directory.
// keep in mind that you can use the same process for simpler commands too in case you
// want to standardizing the way you have written your cli
//
// the newly created package usually has the following files
// - `command.go` : it contains the main struct and entry point of the subcommand
// -
package main
