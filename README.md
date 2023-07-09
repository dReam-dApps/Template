# Template
Template for building [Dero](https://dero.io) and [dReam dApps](https://dreamdapps.io)

![dReamer](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/dReamerUp.png)

This app serves no purpose other then to help you build on Dero. It uses [Fyne](https://fyne.io/) and [dReams packages](https://github.com/dReam-dApps/dReams) to make a  GUI shell app pre-configured with Dero RPC connection and [Gnomon](https://github.com/civilware/Gnomon) indexer. 

The repo is structured in a way so you can build isolated dApps, as well as ones that can be integrated with other Dero dApps. As the dReams platform evolves this Template will evolve with it. The code comments will currently serve as the main guide on how to build with this Template and integrate into dReams.

### Files
Suggested order of viewing this repo:

1. `template/template_layout.go` 
    - Demonstrates the concept of creating a LayoutAllItems() to create importable dApps. Visit Fyne's [developer guides](https://developer.fyne.io/) for docs on building Fyne layouts.
2. `template/template.go` 
    - Demonstrates tying a dApps process loop into the main dReams process if imported, and has various function examples.
3. `template/start.go` 
     - Contains the StartApp() used to run Template as a stand alone application.

### Usage 

- Install latest [Go version](https://go.dev/doc/install)
- Install [Fyne](https://developer.fyne.io/started/) dependencies
- In your GOPATH clone repo and move into directory with:
```
https://github.com/dReam-dApps/Template.git
cd Template
```
- Make changes to files inside of `Template/template` package directory and run with:
```
go run .
```
- Change directory and package name(s) and import paths accordingly for your import

#### Donations
- *Dero Address*: dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn

![DeroDonations](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/DeroDonations.jpg)

---

#### Licensing

dReams platform and packages are free and open source.    
The source code is published under the [MIT](https://github.com/dReam-dApps/Template/blob/main/LICENSE) License.   
Copyright © 2023 dReam dApps   
