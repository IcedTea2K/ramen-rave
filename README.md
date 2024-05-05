# Ramen Rave
Teleparty Alternative

# How to run the project
1. First, you must build and compile the necessary Go binaries, before running the extension. There's a convenient script at the root of the project for this
```bash
./build.sh
```
2. Next, navigate to `fire-extension` folder and run `web-ext run`. Make sure to have the tool installed with this [link](https://extensionworkshop.com/documentation/develop/getting-started-with-web-ext/)
```bash
cd fire-extension
web-ext run
```
3. Navigate to supported streaming platform (e.g., [netflix.com](https://www.netflix.com/)), and allow the extension to have Read & Write permission on the website
