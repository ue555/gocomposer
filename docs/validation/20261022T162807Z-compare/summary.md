# Composer vs gocomposer comparison

- gocomposer_sha256: `275660250c9286fd7d3ab03b95ceb956884427c47667e92fd8c970a5db48e5c7`
- host_os: `darwin/arm64`
- oracle_composer_lock_sha256: `7e7ed6afe532a728d4c2cccec04f476c382048aab6148e0ceb1e7594c7f0d34c`
- oracle_image: `localhost/gocomposer-oracle:8.4`
- oracle_image_digest: `sha256:7c3d859839a08d290aa1b5cce89e88d709ddbe6421525d715a3464a651103607`
- oracle_source_head: `f83588392df43d4fd027d84d26de459946ef4e80`
- php_version: `8.4.27`
- upstream_commit: `f83588392df43d4fd027d84d26de459946ef4e80`

| scenario | status | oracle exit | go exit |
|---|---|---|---|
| depends-basic | MATCH | 0 | 0 |
| depends-recursive | DIFF | 0 | 0 |
| depends-tree | DIFF | 0 | 0 |
| depends-php | MATCH | 0 | 0 |
| depends-notfound | MATCH | 1 | 1 |
| depends-locked | MATCH | 0 | 0 |
| prohibits-basic | MATCH | 1 | 1 |
| prohibits-installed | MATCH | 0 | 0 |
| prohibits-php | MATCH | 1 | 1 |
| prohibits-tree | MATCH | 1 | 1 |
| why-alias | MATCH | 0 | 0 |
| depends-nothing-installed | MATCH | 1 | 1 |

## depends-recursive

```
stdout:
+ (go) test/project  -     requires    acme/cli (^1.0)         
+ (go) test/project  -     requires    acme/http (^1.0)        
+ (go) test/project  -     devRequires acme/devtools (^1.0)    
+ (go) test/project  -     requires    acme/log-impl (^1.0)    
+ (go) acme/cli      1.0.0 requires    acme/util (^1.0)        
+ (go) acme/http     1.2.0 requires    acme/util (^1.0)        
+ (go) acme/devtools 1.0.0 requires    acme/log (^1.0 || ^2.0) 
+ (go) acme/log-impl 1.0.0 requires    acme/log (^1.0)         
+ (go) acme/util     1.1.0 requires    acme/log (^1.1)         
- (oracle) test/project  -     requires                   acme/cli (^1.0)         
- (oracle) test/project  -     requires                   acme/http (^1.0)        
- (oracle) test/project  -     requires (for development) acme/devtools (^1.0)    
- (oracle) test/project  -     requires                   acme/log-impl (^1.0)    
- (oracle) acme/cli      1.0.0 requires                   acme/util (^1.0)        
- (oracle) acme/http     1.2.0 requires                   acme/util (^1.0)        
- (oracle) acme/devtools 1.0.0 requires                   acme/log (^1.0 || ^2.0) 
- (oracle) acme/log-impl 1.0.0 requires                   acme/log (^1.0)         
- (oracle) acme/util     1.1.0 requires                   acme/log (^1.1)         
```

## depends-tree

```
stdout:
+ (go) |  `--test/project (devRequires acme/devtools ^1.0)
- (oracle) |  `--test/project (requires (for development) acme/devtools ^1.0)
```
