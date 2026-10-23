# vcs/gitodd: a tag whose composer.json version does not match, a tag with
# a dev suffix, a branch without composer.json and a branch with invalid json.
repo gitodd
put composer.json '{"name": "vcs/gitodd", "version": "3.0.0"}'
commit "Mismatching version"
tag 1.0.0
put composer.json '{"name": "vcs/gitodd"}'
commit "Good"
tag 1.1.0
tag 1.2.0-dev
branch broken
put composer.json '{"name": "vcs/gitodd",'
commit "Broken json"
switch main
branch nojson
git rm -q composer.json
commit "No composer.json"
switch main
publish gitodd
