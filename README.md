# gvm
Go Version Manager


### User manual

1. Set up gvm and Go environments (only need to set once)
    - linux/mac (mac needs to be allowed in setting security): `./gvm setup`
    - windows: `gvm.exe setup`

2. Set the download Go source (optional)
    - view common source: `gvm mirror`
    - set source: `gvm mirror url`
        - example: `gvm mirror https://mirrors.aliyun.com/golang`

2. Specify the version to install Go (example: 1.22.2)
    - `gvm install 1.22.2`

3. View the current Go version
    - `gvm version`

4. Install other versions of Go (example: 1.22.1)
    - `gvm install 1.22.1`

5. Switch Go version
    - `gvm use 1.22.2`
