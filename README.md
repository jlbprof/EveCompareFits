# EveCompareFits
Compare 2 Eve Online Ship Fits to find out what new modules you need to buy.

I am responsible for handing out ships to new folk in our corporation.  Often we change the fit, usually just one module, and I need to update the ships in the corp hanger to make sure they conform to the new fit.

Well often it is not easy to determine exactly what changed.

To that end I created this command line tool, written in Go, which allows it to be deployed on any machine, Windows, MacOS or Linux with ease.

`Note:` This program is written in Golang.  It is beyond the scope of the README.md to instruct you on how to install and build go.  But you can go to this website, it is pretty easy:

[Download and install](https://go.dev/doc/install)

To use the program:

    ./EveCompareFits path_to_original_fit path_to_new_fit

If on Windows

    .\EveCompareFits.exe path_to_original_fit path_to_new_fit

Example with real files:

    .\EveCompareFits.exe .\Fits\Moa_Roamer_II.fit  .\Fits\Moa_6th_Rate_Fireship_Hv.fit

It will list 2 things of import:

* What modules/rigs/drones/cargo to remove from the ship.
* What modules/rigs/drones/cargo to add to the ship.

To build for your platform after installing golang.

On Windows:

    go build -o "EveCompareFits.exe" main.go

On other platforms

    go build -o "EveCompareFits" main.go

Go will create the executable file and you can execute as described above.

## Other things you can do

* EveCompareFits --version, shows the version number
* EveCompareFits --justparse fit1 fit2 fit3, will parse each fit, and print out the fit.
  Used mostly to debug parsing.

Things that need to be done:

* I need to learn how to create releases in Go, and generate builds for each platform.
* I know I can use github actions to create the executables for all 3 platforms, what I do not know how to do is tell Github this is the new release and here are the executables.
* Help for above would be appreciated.
* I am not fully confident in the fit parser, so I want to add unit tests for any fits that might be problematic.
* PR's will be reviewed and if I like it, will merge it.

