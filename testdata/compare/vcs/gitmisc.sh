# vcs/gitmisc: tags 1.0 and 1.0.0 resolving to the same version, a
# release- prefixed tag, a numeric branch 2.0, a branch with a slash and an
# old tag carrying a different package name (renamed package).
repo gitmisc
put composer.json '{"name": "vcs/oldname"}'
commit "Old name"
tag 0.9.0
put composer.json '{"name": "vcs/gitmisc", "autoload": {"files": ["helpers.php"]}}'
put helpers.php '<?php
function vcs_gitmisc_hello() { return "hi"; }'
commit "Rename"
tag 1.0.0
tag 1.0
put README.md 'release'
commit "Release 1.2"
tag release-1.2.0
branch 2.0
put README.md '2.0 work'
commit "2.0 work"
switch main
branch feature/x
put README.md 'feature'
commit "Feature"
switch main
publish gitmisc
