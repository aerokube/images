#browsers.json

```
	"MicrosoftEdge": {
                "default": "77",
                "versions": {
                        "77": {
                                "image": "selenoid/external-host:1.0.0",
                                "port": "4444",
                                "path": "/",
				"env": ["URLS=[\"http://172.17.0.1:4444/\"]","VNC_PASSWORD=longpassword","SCREEN_RESOLUTION=1921x1080x24"]
                        }
                }
        }

```
