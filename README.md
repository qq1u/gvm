# gvm
Go Version Manager


### User manual

1. Set up gvm and GO environments (only need to set once)
    - linux/mac: `./gvm setup`
    - windows: `gvm.exe setup`

2. Set the download GO source (optional)
    - view common source: `gvm mirror`
    - set source: `gvm mirror url`
        - example: `gvm mirror https://mirrors.aliyun.com/golang`

2. Specify the version to install GO (example: 1.22.2)
    - `gvm install 1.22.2`

3. View the current go version
    - `gvm version`

4. Install other versions of GO (example: 1.22.1)
    - `gvm install 1.22.1`

5. Switch GO version
    - `gvm use 1.22.2`
