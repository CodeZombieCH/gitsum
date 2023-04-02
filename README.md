# gitsum

A CLI application that generates checksum of files contained in a *store* directory and tracks the corresponding checksums files in *repo* directory using git, mirroring the directory structure of the store directory.


## Example

Store:
```
├── dir-a
│   └── file-a:     actual `file-a`
│── dir-b
└── file-c:         actual `file-c`
```

Repo:
```
├── dir-a
│   │── file-a.md5:         md5 checksum of `file-a`
│   │── file-a.sha256:      sha256 checksum of `file-a`
│   └── file-a.sha512:      sha512 checksum of `file-a`
│── dir-b
│── file-c.md5:             md5 checksum of `file-c`
│── file-c.sha256:          sha256 checksum of `file-c`
└── file-c.sha512:          sha512 checksum of `file-c`
```


## Usage

### Create Checksum Files

Create checksum of store in git repository:

    gitsum -repo <repo-dir> -store <store-dir>

where

- `<repo-dir>`: path to the git repository to store the checksum files
- `<store-dir>`: path to the directory of the store, containing all the files to create checksum files for


### Export Checksums

Export checksums from repository to (md5|sha256|sha512sum) compatible checksums files:

    gitsum-export -repo <repo-dir> -hashfile <hash-file> -algorithm <algorithm>

where

- `<repo-dir>`: path to the git repository where the checksum files are stored
- `<hash-file>`: path to checksums file to generate, e.g. MD5SUMS, SHA256SUMS, SHA512SUMS
- `<algorithm>`: hash algorithm to generate export file for. Currently on md5, sha256, sha512 are supported
