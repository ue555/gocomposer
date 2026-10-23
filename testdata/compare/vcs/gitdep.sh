# vcs/gitdep: requires vcs/gitlib ^1.0, tag 1.0.0, branch main.
repo gitdep
put composer.json '{
    "name": "vcs/gitdep",
    "description": "depends on gitlib",
    "require": {"vcs/gitlib": "^1.0"},
    "autoload": {"classmap": ["lib/"]}
}'
put lib/Dep.php '<?php
class VcsDep {}'
commit "Initial import"
tag 1.0.0
put README.md 'dev'
commit "Dev work"
publish gitdep
