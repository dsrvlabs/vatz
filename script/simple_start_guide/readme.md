# How to install VATZ & Vatz Official plugins with scripts
> Those following scripts helps user install VATZ in simple way to get started fast. <br>
> Please, be advised that you need to update several variable in scripts to match with your current setting such as cloned VATZ path.

## Here's simple instruction to get started
### 1. Go to script folder 
- You can simply setup and run VATZ with following scripts.

### 2. install_vatz&official-plugins.sh
- Please, update to DEFAULT_VATZ_PATH to your vatz path and run this script to install vatz and vatz-official plugin.
```
bash install_vatz&official-plugins.sh
```

### 3. default_config.yaml
> Please, refer to [installation guide](../../docs/installation.md) for more detailed configs. 
- Replace default.yaml with this file after execute `install_vatz&official-plugins.sh`.
- You must enter hostname and webhook or add more dispatchers such as telegram, pagerduty
- Change the port if necessary.
```
cp default_config.yaml /<vatz_path>/default.yaml
```

### 4. vatz_start.sh
- Running this script will run vatz and vatz-plugin-sei.
- The log is stored in `/var/log/vatz`.
	- You can change the log path if necessary.
- Enter VALOPER_ADDRESS and VOTER_ADDRESS.
- Modify home path to your current setting.
```
bash vatz_start.sh
```

### 5. vatz_stop.sh
- Running this script will stop both vatz and vatz-official-plugin overall that is currently running.
```
bash vatz_stop.sh
```