# What is okuson-ex?

okuson-ex ("OKUSON extractor") is a program that allows you to fetch all your exercises from [OKUSON](https://github.com/frankluebeck/OKUSON) (a system for online math exercises used at RWTH Aachen University),
and save those exercises in a json or tex file.

Because OKUSON displays different exercises for all users, another function of okuson-ex is combining multiple json files created with okuson-ex by multiple people into one tex file and thus creating a collection of exercises, which could be used for studying.

# How do I compile okuson-ex?

1. [Get and install go](https://golang.org/doc/install)
2. `git clone https://github.com/ninov/okuson-ex`
3. `cd okuson-ex`
4. `go get`
5. `go build`

# How do I use okuson-ex?

Some examples for common use cases:

* To combine all your exercises into a json file called `out.json`, run:  
`okuson-ex -a get -o json -url [your okuson's base url] > out.json`,  
and enter your okuson login credentials
* To combine files out1.json, out2.json and out3.json into a tex file (out.tex), run:  
`okuson-ex -a combine -i out1.json,out2.json,out3.json -tpl [a template file] -o tex > out.tex`  
The template files specifies how the exercises are layouted, see [templates/lainf17.tex](templates/lainf17.tex) for an example.

For further explanation on the command line options, run `okuson-ex` without options.

# Notes

okuson-ex is written by Nino van der Linden and is in no way affiliated with OKUSON or its creators.  
okuson-ex is released under [MIT license](LICENSE).  
okuson-ex was tested with OKUSON 1.4.3.
