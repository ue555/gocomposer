# vcs/gitlib: tags v1.0.0 (lightweight), v1.1.0 (annotated), an invalid tag,
# branch 1.x and the default branch main (branch-alias 2.x-dev).
repo gitlib
put composer.json '{
    "name": "vcs/gitlib",
    "description": "git fixture library",
    "license": "MIT",
    "autoload": {"psr-4": {"VcsLib\\": "src/"}}
}'
put src/Lib.php '<?php
namespace VcsLib;

class Lib { const VERSION = "1.0.0"; }'
commit "Initial import"
tag v1.0.0
put src/Lib.php '<?php
namespace VcsLib;

class Lib { const VERSION = "1.1.0"; }'
commit "Release 1.1.0"
tag v1.1.0 "Release 1.1.0"
tag not-a-version
branch 1.x
put src/Lib.php '<?php
namespace VcsLib;

class Lib { const VERSION = "1.x-dev"; }'
commit "Work on 1.x"
switch main
put composer.json '{
    "name": "vcs/gitlib",
    "description": "git fixture library",
    "license": "MIT",
    "require": {"php": ">=7.4"},
    "autoload": {"psr-4": {"VcsLib\\": "src/"}},
    "extra": {"branch-alias": {"dev-main": "2.x-dev"}}
}'
put src/Lib.php '<?php
namespace VcsLib;

class Lib { const VERSION = "2.x-dev"; }'
commit "Start 2.0"
publish gitlib
