# ReforgedLauncher

## Hosting ModPacks
Modpacks are hosted on github repositories.

#### Versioning
The launcher uses github's 'releases' api in order to determine versions of a given modpack. There must be at least
 one release (drafts included) on the repository in order to be able to install it.

Release tags should conform to the semantic versioning scheme `MAJOR.MINOR.PATCH` (example: `2.1.3`).

The launcher will interpret `MINOR` and `PATCH` increments as updates that can be safely applied to an existing instance
 without breaking the end-user's saves etc. (so should generally only include iterative, additive, & bug-fix changes to
 mods/resources).

A `MAJOR` increment to a modpack will result in the launcher creating an installation separately to the end-user's current
 game files. This would generally be used for changes to the modpack that could break world-saves/data generated on earlier
 versions of the modpack. 

#### Structure
The basic filetree structure should be as follows:
```
/loader/..
/meta/..
/pack/..
/dependencies.json
```

##### loader
Contains the modloader installer jar if required.  
_Currently only forge is supported._

##### meta
Contains additional information/files related to the modpack.  
_Currently only used to store custom cover images._

##### pack
Contains the files & folders that are copied to the instance's game directory during installation.

##### dependencies.json
Specifies other github repository modpacks that should be installed as part of the current modpack.

#### Optional & Dependant Mods
Folders under the `/packs/` directory that are prefixed with an underscore ('`_`') are interpreted as
 'optional' - ie something that the end-user can choose _not_ to install if they don't want it.

The underscored directory itself is not included in the filepath copied to the instance game dir:
```
repository:
/pack/mods/_My Optional Mod/MyOptionalMod.jar

instance:
../mods/MyOptionalMod.jar
```

Any files or sub-folders within an optional directory become dependent on that option being enabled, and can themselves
 be optional:
```
# OptionalLibraryMod can be enabled/disbaled by the user
# OptionalDependantMod can only be enabled by the user if OptionalLibraryMod is enabled
# DependantModResourcepack is only installed when OptionalDependantMod is enabled

repository:
/pack/_Optional Library Mod/mods/OptionalLibraryMod.jar
/pack/_Optional Library Mod/_Optional_Dependant_Mod/mods/OptionalDependantMod.jar
/pack/_Optional Library Mod/_Optional_Dependant_Mod`/resourcepacks/DependantModResourcepack.zip

instance:
../mods/OptionalLibraryMod.jar
../mods/OptionalDependantMod.jar
../resourcepacks/DependantModResourcepack.zip
```